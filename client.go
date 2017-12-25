package icheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

var apiURL = "https://core.icheck.com.vn"
var apiDevURL = "http://sandbox.icheck.com.vn:4336"
var Dev = true

var AppID string
var Secret string

type Backend interface {
	Call(method, path string, form *RequestValues, params *Params, v interface{}) error
}

// BackendConfiguration is the internal implementation for making HTTP calls to Icheck.
type BackendConfiguration struct {
	URL        string
	HTTPClient *http.Client
}

func GetBackend() Backend {
	api := apiURL
	if Dev == true {
		api = apiDevURL
	}
	return &BackendConfiguration{
		URL:        api,
		HTTPClient: &http.Client{},
	}
}

// Call is the Backend.Call implementation for invoking Icheck APIs.
func (s BackendConfiguration) Call(method, path string, form *RequestValues, params *Params, v interface{}) error {
	var body io.Reader
	if form != nil && !form.Empty() {
		data := form.Encode()
		if strings.ToUpper(method) == "GET" {
			path += "?" + data
		} else {
			body = bytes.NewBufferString(data)
		}
	}

	req, err := s.NewRequest(method, path, "application/x-www-form-urlencoded", body, params)
	if err != nil {
		return err
	}

	if err := s.Do(req, v); err != nil {
		return err
	}

	return nil
}

// NewRequest is used by Call to generate an http.Request. It handles encoding
// parameters and attaching the appropriate headers.
func (s *BackendConfiguration) NewRequest(method, path, contentType string, body io.Reader, params *Params) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = s.URL + path

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		logrus.Debugf("Cannot create Icheck request: %v\n", err)
		return nil, err
	}

	req.SetBasicAuth("icheck", "iYAF&;cBe#G3a~D:#heck")
	req.Header.Add("Content-Type", contentType)

	if params != nil {
		if params.AccessToken != "" {
			req.Header.Add("access-token", params.AccessToken)
		}
		for k, v := range params.Headers {
			for _, line := range v {
				req.Header.Add(k, line)
			}
		}
	}

	return req, nil
}

type Response struct {
	Status int
}

// Do is used by Call to execute an API request and parse the response. It uses
// the backend's HTTP client to execute the request and unmarshals the response
// into v. It also handles unmarshaling errors returned by the API.
func (s *BackendConfiguration) Do(req *http.Request, v interface{}) error {
	logrus.Debugf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)

	start := time.Now()

	res, err := s.HTTPClient.Do(req)

	logrus.Debugf("Completed in %v\n", time.Since(start))

	if err != nil {
		logrus.Debugf("Request to Icheck failed: %v\n", err)
		return err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Debugf("Cannot parse Icheck response: %v\n", err)
		return err
	}
	logrus.Debugf("Icheck Response: %q\n", resBody)
	resData := &Response{}
	if err := json.Unmarshal(resBody, resData); err != nil {
		logrus.Debugf("Cannot parse Icheck response: %v\n", err)
		return err
	}
	if resData.Status >= 400 {
		if resData.Status == 400 {
			badRequest := &ErrBadRequest{}
			if err := json.Unmarshal(resBody, badRequest); err != nil {
				logrus.Debugf("Cannot parse Icheck response: %v\n", err)
				return err
			}
			return badRequest
		}

		return s.ResponseToError(res, resBody)
	}

	if err := json.Unmarshal(resBody, v); err != nil {
		logrus.Debugf("Cannot parse Icheck response: %v\n", err)
		return err
	}
	return nil
}

func (s *BackendConfiguration) ResponseToError(res *http.Response, resBody []byte) error {
	// for some odd reason, the Erro structure doesn't unmarshal
	// initially I thought it was because it's a struct inside of a struct
	// but even after trying that, it still didn't work
	// so unmarshalling to a map for now and parsing the results manually
	// but should investigate later
	err := &Error{}
	json.Unmarshal(resBody, err)
	return err
}

// Error invalid
type ErrBadRequest struct {
	Status            int
	RError            string `json:"error"`
	Summary           string
	InvalidAttributes map[string][]Rule
}

func (e *ErrBadRequest) Error() string {
	for _, value := range e.InvalidAttributes {
		return value[0].Message
	}

	return fmt.Sprintf("Invalid attributes =%s", e.InvalidAttributes)
}

func (e *ErrBadRequest) Invalid() {}

type Rule struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type Error struct {
	Status  int
	Message string
}

// Error serializes the error object to JSON and returns it as a string.
func (e *Error) Error() string {
	return e.Message
}

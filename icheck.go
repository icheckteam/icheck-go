package icheck

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	apiURL     = "https://core.icheck.com.vn/feed"
	uploadsURL = "https://core.icheck.com.vn/feed"
)

// defaultHTTPTimeout is the default timeout on the http.Client used by the library.
// This is chosen to be consistent with the other Icheck language libraries and
// to coordinate with other timeouts configured in the Icheck infrastructure.
const defaultHTTPTimeout = 80 * time.Second

// TotalBackends is the total number of Icheck API endpoints supported by the
// binding.
const TotalBackends = 2

// UnknownPlatform is the string returned as the system name if we couldn't get
// one from `uname`.
const UnknownPlatform = "unknown platform"

// clientversion is the binding version
const clientversion = "24.3.0"

// apiversion is the currently supported API version
const apiversion = "2017-08-12"

// AppInfo contains information about the "app" which this integration belongs
// to. This should be reserved for plugins that wish to identify themselves
// with Icheck.
type AppInfo struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

// formatUserAgent formats an AppInfo in a way that's suitable to be appended
// to a User-Agent string. Note that this format is shared between all
// libraries so if it's changed, it should be changed everywhere.
func (a *AppInfo) formatUserAgent() string {
	str := a.Name
	if a.Version != "" {
		str += "/" + a.Version
	}
	if a.URL != "" {
		str += " (" + a.URL + ")"
	}
	return str
}

// Backend is an interface for making calls against a Icheck service.
// This interface exists to enable mocking for during testing if needed.
type Backend interface {
	Call(method, path, key string, body *RequestValues, params *Params, v interface{}) error
	CallMultipart(method, path, key, boundary string, body io.Reader, params *Params, v interface{}) error
}

// BackendConfiguration is the internal implementation for making HTTP calls to Icheck.
type BackendConfiguration struct {
	Type       SupportedBackend
	URL        string
	HTTPClient *http.Client
}

// SupportedBackend is an enumeration of supported Icheck endpoints.
// Currently supported values are "api" and "uploads".
type SupportedBackend string

const (
	// APIBackend is a constant representing the API service backend.
	APIBackend SupportedBackend = "api"

	// APIURL is the URL of the API service backend.
	APIURL string = "https://core.icheck.com.vn/feed"

	// UploadsBackend is a constant representing the uploads service backend.
	UploadsBackend SupportedBackend = "uploads"

	// UploadsURL is the URL of the uploads service backend.
	UploadsURL string = "https://core.icheck.com.vn/feed"
)

// Backends are the currently supported endpoints.
type Backends struct {
	API Backend
}

// icheckClientUserAgent contains information about the current runtime which
// is serialized and sent in the `X-Icheck-Client-User-Agent` as additional
// debugging information.
type icheckClientUserAgent struct {
	Application     *AppInfo `json:"application"`
	BindingsVersion string   `json:"bindings_version"`
	Language        string   `json:"language"`
	LanguageVersion string   `json:"language_version"`
	Publisher       string   `json:"publisher"`
	Uname           string   `json:"uname"`
}

// Key is the Icheck API key used globally in the binding.
var Key string

// LogLevel is the logging level for this library.
// 0: no logging
// 1: errors only
// 2: errors + informational (default)
// 3: errors + informational + debug
var LogLevel = 2

// Logger controls how Icheck performs logging at a package level. It is useful
// to customise if you need it prefixed for your application to meet other
// requirements
var Logger Printfer

// Printfer is an interface to be implemented by Logger.
type Printfer interface {
	Printf(format string, v ...interface{})
}

func init() {
	Logger = log.New(os.Stderr, "", log.LstdFlags)
	initUserAgent()
}

var appInfo *AppInfo
var httpClient = &http.Client{Timeout: defaultHTTPTimeout}
var backends Backends
var encodedIcheckUserAgent string
var encodedUserAgent string

// SetHTTPClient overrides the default HTTP client.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.
func SetHTTPClient(client *http.Client) {
	httpClient = client
}

// NewBackends creates a new set of backends with the given HTTP client. You
// should only need to use this for testing purposes or on App Engine.
func NewBackends(httpClient *http.Client) *Backends {
	return &Backends{
		API: BackendConfiguration{
			APIBackend, APIURL, httpClient},
	}
}

// GetBackend returns the currently used backend in the binding.
func GetBackend(backend SupportedBackend) Backend {
	var ret Backend
	switch backend {
	case APIBackend:
		if backends.API == nil {
			backends.API = BackendConfiguration{backend, apiURL, httpClient}
		}

		ret = backends.API
	}

	return ret
}

// SetBackend sets the backend used in the binding.
func SetBackend(backend SupportedBackend, b Backend) {
	switch backend {
	case APIBackend:
		backends.API = b
	}
}

// Call is the Backend.Call implementation for invoking Icheck APIs.
func (s BackendConfiguration) Call(method, path, key string, form *RequestValues, params *Params, v interface{}) error {
	var body io.Reader
	var contentType = "application/x-www-form-urlencoded"
	if form != nil && !form.Empty() {
		data := form.Encode()
		if strings.ToUpper(method) == "GET" {
			path += "?" + data
		} else {
			body = bytes.NewBufferString(data)
			contentType = "application/json"
		}
	}

	req, err := s.NewRequest(method, path, key, contentType, body, params)
	if err != nil {
		return err
	}

	if err := s.Do(req, v); err != nil {
		return err
	}

	return nil
}

// CallMultipart is the Backend.CallMultipart implementation for invoking Icheck APIs.
func (s BackendConfiguration) CallMultipart(method, path, key, boundary string, body io.Reader, params *Params, v interface{}) error {
	contentType := "multipart/form-data; boundary=" + boundary

	req, err := s.NewRequest(method, path, key, contentType, body, params)
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
func (s *BackendConfiguration) NewRequest(method, path, key, contentType string, body io.Reader, params *Params) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = s.URL + path

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		if LogLevel > 0 {
			Logger.Printf("Cannot create Icheck request: %v\n", err)
		}
		return nil, err
	}
	req.SetBasicAuth("icheck", "iYAF&;cBe#G3a~D:#heck")
	req.Header.Add("Icheck-Version", apiversion)
	req.Header.Add("User-Agent", encodedUserAgent)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("X-Icheck-Client-User-Agent", encodedIcheckUserAgent)

	if params != nil {
		for k, v := range params.Headers {
			for _, line := range v {
				req.Header.Add(k, line)
			}
		}
	}

	return req, nil
}

// Do is used by Call to execute an API request and parse the response. It uses
// the backend's HTTP client to execute the request and unmarshals the response
// into v. It also handles unmarshaling errors returned by the API.
func (s *BackendConfiguration) Do(req *http.Request, v interface{}) error {
	if LogLevel > 1 {
		Logger.Printf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)
	}

	start := time.Now()

	res, err := s.HTTPClient.Do(req)

	if LogLevel > 2 {
		Logger.Printf("Completed in %v\n", time.Since(start))
	}

	if err != nil {
		if LogLevel > 0 {
			Logger.Printf("Request to Icheck failed: %v\n", err)
		}
		return err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if LogLevel > 0 {
			Logger.Printf("Cannot parse Icheck response: %v\n", err)
		}
		return err
	}

	if res.StatusCode >= 400 {
		return s.ResponseToError(res, resBody)
	}

	if LogLevel > 2 {
		Logger.Printf("Icheck Response: %q\n", resBody)
	}

	if v != nil {
		return json.Unmarshal(resBody, v)
	}

	return nil
}

// ResponseToError ...
func (s *BackendConfiguration) ResponseToError(res *http.Response, resBody []byte) error {
	// for some odd reason, the Erro structure doesn't unmarshal
	// initially I thought it was because it's a struct inside of a struct
	// but even after trying that, it still didn't work
	// so unmarshalling to a map for now and parsing the results manually
	// but should investigate later
	var errMap map[string]interface{}
	json.Unmarshal(resBody, &errMap)

	e, ok := errMap["error"]
	if !ok {
		err := errors.New(string(resBody))
		if LogLevel > 0 {
			Logger.Printf("Unparsable error returned from Icheck: %v\n", err)
		}
		return err
	}

	root := e.(map[string]interface{})

	icheckErr := &Error{
		Type:           ErrorType(root["type"].(string)),
		Msg:            root["message"].(string),
		HTTPStatusCode: res.StatusCode,
		RequestID:      res.Header.Get("Request-Id"),
	}

	if code, ok := root["code"]; ok {
		icheckErr.Code = ErrorCode(code.(string))
	}

	if param, ok := root["param"]; ok {
		icheckErr.Param = param.(string)
	}

	if charge, ok := root["charge"]; ok {
		icheckErr.ChargeID = charge.(string)
	}

	switch icheckErr.Type {
	case ErrorTypeAPI:
		icheckErr.Err = &APIError{icheckErr: icheckErr}

	case ErrorTypeAPIConnection:
		icheckErr.Err = &APIConnectionError{icheckErr: icheckErr}

	case ErrorTypeAuthentication:
		icheckErr.Err = &AuthenticationError{icheckErr: icheckErr}

	case ErrorTypeCard:
		cardErr := &CardError{icheckErr: icheckErr}
		icheckErr.Err = cardErr

		if declineCode, ok := root["decline_code"]; ok {
			cardErr.DeclineCode = declineCode.(string)
		}

	case ErrorTypeInvalidRequest:
		icheckErr.Err = &InvalidRequestError{icheckErr: icheckErr}

	case ErrorTypePermission:
		icheckErr.Err = &PermissionError{icheckErr: icheckErr}

	case ErrorTypeRateLimit:
		icheckErr.Err = &RateLimitError{icheckErr: icheckErr}
	}

	if LogLevel > 0 {
		Logger.Printf("Error encountered from Icheck: %v\n", icheckErr)
	}

	return icheckErr
}

// SetAppInfo sets app information. See AppInfo.
func SetAppInfo(info *AppInfo) {
	if info != nil && info.Name == "" {
		panic(fmt.Errorf("App info name cannot be empty"))
	}
	appInfo = info

	// This is run in init, but we need to reinitialize it now that we have
	// some app info.
	initUserAgent()
}

// getUname tries to get a uname from the system, but not that hard. It tries
// to execute `uname -a`, but swallows any errors in case that didn't work
// (i.e. non-Unix non-Mac system or some other reason).
func getUname() string {
	path, err := exec.LookPath("uname")
	if err != nil {
		return UnknownPlatform
	}

	cmd := exec.Command(path, "-a")
	var out bytes.Buffer
	cmd.Stderr = nil // goes to os.DevNull
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return UnknownPlatform
	}

	return out.String()
}

func initUserAgent() {
	encodedIcheckUserAgent = "Icheck/v1 GoBindings/" + clientversion
	if appInfo != nil {
		encodedUserAgent += " " + appInfo.formatUserAgent()
	}

	icheckUserAgent := &icheckClientUserAgent{
		Application:     appInfo,
		BindingsVersion: clientversion,
		Language:        "go",
		LanguageVersion: runtime.Version(),
		Publisher:       "icheck",
		Uname:           getUname(),
	}
	marshaled, err := json.Marshal(icheckUserAgent)
	// Encoding this struct should never be a problem, so we're okay to panic
	// in case it is for some reason.
	if err != nil {
		panic(err)
	}
	encodedIcheckUserAgent = string(marshaled)
}

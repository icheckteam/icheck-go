package icheck

import (
	"log"
	"net/http"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestLogin(t *testing.T) {
	client := &BackendConfiguration{
		URL:        "https://core.icheck.com.vn",
		HTTPClient: &http.Client{},
	}

	logrus.SetLevel(logrus.DebugLevel)
	form := &RequestValues{}
	form.Add("username", "0977465849")
	form.Add("password", "12345678")

	res := &LoginResponse{}
	err := client.Call("POST", "login", form, nil, res)

	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", res)
}

package account

import (
	"log"
	"net/http"
	"testing"

	"github.com/Sirupsen/logrus"
	icheck "github.com/icheckteam/icheck-go"
)

func TestLogin(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	backend := &icheck.BackendConfiguration{
		URL:        "https://core.icheck.com.vn",
		HTTPClient: &http.Client{},
	}

	client := &Client{B: backend}
	accessToken, err := client.Login(&icheck.LoginParams{
		Username: "0977465849",
		Password: "12345678",
	})

	if err != nil {
		t.Fatal(err)
	}

	user, err := client.Me(&icheck.Params{
		AccessToken: accessToken.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	log.Print(user)
}

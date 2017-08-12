package icheck

import "testing"

func TestCheckinBackendConfigurationNewRequestWithStripeAccount(t *testing.T) {
	c := BackendConfiguration{URL: APIURL}
	p := &Params{}

	_, err := c.NewRequest("", "", "", "", nil, p)
	if err != nil {
		t.Fatal(err)
	}
}

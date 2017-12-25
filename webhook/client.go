package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"time"
)

// Computes a webhook signature using Stripe's v1 signing method. See
// https://stripe.com/docs/webhooks#signatures
func computeSignature(t time.Time, payload []byte, appID string, secret string) []byte {
	mac := hmac.New(sha1.New, []byte(appID+":"+secret))
	mac.Write([]byte(fmt.Sprintf("%d", t.Unix())))
	mac.Write([]byte("."))
	mac.Write(payload)
	return mac.Sum(nil)
}

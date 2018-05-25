package twirest

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"sort"
	"strings"
)

// ErrSignature is used to signify that ValidateRequest was not able
// to confirm the received signature as being correct.
var ErrSignature = errors.New("Invalid X-Twilio-Signature")

// ValidateRequest will take the incoming http.Request on a Twilio callback
// URL and verify that the X-Twilio-Signature is correct. fullUrl must be
// in the format of https://www.exampe.com/twilio/callback?param1=exists
//
// This can be used to create custom middleware to make it easy to verify
// the request is sent from Twilio. Documentation about this method is
// available at https://www.twilio.com/docs/api/security#validating-requests
func (twiClient *TwilioClient) ValidateRequest(fullUrl string, req *http.Request) error {
	twilioSignature := req.Header.Get("X-Twilio-Signature")
	if twilioSignature == "" {
		return ErrSignature
	}

	var toEncode []string

	toEncode = append(toEncode, fullUrl)

	if req.Method == "POST" {
		err := req.ParseForm()
		if err != nil {
			return err
		}

		var sortedParams []string
		for value, _ := range req.PostForm {
			sortedParams = append(sortedParams, value)
		}
		sort.Strings(sortedParams)

		for _, key := range sortedParams {
			toEncode = append(toEncode, key)
			toEncode = append(toEncode, req.FormValue(key))
		}
	}

	validateString := strings.Join(toEncode, "")

	mac := hmac.New(sha1.New, []byte(twiClient.authToken))
	mac.Write([]byte(validateString))
	hashed := mac.Sum(nil)
	expectedSignature := base64.StdEncoding.EncodeToString(hashed)

	if twilioSignature == expectedSignature {
		return nil
	} else {
		return ErrSignature
	}
}

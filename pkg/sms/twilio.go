package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// TwilioMessager is a Twilio message sending struct
type TwilioMessager struct {
	client *http.Client
	sid    string
	token  string
	from   string
}

// NewTwilioMessager configures and returns a new TwilioMessager
func NewTwilioMessager(c *http.Client, sid, token, from string) TwilioMessager {
	return TwilioMessager{
		client: c,
		sid:    sid,
		token:  token,
		from:   from,
	}
}

// Send sends an SMS containing `message` via twilio to `number`
func (t TwilioMessager) Send(ctx context.Context, number, message string) (string, error) {
	req, err := t.SMSRequest(number, message)
	if err != nil {
		return "", errors.Wrap(err, "failed to get sms request")
	}
	res, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode >= 300 {
		b, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return "", errors.Wrap(e, "failed to read twilio response body")
		}
		return "", errors.Errorf("failed to send twilio message: %s", string(b))
	}

	var d map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&d)
	if err != nil {
		return "", errors.Wrap(err, "unable to decode the twilio API response")
	}
	// Return the SID from the response
	return fmt.Sprintf("%v", d["sid"]), nil
}

// SMSRequest configures an HTTP request to send to the Twilio REST API
func (t TwilioMessager) SMSRequest(to, body string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf(`https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json`, t.sid))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse sms request url")
	}
	d := url.Values{}
	d.Set("To", to)
	d.Set("From", t.from)
	d.Set("Body", body)

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(d.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(t.sid, t.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

// ValidateIncomingRequest returns an error if the incoming req could not be
// validated as coming from Twilio.
//
// This process is frequently error prone, especially if you are running behind
// a proxy, or Twilio is making requests with a port in the URL.
// See https://www.twilio.com/docs/security#validating-requests for more information
func ValidateIncomingRequest(host string, authToken string, req *http.Request) (err error) {
	err = req.ParseForm()
	if err != nil {
		return
	}
	err = validateIncomingRequest(host, authToken, req.URL.String(), req.Form, req.Header.Get("X-Twilio-Signature"))
	if err != nil {
		return
	}

	return
}

func validateIncomingRequest(host string, authToken string, URL string, postForm url.Values, xTwilioSignature string) (err error) {
	expectedTwilioSignature := GetExpectedTwilioSignature(host, authToken, URL, postForm)
	if xTwilioSignature != expectedTwilioSignature {
		err = errors.New("Bad X-Twilio-Signature")
		return
	}

	return
}

func GetExpectedTwilioSignature(host string, authToken string, URL string, postForm url.Values) (expectedTwilioSignature string) {
	// Take the full URL of the request URL you specify for your
	// phone number or app, from the protocol (https...) through
	// the end of the query string (everything after the ?).
	str := host + URL

	// If the request is a POST, sort all of the POST parameters
	// alphabetically (using Unix-style case-sensitive sorting order).
	keys := make([]string, 0, len(postForm))
	for key := range postForm {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Iterate through the sorted list of POST parameters, and append
	// the variable name and value (with no delimiters) to the end
	// of the URL string.
	for _, key := range keys {
		str += key + postForm[key][0]
	}

	// Sign the resulting string with HMAC-SHA1 using your AuthToken
	// as the key (remember, your AuthToken's case matters!).
	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(str))
	expectedMac := mac.Sum(nil)

	// Base64 encode the resulting hash value.
	expectedTwilioSignature = base64.StdEncoding.EncodeToString(expectedMac)

	// Compare your hash to ours, submitted in the X-Twilio-Signature header.
	// If they match, then you're good to go.
	return expectedTwilioSignature
}

package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

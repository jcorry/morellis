package sms_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jcorry/morellis/pkg/sms"
	"github.com/stretchr/testify/require"
)

var c *http.Client

func TestTwilioMessager_SMSRequest(t1 *testing.T) {
	type fields struct {
		client *http.Client
		sid    string
		token  string
		from   string
	}
	type args struct {
		to   string
		body string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: c,
				sid:    "foo",
				token:  "bar",
				from:   "212-867-5309",
			},
			args: args{
				to:   "404-515-0400",
				body: "test message",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t1.Run(tt.name, func(t *testing.T) {
			tm := sms.NewTwilioMessager(tt.fields.client, tt.fields.sid, tt.fields.token, tt.fields.from)
			req, err := tm.SMSRequest(tt.args.to, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("SMSRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Request should...
			// have headers set
			require.Equal(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
			require.Equal(t, req.Header.Get("Accept"), "application/json")
			// have URL set with sid
			require.Equal(t, req.URL.String(), `https://api.twilio.com/2010-04-01/Accounts/foo/Messages.json`)
			// have values in the body
			body, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			require.Equal(t, `Body=test+message&From=212-867-5309&To=404-515-0400`, string(body))
		})
	}
}

func TestTwilioMessager_Send(t *testing.T) {
	type fields struct {
		client *http.Client
		sid    string
		token  string
		from   string
	}
	type args struct {
		ctx     context.Context
		number  string
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: NewTestClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"sid": "message-sid"}`)),
						Header:     make(http.Header),
					}
				}),
				sid:   "foo",
				token: "bar",
				from:  "404-555-1212",
			},
			args: args{
				ctx:     context.TODO(),
				number:  "212-867-5309",
				message: "Test Message",
			},
			want:    "message-sid",
			wantErr: false,
		},
		{
			name: "err: bad twilio request",
			fields: fields{
				client: NewTestClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: 400,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`bad request`)),
						Header:     make(http.Header),
					}
				}),
				sid:   "foo",
				token: "bar",
				from:  "404-555-1212",
			},
			args: args{
				ctx:     context.TODO(),
				number:  "212-867-5309",
				message: "Test Message",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "err: undecodable response",
			fields: fields{
				client: NewTestClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"sid": "message-sid"`)),
						Header:     make(http.Header),
					}
				}),
				sid:   "foo",
				token: "bar",
				from:  "404-555-1212",
			},
			args: args{
				ctx:     context.TODO(),
				number:  "212-867-5309",
				message: "Test Message",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tm := sms.NewTwilioMessager(tt.fields.client, tt.fields.sid, tt.fields.token, tt.fields.from)
			got, err := tm.Send(tt.args.ctx, tt.args.number, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Logf("%v", err)
			}
			if got != tt.want {
				t.Errorf("Send() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// RoundTripFunc to allow mocking the http.Client Transport
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient provides a client with a mock transport layer
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

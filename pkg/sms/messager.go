package sms

import "context"

// Messager provides and interface for sendng SMS messages
//go:generate counterfeiter . Messager
type Messager interface {
	Send(ctx context.Context, number, message string) (string, error)
}

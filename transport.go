package iap

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"golang.org/x/oauth2"
)

var (
	errOauth2TransportSource = errors.New("oauth2: Transport's Source is nil")
	errMissingIDToken        = errors.New("missing IDToken")
)

// from oauth2.Transport
// to Use IdToken for Authorization Header

type Transport struct {
	Base   http.RoundTripper
	Source oauth2.TokenSource
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}

	return http.DefaultTransport
}

func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}

	return r2
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false

	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	if t.Source == nil {
		return nil, errOauth2TransportSource
	}

	token, err := t.Source.Token()
	if err != nil {
		return nil, err
	}

	req2 := cloneRequest(req) // per RoundTripper contract

	// token.SetAuthHeader(req2)
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errMissingIDToken
	}

	req2.Header.Set("Authorization", token.Type()+" "+idToken)

	// req.Body is assumed to be closed by the base RoundTripper.
	reqBodyClosed = true

	return t.base().RoundTrip(req2)
}

var cancelOnce sync.Once

// CancelRequest does nothing. It used to be a legacy cancellation mechanism
// but now only it only logs on first use to warn that it's deprecated.
//
// Deprecated: use contexts for cancellation instead.
func (t *Transport) CancelRequest(req *http.Request) {
	cancelOnce.Do(func() {
		log.Printf("deprecated: golang.org/x/oauth2: Transport.CancelRequest no longer does anything; use contexts")
	})
}

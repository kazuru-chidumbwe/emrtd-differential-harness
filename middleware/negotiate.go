// Package middleware implements §VIII explicit-reject on PACE failure before BAC fallback.
package middleware

import (
	"errors"
	"fmt"

	"github.com/gmrtd/gmrtd/bac"
	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/pace"
)

// ErrPaceFailed is returned when PACE fails and BAC fallback is not explicitly allowed.
var ErrPaceFailed = errors.New("pace failed: explicit reject (middleware §VIII)")

// Options controls downgrade behaviour for NegotiatePACEBAC.
type Options struct {
	// AllowBACFallback permits gmrtd-style silent PACE→BAC when PACE fails.
	AllowBACFallback bool
}

// Session holds negotiation results after NegotiatePACEBAC.
type Session struct {
	PaceErr       error
	BacErr        error
	BacSuccess    bool
	SurfacedError error
}

// NegotiatePACEBAC runs PACE then optionally BAC, enforcing explicit-reject unless fallback is allowed.
func NegotiatePACEBAC(nfc *iso7816.NfcSession, doc *document.Document, pass *password.Password, opts Options) Session {
	var sess Session

	_, _, paceErr := pace.NewPace(nfc, doc, pass).DoPACE()
	sess.PaceErr = paceErr

	if paceErr != nil && !opts.AllowBACFallback {
		sess.SurfacedError = fmt.Errorf("%w: %v", ErrPaceFailed, paceErr)
		return sess
	}

	if nfc.SM() == nil {
		bacResult, bacErr := bac.NewBAC(nfc, doc, pass).DoBAC()
		sess.BacErr = bacErr
		sess.BacSuccess = bacResult != nil && bacResult.Success
	}

	return sess
}

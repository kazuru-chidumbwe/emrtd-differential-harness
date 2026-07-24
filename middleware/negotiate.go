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

// ErrPaceFailed is returned when PACE was advertised, fails, and BAC fallback is not allowed.
var ErrPaceFailed = errors.New("pace failed: explicit reject (middleware §VIII)")

// Options controls downgrade behaviour for NegotiatePACEBAC.
type Options struct {
	// AllowBACFallback permits silent PACE→BAC when advertised PACE fails.
	AllowBACFallback bool
}

// Session holds negotiation results after NegotiatePACEBAC.
type Session struct {
	PaceErr       error
	BacErr        error
	BacSuccess    bool
	SurfacedError error
	PaceOffered   bool
}

// runPace is the PACE entry used by NegotiatePACEBAC (overridable in tests).
var runPace = func(nfc *iso7816.NfcSession, doc *document.Document, pass *password.Password) error {
	_, _, err := pace.NewPace(nfc, doc, pass).DoPACE()
	return err
}

// paceAdvertised reports whether EF.CardAccess lists PACE parameters.
func paceAdvertised(doc *document.Document) bool {
	if doc == nil || doc.Mf.CardAccess == nil || doc.Mf.CardAccess.SecurityInfos == nil {
		return false
	}
	return len(doc.Mf.CardAccess.SecurityInfos.PaceInfos) > 0
}

// NegotiatePACEBAC runs PACE when advertised, then optionally BAC.
// Explicit-reject applies only when PACE was advertised and failed and AllowBACFallback is false.
// BAC-only chips (no PACE params) proceed to BAC without a false-positive reject.
func NegotiatePACEBAC(nfc *iso7816.NfcSession, doc *document.Document, pass *password.Password, opts Options) Session {
	var sess Session
	sess.PaceOffered = paceAdvertised(doc)

	if sess.PaceOffered {
		paceErr := runPace(nfc, doc, pass)
		sess.PaceErr = paceErr

		if paceErr != nil && !opts.AllowBACFallback {
			sess.SurfacedError = fmt.Errorf("%w: %v", ErrPaceFailed, paceErr)
			return sess
		}
	}

	if nfc.SM() == nil {
		bacResult, bacErr := bac.NewBAC(nfc, doc, pass).DoBAC()
		sess.BacErr = bacErr
		sess.BacSuccess = bacResult != nil && bacResult.Success
	}

	return sess
}

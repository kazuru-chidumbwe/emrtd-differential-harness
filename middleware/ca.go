// CA middleware: explicit-reject when chip authentication fails (§VIII analogue).
package middleware

import (
	"errors"
	"fmt"

	"github.com/gmrtd/gmrtd/chipauth"
	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
)

// ErrChipAuthFailed is returned when chip authentication fails and continue is not allowed.
var ErrChipAuthFailed = errors.New("chip authentication failed: explicit reject (middleware §VIII)")

// CAOptions controls behaviour after chip authentication failure.
type CAOptions struct {
	AllowContinue bool
}

// CAResult holds chip authentication outcomes.
type CAResult struct {
	ChipAuthErr   error
	ChipAuthOK    bool
	SurfacedError error
}

// runChipAuth is the CA entry used by PerformChipAuth (overridable in tests).
var runChipAuth = func(nfc *iso7816.NfcSession, doc *document.Document) (ok bool, err error) {
	result, err := chipauth.NewChipAuth(nfc, doc).DoChipAuth()
	return result != nil && result.Success, err
}

// PerformChipAuth runs EAC-CA and optionally surfaces failure to the caller.
func PerformChipAuth(nfc *iso7816.NfcSession, doc *document.Document, opts CAOptions) CAResult {
	var res CAResult
	ok, err := runChipAuth(nfc, doc)
	res.ChipAuthErr = err
	res.ChipAuthOK = ok

	if err != nil && !opts.AllowContinue {
		res.SurfacedError = fmt.Errorf("%w: %v", ErrChipAuthFailed, err)
	}
	return res
}

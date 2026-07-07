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

// PerformChipAuth runs EAC-CA and optionally surfaces failure to the caller.
func PerformChipAuth(nfc *iso7816.NfcSession, doc *document.Document, opts CAOptions) CAResult {
	var res CAResult
	result, err := chipauth.NewChipAuth(nfc, doc).DoChipAuth()
	res.ChipAuthErr = err
	res.ChipAuthOK = result != nil && result.Success

	if err != nil && !opts.AllowContinue {
		res.SurfacedError = fmt.Errorf("%w: %v", ErrChipAuthFailed, err)
	}
	return res
}

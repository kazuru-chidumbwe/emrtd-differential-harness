// AA middleware: explicit-reject when Active Authentication fails.
package middleware

import (
	"errors"
	"fmt"

	"github.com/gmrtd/gmrtd/activeauth"
	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
)

// ErrActiveAuthFailed is returned when AA fails and continue is not allowed.
var ErrActiveAuthFailed = errors.New("active authentication failed: explicit reject (middleware AA)")

// AAOptions controls behaviour after Active Authentication failure.
type AAOptions struct {
	AllowContinue bool
}

// AAResult holds Active Authentication outcomes.
type AAResult struct {
	ActiveAuthErr error
	ActiveAuthOK  bool
	SurfacedError error
}

// runActiveAuth is the AA entry used by PerformActiveAuth (overridable in tests).
var runActiveAuth = func(nfc *iso7816.NfcSession, doc *document.Document) (ok bool, err error) {
	result, err := activeauth.NewActiveAuth(nfc, doc).DoActiveAuth()
	return result != nil && result.Success, err
}

// PerformActiveAuth runs AA and optionally surfaces failure to the caller.
func PerformActiveAuth(nfc *iso7816.NfcSession, doc *document.Document, opts AAOptions) AAResult {
	var res AAResult
	ok, err := runActiveAuth(nfc, doc)
	res.ActiveAuthErr = err
	res.ActiveAuthOK = ok

	if err != nil && !opts.AllowContinue {
		res.SurfacedError = fmt.Errorf("%w: %v", ErrActiveAuthFailed, err)
	}
	return res
}

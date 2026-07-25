// Package normfail defines a cross-library failure taxonomy for AA/TA/EAC result JSON.
package normfail

import (
	"regexp"
	"strings"
)

// Failure is embedded in harness run JSON as normalized_failure.
type Failure struct {
	Mechanism    string `json:"mechanism"`
	Step         string `json:"step"`
	Iso7816SW    string `json:"iso7816_sw,omitempty"`
	FailureClass string `json:"failure_class"`
	Surfaced     bool   `json:"surfaced"`
}

const (
	ClassChipSWReject     = "chip_sw_reject"
	ClassPeerUnsupported  = "peer_unsupported"
	ClassProtocolException = "protocol_exception"
)

var (
	reAdvertised = regexp.MustCompile(`(?i)status[:\s]*([0-9a-f]{4})`)
	reHexSW      = regexp.MustCompile(`(?i)(?:SW\s*=\s*0x|0x)([0-9a-f]{4})`)
)

// ExtractSW pulls an ISO 7816 status word from a library error string when present.
func ExtractSW(msg string) string {
	if msg == "" {
		return ""
	}
	if m := reHexSW.FindStringSubmatch(msg); len(m) == 2 {
		return strings.ToUpper(m[1])
	}
	if m := reAdvertised.FindStringSubmatch(msg); len(m) == 2 {
		return strings.ToUpper(m[1])
	}
	return ""
}

// ChipSW builds a chip status-word rejection record.
func ChipSW(mechanism, step, sw string, surfaced bool) Failure {
	return Failure{
		Mechanism:    mechanism,
		Step:         step,
		Iso7816SW:    strings.ToUpper(sw),
		FailureClass: ClassChipSWReject,
		Surfaced:     surfaced,
	}
}

// FromErr builds a chip_sw_reject when SW can be parsed, else protocol_exception.
func FromErr(mechanism, step, errMsg string, surfaced bool) Failure {
	sw := ExtractSW(errMsg)
	if sw != "" {
		return ChipSW(mechanism, step, sw, surfaced)
	}
	return Failure{
		Mechanism:    mechanism,
		Step:         step,
		FailureClass: ClassProtocolException,
		Surfaced:     surfaced,
	}
}

// PeerUnsupported marks a library that lacks the mechanism.
func PeerUnsupported(mechanism string) Failure {
	return Failure{
		Mechanism:    mechanism,
		Step:         "n/a",
		FailureClass: ClassPeerUnsupported,
		Surfaced:     true,
	}
}

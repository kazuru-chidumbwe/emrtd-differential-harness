package profile

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type MRZ struct {
	DocumentNumber string `json:"document_number"`
	DateOfBirth    string `json:"date_of_birth"`
	DateOfExpiry   string `json:"date_of_expiry"`
}

type Injection struct {
	PaceFailOn string `json:"pace_fail_on"`
	PaceSW     string `json:"pace_sw"`
	// PaceChannel selects how PACE fails at the APDU boundary:
	//   "" / "sw" (default): return configured pace_sw (honest synthetic chip).
	//   "timeout": brief delay then empty response (timed exchange abstraction).
	//   "no_response": empty frame / incomplete exchange (jammed mid-handshake).
	//   "transport_abort": empty response treated as hard transport failure.
	// Channel modes are lab-only synthetic abstractions — not RF/NFC fidelity.
	PaceChannel string `json:"pace_channel,omitempty"`
	Notes       string `json:"notes"`
}

type CAInjection struct {
	CaFailOn string `json:"ca_fail_on"`
	CaSW     string `json:"ca_sw"`
	Notes    string `json:"notes"`
}

type AAInjection struct {
	AaFailOn string `json:"aa_fail_on"`
	AaSW     string `json:"aa_sw"`
	Notes    string `json:"notes"`
}

type TAInjection struct {
	TaFailOn string `json:"ta_fail_on"`
	TaSW     string `json:"ta_sw"`
	Notes    string `json:"notes"`
}

type Profile struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Mechanism     string       `json:"mechanism"`
	Condition     string       `json:"condition"`
	Tier          string       `json:"tier"`
	MRZ           MRZ          `json:"mrz"`
	CardAccessHex string       `json:"card_access_hex"`
	Dg14HexPath   string       `json:"dg14_hex_path,omitempty"`
	Dg15HexPath   string       `json:"dg15_hex_path,omitempty"`
	Injection     Injection    `json:"injection,omitempty"`
	CAInjection   CAInjection  `json:"ca_injection,omitempty"`
	AAInjection   AAInjection  `json:"aa_injection,omitempty"`
	TAInjection   TAInjection  `json:"ta_injection,omitempty"`
	// PeerSupport documents asymmetric mechanisms (e.g. TA: gmrtd unsupported).
	PeerSupport map[string]string `json:"peer_support,omitempty"`
}

func Load(path string) (*Profile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read profile: %w", err)
	}
	var p Profile
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("parse profile: %w", err)
	}
	if p.ID == "" || p.MRZ.DocumentNumber == "" {
		return nil, fmt.Errorf("profile missing required fields (id, mrz)")
	}
	if p.CardAccessHex == "" && p.Dg14HexPath == "" && p.Dg15HexPath == "" {
		// BAC-only success-path profiles intentionally omit CardAccess.
		// TA/EAC asymmetric profiles may omit DG fixtures when using synthetic PKI paths.
		if p.Condition != "bac_only_success" && p.Mechanism != "BAC" &&
			p.Mechanism != "TA" && p.Mechanism != "EAC" {
			return nil, fmt.Errorf("profile needs card_access_hex, dg14_hex_path, and/or dg15_hex_path")
		}
	}
	return &p, nil
}

func LoadHexFile(path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, string(raw))
	return hex.DecodeString(s)
}

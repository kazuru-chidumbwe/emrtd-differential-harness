package middleware

import (
	"errors"
	"testing"

	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
)

func TestPerformChipAuth_explicitReject(t *testing.T) {
	old := runChipAuth
	runChipAuth = func(*iso7816.NfcSession, *document.Document) (bool, error) {
		return false, errors.New("ca inject fail")
	}
	defer func() { runChipAuth = old }()

	res := PerformChipAuth(nil, &document.Document{}, CAOptions{AllowContinue: false})
	if res.SurfacedError == nil {
		t.Fatal("expected surfaced error on CA failure")
	}
	if !errors.Is(res.SurfacedError, ErrChipAuthFailed) {
		t.Fatalf("want ErrChipAuthFailed, got %v", res.SurfacedError)
	}
	if res.ChipAuthOK {
		t.Fatal("ChipAuthOK should be false")
	}
}

func TestPerformChipAuth_allowContinue(t *testing.T) {
	old := runChipAuth
	runChipAuth = func(*iso7816.NfcSession, *document.Document) (bool, error) {
		return false, errors.New("ca inject fail")
	}
	defer func() { runChipAuth = old }()

	res := PerformChipAuth(nil, &document.Document{}, CAOptions{AllowContinue: true})
	if res.SurfacedError != nil {
		t.Fatalf("unexpected surface when AllowContinue: %v", res.SurfacedError)
	}
	if res.ChipAuthErr == nil {
		t.Fatal("expected ChipAuthErr recorded")
	}
}

func TestPerformChipAuth_successNoFalsePositive(t *testing.T) {
	old := runChipAuth
	runChipAuth = func(*iso7816.NfcSession, *document.Document) (bool, error) {
		return true, nil
	}
	defer func() { runChipAuth = old }()

	res := PerformChipAuth(nil, &document.Document{}, CAOptions{AllowContinue: false})
	if res.SurfacedError != nil {
		t.Fatalf("false positive on CA success: %v", res.SurfacedError)
	}
	if !res.ChipAuthOK {
		t.Fatal("expected ChipAuthOK")
	}
}

func TestPerformActiveAuth_explicitReject(t *testing.T) {
	old := runActiveAuth
	runActiveAuth = func(*iso7816.NfcSession, *document.Document) (bool, error) {
		return false, errors.New("aa inject fail")
	}
	defer func() { runActiveAuth = old }()

	res := PerformActiveAuth(nil, &document.Document{}, AAOptions{AllowContinue: false})
	if res.SurfacedError == nil {
		t.Fatal("expected surfaced error on AA failure")
	}
	if !errors.Is(res.SurfacedError, ErrActiveAuthFailed) {
		t.Fatalf("want ErrActiveAuthFailed, got %v", res.SurfacedError)
	}
}

func TestPerformActiveAuth_allowContinue(t *testing.T) {
	old := runActiveAuth
	runActiveAuth = func(*iso7816.NfcSession, *document.Document) (bool, error) {
		return false, errors.New("aa inject fail")
	}
	defer func() { runActiveAuth = old }()

	res := PerformActiveAuth(nil, &document.Document{}, AAOptions{AllowContinue: true})
	if res.SurfacedError != nil {
		t.Fatalf("unexpected surface when AllowContinue: %v", res.SurfacedError)
	}
}

func TestPerformActiveAuth_successNoFalsePositive(t *testing.T) {
	old := runActiveAuth
	runActiveAuth = func(*iso7816.NfcSession, *document.Document) (bool, error) {
		return true, nil
	}
	defer func() { runActiveAuth = old }()

	res := PerformActiveAuth(nil, &document.Document{}, AAOptions{AllowContinue: false})
	if res.SurfacedError != nil {
		t.Fatalf("false positive on AA success: %v", res.SurfacedError)
	}
	if !res.ActiveAuthOK {
		t.Fatal("expected ActiveAuthOK")
	}
}

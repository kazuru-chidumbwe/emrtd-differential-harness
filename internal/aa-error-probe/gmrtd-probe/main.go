package main

import (
	"fmt"
	"os"

	"github.com/gmrtd/gmrtd/activeauth"
	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/utils"
)

func main() {
	dg15 := utils.HexToBytes("6f8201023081ff300d06092a864886f70d01010105000381ed003081e90281e100bb8f93f4dc95e205cda17c6927ab1e365b13065d03cd12e0fce95d96840529453202f56cc4c13f77cd062930c8bc89a2873b257045c286e601cf3c09323a53103314902804aa10a314628ce222206a8866946a36b442041bb54ac81e6855dd1d6e16101833d65a191c20ac8b33b8a1a32920f46043f8031cf2bc17417030865fc5be5a39dee423bcba3ca8177168eb23cfe01ba43ec87711b1cfff85db46f300dd8ae317b50d543b573e119e23af7070d0b2fed6a3b2313a5ec02a531aaed1741f4390d1013e2a0f081eac5dc8b0a1b2c6bdb1206f08d30e3643e1e5bdf536110203010001")

	fmt.Println("=== PROBE A: gmrtd DoActiveAuth with chip SW=6982 ===")
	var doc document.Document
	if err := doc.NewDG(15, dg15); err != nil {
		fmt.Println("setup DG15 fail:", err)
		os.Exit(2)
	}
	nfc := iso7816.NewNfcSession(&iso7816.StaticTransceiver{RApdu: utils.HexToBytes("6982")})
	aa := activeauth.NewActiveAuth(nfc, &doc)
	aa2, err := aa.WithChallenge([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	if err != nil {
		fmt.Println("WithChallenge:", err)
		os.Exit(2)
	}
	result, err := aa2.DoActiveAuth()
	if result != nil {
		fmt.Printf("result.Success=%v\n", result.Success)
	} else {
		fmt.Println("result=nil")
	}
	fmt.Printf("err=%v\n", err)
	if err == nil {
		fmt.Println("VERDICT_LIB: ERROR NOT RETURNED")
	} else {
		fmt.Println("VERDICT_LIB: error returned to caller")
	}

	fmt.Println()
	fmt.Println("=== PROBE B: ValidateActiveAuthSignature empty response ===")
	dg15obj, e := document.NewDG15(dg15)
	if e != nil {
		fmt.Println(e)
		os.Exit(2)
	}
	r2, err2 := activeauth.ValidateActiveAuthSignature(dg15obj, nil, []byte{1, 2, 3, 4, 5, 6, 7, 8})
	if r2 != nil {
		fmt.Printf("result.Success=%v\n", r2.Success)
	} else {
		fmt.Println("result=nil")
	}
	fmt.Printf("err=%v\n", err2)
	if err2 == nil {
		fmt.Println("VERDICT_SIG: ERROR NOT RETURNED")
	} else {
		fmt.Println("VERDICT_SIG: signature validation returns error")
	}

	fmt.Println()
	fmt.Println("=== PROBE C: reader-shaped swallow (same pattern as performChipAuthentication) ===")
	// Mirrors reader.go: record result+err, always continue (return nil to ReadDocument).
	sessionActiveAuthResult, sessionActiveAuthErr := aa2.DoActiveAuth()
	readerStepErr := error(nil) // performChipAuthentication returns nil
	fmt.Printf("Session.ActiveAuthResult.Success=%v\n", sessionActiveAuthResult != nil && sessionActiveAuthResult.Success)
	fmt.Printf("Session.ActiveAuthErr=%v\n", sessionActiveAuthErr)
	fmt.Printf("performChipAuthentication return err=%v\n", readerStepErr)
	if sessionActiveAuthErr != nil && readerStepErr == nil {
		fmt.Println("VERDICT_READER: AA failure recorded but step returns nil (silent continue — same class as CA/PACE record-only)")
	}

	fmt.Println()
	fmt.Println("=== PROBE D: CLI surfacing gap (static) ===")
	fmt.Println("cmd/gmrtd-reader has surfacePaceErr; grep surfaceActiveAuthErr => absent")
	fmt.Println("VERDICT_CLI: ActiveAuthErr would NOT abort process exit unless host adds surfacing")
}

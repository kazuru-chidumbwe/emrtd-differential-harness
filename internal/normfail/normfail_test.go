package normfail

import "testing"

func TestExtractSW(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Advertised(insecure) status:6982", "6982"},
		{"Exception (SW = 0x6982: SECURITY STATUS NOT SATISFIED)", "6982"},
		{"forced SW 6FFF", ""},
		{"", ""},
	}
	for _, tc := range cases {
		got := ExtractSW(tc.in)
		if got != tc.want {
			t.Fatalf("ExtractSW(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestFromErr(t *testing.T) {
	f := FromErr("AA", "internal_authenticate", "status:6982", false)
	if f.FailureClass != ClassChipSWReject || f.Iso7816SW != "6982" || f.Surfaced {
		t.Fatalf("unexpected %#v", f)
	}
	u := PeerUnsupported("TA")
	if u.FailureClass != ClassPeerUnsupported || !u.Surfaced {
		t.Fatalf("unexpected %#v", u)
	}
}

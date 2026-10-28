package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckBundledGoogleSans(t *testing.T) {
	shellPath := t.TempDir()
	fontPath := filepath.Join(shellPath, "DCommon", "assets", "fonts", "google-sans-flex", "GoogleSansFlex.ttf")
	missing := checkConfiguredFont("UI font", "Google Sans Flex", shellPath, nil, false, "")
	if missing.status != statusWarn {
		t.Fatalf("missing bundled font status = %v, want %v", missing.status, statusWarn)
	}
	if err := os.MkdirAll(filepath.Dir(fontPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fontPath, []byte("font"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, family := range []string{"Google Sans Flex", " google sans flex "} {
		result := checkConfiguredFont("UI font", family, shellPath, nil, false, "")
		if result.status != statusOK {
			t.Fatalf("bundled font %q status = %v, want %v", family, result.status, statusOK)
		}
	}
}

func TestQtFontFamilies(t *testing.T) {
	fcOutput := "RobotoMono Nerd Font\tBold\tGOOG\n" +
		"RobotoMono Nerd Font\tBold Italic\tpyrs\n" +
		"Nimbus Sans\tRegular\tURW \n" +
		"Nimbus Sans\tBold\tUKWN\n" +
		"Plex,Plex Lt\tLight,Regular\tIBM\n" +
		"Mono,Mono NF,Mono NF Light\tLight,Regular\tNONE\n"
	families := parseQtFontFamilies(fcOutput)

	cases := map[string]bool{
		"RobotoMono Nerd Font [GOOG]": true,
		"robotomono nerd font [pyrs]": true,
		"RobotoMono Nerd Font":        true,
		"Nimbus Sans [URW ]":          true,
		"Plex Lt":                     true,
		"Mono NF":                     true,
		"Mono NF Light":               false,
		"Monospace":                   true,
		"Missing Font [GOOG]":         false,
	}
	for name, want := range cases {
		if got := fontInFcList(name, families); got != want {
			t.Errorf("fontInFcList(%q) = %v, want %v", name, got, want)
		}
	}
}

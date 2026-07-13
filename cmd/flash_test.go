package cmd

import (
	"testing"

	"github.com/zarhus/benchctl/platform"
)

func TestParseFlashTarget(t *testing.T) {
	cases := []struct {
		args    []string
		want    platform.FlashTarget
		wantErr bool
	}{
		{nil, platform.FlashHost, false},
		{[]string{"host"}, platform.FlashHost, false},
		{[]string{"bmc"}, platform.FlashBMC, false},
		{[]string{"wat"}, 0, true},
	}
	for _, tc := range cases {
		got, err := parseFlashTarget(tc.args)
		if tc.wantErr {
			if err == nil {
				t.Errorf("parseFlashTarget(%q) = nil error, want error", tc.args)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseFlashTarget(%q) error: %v", tc.args, err)
			continue
		}
		if got != tc.want {
			t.Errorf("parseFlashTarget(%q) = %v, want %v", tc.args, got, tc.want)
		}
	}
}

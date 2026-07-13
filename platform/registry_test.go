package platform

import (
	"reflect"
	"testing"

	"github.com/zarhus/benchctl/exec"
)

func TestRegisterAndGet(t *testing.T) {
	spec := Spec{
		Name:            "test-register-and-get",
		DefaultUser:     "root",
		DefaultPassword: "secret",
		New:             func(exec.Runner, Config) Platform { return nil },
	}
	Register(spec)

	got, ok := Get("test-register-and-get")
	if !ok {
		t.Fatal("Get returned ok=false for a registered driver")
	}
	if got.DefaultPassword != "secret" || got.DefaultUser != "root" {
		t.Errorf("Get returned %+v, want the registered spec", got)
	}
}

func TestGetUnknownReturnsFalse(t *testing.T) {
	if _, ok := Get("no-such-platform"); ok {
		t.Error("Get returned ok=true for an unregistered name")
	}
}

func TestDefaultEmptyUntilSet(t *testing.T) {
	if got := Default(); got != "" {
		t.Errorf("Default() = %q, want empty before SetDefault", got)
	}
	SetDefault("default-test")
	defer SetDefault("")
	if got := Default(); got != "default-test" {
		t.Errorf("Default() = %q, want default-test", got)
	}
}

func TestNamesAreSorted(t *testing.T) {
	Register(Spec{Name: "zeta-names-test"})
	Register(Spec{Name: "alpha-names-test"})

	names := Names()
	// The two test names must appear, alpha before zeta.
	var ai, zi = -1, -1
	for i, name := range names {
		switch name {
		case "alpha-names-test":
			ai = i
		case "zeta-names-test":
			zi = i
		}
	}
	if ai == -1 || zi == -1 || ai > zi {
		t.Errorf("Names() = %v, want alpha-names-test before zeta-names-test", names)
	}
	if !reflect.DeepEqual(names, append([]string(nil), names...)) {
		t.Fatal("unreachable") // guards against accidental nil handling
	}
}

func TestFlashTargetString(t *testing.T) {
	if FlashHost.String() != "host" {
		t.Errorf("FlashHost.String() = %q, want host", FlashHost.String())
	}
	if FlashBMC.String() != "bmc" {
		t.Errorf("FlashBMC.String() = %q, want bmc", FlashBMC.String())
	}
}

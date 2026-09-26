//go:build linux

package systray

import (
	"testing"

	"github.com/godbus/dbus/v5/introspect"
)

func hasMethod(methods []introspect.Method, name string) bool {
	for _, m := range methods {
		if m.Name == name {
			return true
		}
	}
	return false
}

func TestSNIIntrospectionHidesUnhandledMethods(t *testing.T) {
	prevLeft, prevRight := tappedLeft, tappedRight
	t.Cleanup(func() { tappedLeft, tappedRight = prevLeft, prevRight })

	tappedLeft, tappedRight = nil, nil
	iface := sniIntrospection()
	if hasMethod(iface.Methods, "Activate") {
		t.Error("Activate must not be advertised without SetOnTapped")
	}
	if hasMethod(iface.Methods, "SecondaryActivate") {
		t.Error("SecondaryActivate must not be advertised without SetOnSecondaryTapped")
	}

	tappedLeft = func() {}
	tappedRight = func() {}
	iface = sniIntrospection()
	if !hasMethod(iface.Methods, "Activate") {
		t.Error("Activate must be advertised after SetOnTapped")
	}
	if !hasMethod(iface.Methods, "SecondaryActivate") {
		t.Error("SecondaryActivate must be advertised after SetOnSecondaryTapped")
	}
}

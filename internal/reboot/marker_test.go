package reboot

import (
	"path/filepath"
	"testing"
	"time"
)

func TestRebootedAfter(t *testing.T) {
	created := time.Date(2026, 11, 6, 12, 47, 9, 0, time.UTC)
	cases := []struct {
		name     string
		decided  *time.Time
		bootTime time.Time
		want     bool
	}{
		{name: "booted after the remediation was created", bootTime: created.Add(30 * time.Second), want: true},
		{name: "booted before, no self reboot recorded", bootTime: created.Add(-10 * time.Second), want: false},
		{name: "rebooted itself shortly before the remediation was created", decided: ptr(created.Add(-13 * time.Second)), bootTime: created.Add(-11 * time.Second), want: true},
		{name: "self reboot recorded but the node did not boot since", decided: ptr(created.Add(-13 * time.Second)), bootTime: created.Add(-time.Hour), want: false},
		{name: "self reboot too long before the remediation was created", decided: ptr(created.Add(-10 * time.Minute)), bootTime: created.Add(-9 * time.Minute), want: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			selfRebootMarkerPath = filepath.Join(t.TempDir(), "var", "lib", "self-node-remediation", "self-reboot")
			if c.decided != nil {
				if err := RecordSelfReboot(*c.decided); err != nil {
					t.Fatal(err)
				}
			}
			got, err := RebootedAfter(c.bootTime, created)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Fatalf("RebootedAfter = %v, want %v", got, c.want)
			}
		})
	}
}

func ptr(t time.Time) *time.Time {
	return &t
}

func TestRebootedItselfRecently(t *testing.T) {
	now := time.Date(2026, 11, 6, 14, 42, 1, 0, time.UTC)
	cases := []struct {
		name     string
		decided  *time.Time
		bootTime time.Time
		want     bool
	}{
		{name: "no self reboot recorded", bootTime: now.Add(-time.Minute), want: false},
		{name: "booted after a self reboot a minute ago", decided: ptr(now.Add(-63 * time.Second)), bootTime: now.Add(-60 * time.Second), want: true},
		{name: "self reboot decided but not booted since", decided: ptr(now.Add(-10 * time.Second)), bootTime: now.Add(-time.Hour), want: false},
		{name: "self reboot too long ago", decided: ptr(now.Add(-10 * time.Minute)), bootTime: now.Add(-9 * time.Minute), want: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			selfRebootMarkerPath = filepath.Join(t.TempDir(), "self-reboot")
			if c.decided != nil {
				if err := RecordSelfReboot(*c.decided); err != nil {
					t.Fatal(err)
				}
			}
			got, err := RebootedItselfRecently(c.bootTime, now)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Fatalf("RebootedItselfRecently = %v, want %v", got, c.want)
			}
		})
	}
}

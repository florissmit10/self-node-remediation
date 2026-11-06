package reboot

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const SelfRebootLeadWindow = 2 * time.Minute

var selfRebootMarkerPath = "/proc/1/root/var/lib/self-node-remediation/self-reboot"

func RecordSelfReboot(at time.Time) error {
	dir := filepath.Dir(selfRebootMarkerPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(selfRebootMarkerPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(at.UTC().Format(time.RFC3339Nano)); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	for _, d := range []string{dir, filepath.Dir(dir)} {
		if err := syncDir(d); err != nil {
			return err
		}
	}
	return nil
}

func syncDir(path string) error {
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}

func LastSelfReboot() (time.Time, bool, error) {
	data, err := os.ReadFile(selfRebootMarkerPath)
	if errors.Is(err, fs.ErrNotExist) {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}
	at, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(string(data)))
	if err != nil {
		return time.Time{}, false, err
	}
	return at, true, nil
}

func RebootedAfter(bootTime, remediationCreated time.Time) (bool, error) {
	if bootTime.After(remediationCreated) {
		return true, nil
	}
	decided, ok, err := LastSelfReboot()
	if err != nil || !ok {
		return false, err
	}
	return bootTime.After(decided) && decided.After(remediationCreated.Add(-SelfRebootLeadWindow)), nil
}

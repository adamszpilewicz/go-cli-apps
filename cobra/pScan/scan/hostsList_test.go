package scan_test

import (
	"errors"
	"pscan/scan"
	"testing"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name   string
		host   string
		expLen int
		expErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostList{}
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}
			err := hl.Add(tc.host)
			// tests if err expected
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %q, got %q instead\n",
						tc.expErr, err)
				}
				return
			}

			// tests if err not expected
			if err != nil {
				t.Fatalf("Expected no error, got %q instead\n", err)
			}
			if len(hl.Hosts) != tc.expLen {
				t.Errorf("Expected list length %d, got %d instead\n",
					tc.expLen, len(hl.Hosts))
			}
			if hl.Hosts[1] != tc.host {
				t.Errorf("Expected host name %q as index 1, got %q instead\n",
					tc.host, hl.Hosts[1])
			}
		})
	}

}

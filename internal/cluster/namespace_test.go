package cluster

import (
	"testing"
)

func TestCmd(t *testing.T) {
	cases := []struct {
		name  string
		isErr bool
	}{
		{
			name:  "joe",
			isErr: false,
		},
		{
			name:  "joes-2nd-feature",
			isErr: false,
		},
		{
			name:  "joes-realllllllyyyyyyy-looooonnnnnggggg-nnaaaammmeeee",
			isErr: true,
		},
		{
			name:  "joE",
			isErr: true,
		},
		{
			name:  "joe/feature",
			isErr: true,
		},
	}

	for _, test := range cases {
		err := CheckName(test.name)
		if err != nil && !test.isErr {
			t.Errorf("unexpected test error for %s: %v", test.name, err)
			continue
		} else if err == nil && test.isErr {
			t.Errorf("expected error, no error received for %s", test.name)
			continue
		}
	}

}

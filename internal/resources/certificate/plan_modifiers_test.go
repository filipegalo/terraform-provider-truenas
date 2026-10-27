// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package certificate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestStorePrivatekeyChanged verifies store_privatekey forces replacement only
// when its effective value changes. The regression: state written before the
// attribute existed holds null, the default plans true, and the certificate
// (plus any ACME certificate chained to it via csr_id) was replaced on
// upgrade. (#49)
func TestStorePrivatekeyChanged(t *testing.T) {
	cases := []struct {
		name  string
		state types.Bool
		plan  types.Bool
		want  bool
	}{
		{"null prior, planned default true (upgrade)", types.BoolNull(), types.BoolValue(true), false},
		{"null prior, planned false", types.BoolNull(), types.BoolValue(false), true},
		{"true to false", types.BoolValue(true), types.BoolValue(false), true},
		{"false to true", types.BoolValue(false), types.BoolValue(true), true},
		{"unchanged true", types.BoolValue(true), types.BoolValue(true), false},
		{"unchanged false", types.BoolValue(false), types.BoolValue(false), false},
		{"unknown plan", types.BoolValue(false), types.BoolUnknown(), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := storePrivatekeyChanged(tc.state, tc.plan); got != tc.want {
				t.Errorf("storePrivatekeyChanged(%v, %v) = %v, want %v", tc.state, tc.plan, got, tc.want)
			}
		})
	}
}

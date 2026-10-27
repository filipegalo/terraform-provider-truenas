// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// replaceIfStorePrivatekeyChanged requires replacement only when the effective
// store_privatekey value changes. State written before the attribute existed
// holds null; the default plans true, and a plain RequiresReplace treated that
// null -> true as a change, replacing every certificate on upgrade (and with a
// CSR, the ACME certificate chained to it through csr_id). A null prior value
// reads as true, so planning true records it in place. (#49)
func replaceIfStorePrivatekeyChanged() planmodifier.Bool {
	return boolplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.BoolRequest, resp *boolplanmodifier.RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = storePrivatekeyChanged(req.StateValue, req.PlanValue)
		},
		"Replaced only when the effective value changes (a null prior value reads as true and records in place).",
		"Replaced only when the effective value changes (a null prior value reads as true and records in place).",
	)
}

// storePrivatekeyChanged reports whether the effective store_privatekey value
// differs between prior state and plan. An unknown plan value is not a change.
func storePrivatekeyChanged(state, plan types.Bool) bool {
	if plan.IsUnknown() {
		return false
	}
	return storesPrivatekey(state) != storesPrivatekey(plan)
}

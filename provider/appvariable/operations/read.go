package operations

import (
	"context"
	"errors"

	"github.com/eaze/terraform-provider-expo/internal/client"
	"github.com/eaze/terraform-provider-expo/internal/eas"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*client.EASClient)

	input := eas.GetAppVariableData{
		Id:    d.Id(),
		AppId: d.Get("app_id").(string),
	}

	data, err := client.AppVariable.Get(input)
	if errors.Is(err, eas.ErrAppVariableNotFound) {
		// Deleted outside Terraform: drop it from state so the next plan
		// recreates it instead of failing.
		tflog.Warn(ctx, "expo_app_variable not found in EAS, removing from state", map[string]any{
			"app_id": input.AppId,
			"id":     d.Id(),
		})
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(data.Id)

	if err := d.Set("id", data.Id); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", data.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	// EAS never returns the value of a SECRET variable; keep what is in
	// state so the plan does not show a perpetual change.
	if data.Visibility != "SECRET" {
		if err := d.Set("value", data.Value); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if err := d.Set("visibility", data.Visibility); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("environments", data.Environments); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("app_id", input.AppId); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

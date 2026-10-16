package operations

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parseImportId(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected import id %q: expected <app_id>/<variable_id>", id)
	}
	return parts[0], parts[1], nil
}

func Import(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	appId, variableId, err := parseImportId(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("app_id", appId); err != nil {
		return nil, err
	}
	if err := d.Set("id", variableId); err != nil {
		return nil, err
	}
	d.SetId(variableId)

	if diags := Read(ctx, d, m); diags.HasError() {
		return nil, fmt.Errorf("reading imported variable %s: %s", variableId, diags[0].Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("environment variable %s not found on app %s", variableId, appId)
	}

	return []*schema.ResourceData{d}, nil
}

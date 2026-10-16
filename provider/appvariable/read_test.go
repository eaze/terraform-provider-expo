package appvariable

import (
	"context"
	"errors"
	"testing"

	"github.com/eaze/terraform-provider-expo/internal/client"
	"github.com/eaze/terraform-provider-expo/internal/eas"
	"github.com/eaze/terraform-provider-expo/provider/appvariable/operations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

type fakeAppVariableService struct {
	data *eas.AppVariableData
	err  error
}

func (f *fakeAppVariableService) Get(eas.GetAppVariableData) (*eas.AppVariableData, error) {
	return f.data, f.err
}
func (f *fakeAppVariableService) GetByName(eas.GetByNameAppVariableData) (*eas.AppVariableData, error) {
	return f.data, f.err
}
func (f *fakeAppVariableService) Create(eas.CreateAppVariableData) (*eas.AppVariableData, error) {
	return f.data, f.err
}
func (f *fakeAppVariableService) Update(eas.UpdateAppVariableData) (*eas.AppVariableData, error) {
	return f.data, f.err
}
func (f *fakeAppVariableService) Delete(string) (*any, error) { return nil, f.err }

func newTestClient(fake *fakeAppVariableService) *client.EASClient {
	return &client.EASClient{EASClient: eas.EASClient{AppVariable: fake}}
}

func newTestState(t *testing.T, visibility string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]any{
		"app_id":       "app-1",
		"name":         "MY_VAR",
		"value":        "from-state",
		"visibility":   visibility,
		"environments": []any{"production"},
	})
	d.SetId("var-1")
	return d
}

func TestReadClearsStateWhenVariableIsGone(t *testing.T) {
	d := newTestState(t, "PUBLIC")
	fake := &fakeAppVariableService{err: eas.ErrAppVariableNotFound}

	diags := operations.Read(context.Background(), d, newTestClient(fake))

	assert.False(t, diags.HasError(), "a deleted remote object is a plan-time recreate, not an error")
	assert.Equal(t, "", d.Id())
}

func TestReadReportsOtherErrors(t *testing.T) {
	d := newTestState(t, "PUBLIC")
	fake := &fakeAppVariableService{err: errors.New("boom")}

	diags := operations.Read(context.Background(), d, newTestClient(fake))

	assert.True(t, diags.HasError())
	assert.Equal(t, "var-1", d.Id())
}

func TestReadKeepsConfiguredValueForSecretVariable(t *testing.T) {
	d := newTestState(t, "SECRET")
	fake := &fakeAppVariableService{data: &eas.AppVariableData{
		Id: "var-1", Name: "MY_VAR", Value: "", Visibility: "SECRET", Environments: []string{"production"},
	}}

	diags := operations.Read(context.Background(), d, newTestClient(fake))

	assert.False(t, diags.HasError())
	assert.Equal(t, "from-state", d.Get("value"), "EAS never returns SECRET values; overwriting would cause a perpetual diff")
}

func TestReadSetsFieldsFromApi(t *testing.T) {
	d := newTestState(t, "PUBLIC")
	fake := &fakeAppVariableService{data: &eas.AppVariableData{
		Id: "var-1", Name: "MY_VAR", Value: "remote", Visibility: "SENSITIVE", Environments: []string{"preview", "development"},
	}}

	diags := operations.Read(context.Background(), d, newTestClient(fake))

	assert.False(t, diags.HasError())
	assert.Equal(t, "remote", d.Get("value"))
	assert.Equal(t, "SENSITIVE", d.Get("visibility"))
	assert.Equal(t, "app-1", d.Get("app_id"))
	assert.ElementsMatch(t, []any{"preview", "development"}, d.Get("environments").(*schema.Set).List())
}

func TestReadReportsMissingAppAsErrorAndKeepsState(t *testing.T) {
	d := newTestState(t, "PUBLIC")
	fake := &fakeAppVariableService{err: eas.ErrAppNotFound}

	diags := operations.Read(context.Background(), d, newTestClient(fake))

	assert.True(t, diags.HasError(), "a missing or inaccessible app is not a deleted variable")
	assert.Equal(t, "var-1", d.Id())
}

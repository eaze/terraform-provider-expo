package operations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImportIdSplitsAppIdAndVariableId(t *testing.T) {
	appId, variableId, err := parseImportId("app-123/var-456")

	assert.NoError(t, err)
	assert.Equal(t, "app-123", appId)
	assert.Equal(t, "var-456", variableId)
}

func TestParseImportIdRejectsMissingSeparator(t *testing.T) {
	_, _, err := parseImportId("var-456")

	assert.ErrorContains(t, err, "<app_id>/<variable_id>")
}

func TestParseImportIdRejectsEmptyParts(t *testing.T) {
	for _, id := range []string{"/var-456", "app-123/", "/"} {
		_, _, err := parseImportId(id)
		assert.Error(t, err, id)
	}
}

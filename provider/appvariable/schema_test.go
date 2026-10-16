package appvariable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueIsSensitive(t *testing.T) {
	assert.True(t, Resource().Schema["value"].Sensitive, "value must be redacted in plans")
}

func TestAppIdForcesNew(t *testing.T) {
	assert.True(t, Resource().Schema["app_id"].ForceNew, "the update mutation has no appId, so a new app means a new variable")
}

func TestImporterIsRegistered(t *testing.T) {
	assert.NotNil(t, Resource().Importer)
}

package graphql

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
)

func newTestGraphQL(t *testing.T, status int, body string) GraphQL {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return &graphQL{
		client: graphql.NewClient(server.URL, graphql.WithHTTPClient(newHTTPClient())),
		token:  "test-token",
	}
}

type probeResponse struct {
	AppByAppId *struct{} `json:"appByAppId"`
}

func TestQueryRejectsNon2xxJsonBodies(t *testing.T) {
	graphQL := newTestGraphQL(t, http.StatusServiceUnavailable, `{"message":"upstream unavailable"}`)

	var response probeResponse
	err := graphQL.Query("query { appByAppId { id } }", nil, &response)

	assert.ErrorContains(t, err, "HTTP 503")
	assert.ErrorContains(t, err, "upstream unavailable")
	assert.Nil(t, response.AppByAppId)
}

func TestQueryStillSurfacesGraphQLErrors(t *testing.T) {
	graphQL := newTestGraphQL(t, http.StatusOK, `{"data":null,"errors":[{"message":"Unauthorized"}]}`)

	var response probeResponse
	err := graphQL.Query("query { appByAppId { id } }", nil, &response)

	assert.ErrorContains(t, err, "Unauthorized")
}

func TestQueryDecodesSuccessfulResponses(t *testing.T) {
	graphQL := newTestGraphQL(t, http.StatusOK, `{"data":{"appByAppId":{}}}`)

	var response probeResponse
	err := graphQL.Query("query { appByAppId { id } }", nil, &response)

	assert.NoError(t, err)
	assert.NotNil(t, response.AppByAppId)
}

package graphql

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eaze/terraform-provider-expo/internal/eas/internal/graphql/machinebox"
	"github.com/machinebox/graphql"
)

type VariableValue interface {
	~int | ~string | ~[]string
}

type GraphQL interface {
	Query(query string, variables map[string]any, response any) error
}

type graphQL struct {
	client machinebox.MachineboxGraphqlClient
	token  string
}

var _ GraphQL = (*graphQL)(nil)

func NewGraphQL(token string) GraphQL {
	client := graphql.NewClient("https://api.expo.dev/graphql", graphql.WithHTTPClient(newHTTPClient()))
	return &graphQL{
		client: client,
		token:  token,
	}
}

func newHTTPClient() *http.Client {
	return &http.Client{Transport: statusCheckingTransport{next: http.DefaultTransport}}
}

type statusCheckingTransport struct {
	next http.RoundTripper
}

const errorBodyLimit = 512

func (t statusCheckingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 400 {
		return res, nil
	}

	body, _ := io.ReadAll(io.LimitReader(res.Body, errorBodyLimit))
	_ = res.Body.Close()
	return nil, fmt.Errorf("expo api returned HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
}

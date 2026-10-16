package appvariable

import (
	"errors"
	"testing"

	"github.com/eaze/terraform-provider-expo/internal/eas/internal/graphql"
	"github.com/eaze/terraform-provider-expo/internal/eas/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var expectedAppId = "test-app-id"

var expectedData = &Data{
	Id:           "test-id",
	Name:         "test-name",
	Value:        "test-value",
	Visibility:   "PUBLIC",
	Environments: []string{"development"},
}

var mockResponse = getResponse{AppByAppId: &appById{Data: []Data{*expectedData}}}

func TestGetByName(t *testing.T) {
	getData := &GetByNameData{
		AppId: expectedAppId,
		Name:  expectedData.Name,
	}
	expectedVariables := map[string]any{"appId": expectedAppId}

	config := utils.TestConfig[GetByNameData, Data, getResponse, Service]{
		NewServiceFunction: NewService,
		FunctionUnderTest:  "GetByName",
		Input:              getData,
		MockResponse:       mockResponse,
		ExpectedQuery:      getQuery,
		ExpectedVariables:  expectedVariables,
		ExpectedData:       expectedData,
	}
	utils.Test(t, config)
}

func TestGetById(t *testing.T) {
	getData := &GetData{
		AppId: expectedAppId,
		Id:    expectedData.Id,
	}
	expectedVariables := map[string]any{"appId": expectedAppId}

	config := utils.TestConfig[GetData, Data, getResponse, Service]{
		NewServiceFunction: NewService,
		FunctionUnderTest:  "Get",
		Input:              getData,
		MockResponse:       mockResponse,
		ExpectedQuery:      getQuery,
		ExpectedVariables:  expectedVariables,
		ExpectedData:       expectedData,
	}
	utils.Test(t, config)
}

func TestGetByIdReturnsErrNotFoundWhenMissing(t *testing.T) {
	service := NewService(graphql.NewGraphQLMock(mockResponse))

	data, err := service.Get(GetData{AppId: expectedAppId, Id: "missing-id"})

	assert.Nil(t, data)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGetByNameReturnsErrNotFoundWhenMissing(t *testing.T) {
	service := NewService(graphql.NewGraphQLMock(mockResponse))

	data, err := service.GetByName(GetByNameData{AppId: expectedAppId, Name: "missing-name"})

	assert.Nil(t, data)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGetReturnsErrAppNotFoundWhenAppIsNull(t *testing.T) {
	service := NewService(graphql.NewGraphQLMock(getResponse{AppByAppId: nil}))

	data, err := service.Get(GetData{AppId: expectedAppId, Id: expectedData.Id})

	assert.Nil(t, data)
	assert.ErrorIs(t, err, ErrAppNotFound)
	assert.NotErrorIs(t, err, ErrNotFound)
}

func TestGetByNameReturnsErrAppNotFoundWhenAppIsNull(t *testing.T) {
	service := NewService(graphql.NewGraphQLMock(getResponse{AppByAppId: nil}))

	data, err := service.GetByName(GetByNameData{AppId: expectedAppId, Name: expectedData.Name})

	assert.Nil(t, data)
	assert.ErrorIs(t, err, ErrAppNotFound)
	assert.NotErrorIs(t, err, ErrNotFound)
}

func TestGetPropagatesQueryErrorsWithoutErrNotFound(t *testing.T) {
	graphQLMock := &graphql.GraphQLMock{}
	graphQLMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("boom"))
	service := NewService(graphQLMock)

	data, err := service.Get(GetData{AppId: expectedAppId, Id: expectedData.Id})

	assert.Nil(t, data)
	assert.EqualError(t, err, "boom")
	assert.NotErrorIs(t, err, ErrNotFound)
	assert.NotErrorIs(t, err, ErrAppNotFound)
}

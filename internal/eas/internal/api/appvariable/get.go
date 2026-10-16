package appvariable

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound    = errors.New("environment variable not found")
	ErrAppNotFound = errors.New("app not found or not accessible with this token")
)

type appById struct {
	Data []Data `json:"environmentVariablesIncludingSensitive"`
}

type getResponse struct {
	AppByAppId *appById `json:"appByAppId"`
}

const getQuery = `
	query ($appId: String!) {
		appByAppId(appId: $appId) {
			environmentVariablesIncludingSensitive {
				environments
				id
				name
				value
				visibility
			}
		}
	}`

func (service *service) getAppEnvironmentVariables(appId string) ([]Data, error) {
	variables := map[string]any{"appId": appId}

	var response getResponse

	if err := service.graphql.Query(getQuery, variables, &response); err != nil {
		return nil, err
	}
	if response.AppByAppId == nil {
		return nil, fmt.Errorf("%w: appId %s", ErrAppNotFound, appId)
	}
	return response.AppByAppId.Data, nil
}

// Retrieves an App Environment Variable from EAS by it's name and appId
func (service *service) GetByName(getByName GetByNameData) (*Data, error) {
	data, err := service.getAppEnvironmentVariables(getByName.AppId)
	if err != nil {
		return nil, err
	}
	return findByName(data, getByName.Name)
}

// Retrieves an App Environment Variable from EAS by it' id and appId
func (service *service) Get(getData GetData) (*Data, error) {
	data, err := service.getAppEnvironmentVariables(getData.AppId)
	if err != nil {
		return nil, err
	}

	return findById(data, getData.Id)
}

func findByName(variables []Data, name string) (*Data, error) {
	for _, variable := range variables {
		if variable.Name == name {
			return &variable, nil
		}
	}
	return nil, fmt.Errorf("%w: name %s", ErrNotFound, name)
}

func findById(variables []Data, id string) (*Data, error) {
	for _, variable := range variables {
		if variable.Id == id {
			return &variable, nil
		}
	}
	return nil, fmt.Errorf("%w: id %s", ErrNotFound, id)
}

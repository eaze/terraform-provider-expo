package android

import (
	"github.com/eaze/terraform-provider-expo/internal/eas/internal/api/android/appbuildcredentials"
	"github.com/eaze/terraform-provider-expo/internal/eas/internal/api/android/appcredentials"
	"github.com/eaze/terraform-provider-expo/internal/eas/internal/api/android/fcmkey"
	"github.com/eaze/terraform-provider-expo/internal/eas/internal/api/android/googleserviceaccountkey"
	"github.com/eaze/terraform-provider-expo/internal/eas/internal/api/android/keystore"
	"github.com/eaze/terraform-provider-expo/internal/eas/internal/graphql"
)

type Service struct {
	GoogleServiceAccountKey googleserviceaccountkey.Service
	AppCredentials          appcredentials.Service
	AppBuildCredentials     appbuildcredentials.Service
	FCMKey                  fcmkey.Service
	Keystore                keystore.Service
}

func NewService(graphQL graphql.GraphQL) Service {
	return Service{
		GoogleServiceAccountKey: googleserviceaccountkey.NewService(graphQL),
		AppCredentials:          appcredentials.NewService(graphQL),
		AppBuildCredentials:     appbuildcredentials.NewService(graphQL),
		FCMKey:                  fcmkey.NewService(graphQL),
		Keystore:                keystore.NewService(graphQL),
	}
}

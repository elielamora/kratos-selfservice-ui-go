package apiclient

import (
	"net/url"
	"os"
	"sync"

	kratos "github.com/ory/kratos-client-go"
)

var (
	publicClientOnce     sync.Once
	publicClientInstance *kratos.APIClient
)

func InitPublicClient(url url.URL) *kratos.APIClient {

	publicClientOnce.Do(func() {
		cfg := &kratos.Configuration{
			Servers: []kratos.ServerConfiguration{
				{
					URL: url.String(),
				},
			},
			Debug: isDebug(),
		}
		publicClientInstance = kratos.NewAPIClient(cfg)
	})

	return publicClientInstance
}

func PublicClient() *kratos.APIClient {
	return publicClientInstance
}

var (
	adminClientOnce     sync.Once
	adminClientInstance *kratos.APIClient
)

func InitAdminClient(url url.URL) *kratos.APIClient {
	adminClientOnce.Do(func() {
		cfg := &kratos.Configuration{
			Servers: []kratos.ServerConfiguration{
				{
					URL: url.String(),
				},
			},
			Debug: isDebug(),
		}
		cfg.Scheme = url.Scheme
		adminClientInstance = kratos.NewAPIClient(cfg)
	})

	return adminClientInstance
}

func isDebug() bool {
	return len(os.Getenv("DEBUG")) > 0
}

func AdminClient() *kratos.APIClient {
	return adminClientInstance
}

package resources

import (
	"github.com/nitrictech/go-sdk/nitric"
)

func GetPublicApi() nitric.Api {
	api, _ := nitric.NewApi("main")
	return api
}

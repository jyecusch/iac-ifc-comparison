package resources

import (
	"github.com/nitrictech/go-sdk/nitric"
)

func GetTopic() nitric.SubscribableTopic {
	return nitric.NewTopic("events")
}

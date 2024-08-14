package eventbridge

import (
	"fmt"
	"strings"
)

// SourceId - Generate a source id for an eventbridge rule
func SourceId(stackId string, topicName string) string {
	return fmt.Sprintf("nitric.%s.topic.%s", stackId, topicName)
}

// BusName - Generate a bus name for an eventbridge bus
func BusName(stackId string, topicName string) string {
	return fmt.Sprintf("nitric-%s-%s", stackId, topicName)
}

func GetTopicNameFromSourceId(sourceId string) (string, error) {
	if sourceId == "" {
		return "", fmt.Errorf("sourceId is empty")
	}

	if !strings.HasPrefix(sourceId, "nitric.") {
		return "", fmt.Errorf("sourceId does not start with 'nitric.'")
	}

	sourceIdParts := strings.Split(sourceId, ".")

	if len(sourceIdParts) != 4 {
		return "", fmt.Errorf("invalid topic sourceId format")
	}

	return sourceIdParts[len(sourceIdParts)-1], nil
}

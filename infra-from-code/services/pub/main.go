package main

import (
	"context"
	"fmt"

	"github.com/nitrictech/go-sdk/handler"
	"github.com/nitrictech/go-sdk/nitric"
	"github.com/nitrictech/templates/go-starter/resources"
)

type RouteHandler = func(ctx *handler.HttpContext, next func(*handler.HttpContext) (*handler.HttpContext, error)) (*handler.HttpContext, error)

func triggerHandler(topic nitric.Topic) RouteHandler {
	return func(ctx *handler.HttpContext, next func(*handler.HttpContext) (*handler.HttpContext, error)) (*handler.HttpContext, error) {
		err := topic.Publish(context.TODO(), map[string]interface{}{
			"data": string(ctx.Request.Data()),
		})
		if err != nil {
			ctx.Response.Body = []byte("Error publishing message")
			ctx.Response.Status = 500
			return next(ctx)
		}

		ctx.Response.Body = []byte("Message published")
		return next(ctx)
	}
}

func main() {
	api := resources.GetPublicApi()
	topic, _ := resources.GetTopic().Allow(nitric.TopicPublish)

	api.Post("/trigger", triggerHandler(topic))

	if err := nitric.Run(); err != nil {
		fmt.Println(err)
	}
}

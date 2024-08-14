package main

import (
	"fmt"

	"github.com/nitrictech/go-sdk/handler"
	"github.com/nitrictech/go-sdk/nitric"
	"github.com/nitrictech/templates/go-starter/resources"
)

func handleMsg(ctx *handler.MessageContext, next handler.MessageHandler) (*handler.MessageContext, error) {
	fmt.Printf("Received message: %+v", ctx.Request.Message())

	return ctx, nil
}

func main() {
	topic := resources.GetTopic()

	topic.Subscribe(handleMsg)

	if err := nitric.Run(); err != nil {
		fmt.Println(err)
	}
}

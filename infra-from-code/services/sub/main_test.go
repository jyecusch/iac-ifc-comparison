package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nitrictech/go-sdk/handler"
	mock_handler "github.com/nitrictech/templates/go-starter/mocks/handler"
)

func noOpNext(ctx *handler.MessageContext) (*handler.MessageContext, error) {
	return ctx, nil
}

func TestHandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRequest := mock_handler.NewMockMessageRequest(ctrl)
	mockRequest.EXPECT().Message().Return(map[string]interface{}{"name": "test"})

	ctx := &handler.MessageContext{
		Request:  mockRequest,
		Response: &handler.MessageResponse{},
	}

	_, err := handleMsg(ctx, noOpNext)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

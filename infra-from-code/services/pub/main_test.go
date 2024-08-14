package main

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nitrictech/go-sdk/handler"
	mock_handler "github.com/nitrictech/templates/go-starter/mocks/handler"
	mock_nitric "github.com/nitrictech/templates/go-starter/mocks/nitric"
)

func noOpNext(ctx *handler.HttpContext) (*handler.HttpContext, error) {
	return ctx, nil
}

func TestTriggerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRequest := mock_handler.NewMockHttpRequest(ctrl)
	mockRequest.EXPECT().Data().Return([]byte(`{"name": "test"}`))

	mockTopic := mock_nitric.NewMockTopic(ctrl)
	mockTopic.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	handleRequest := triggerHandler(mockTopic)

	ctx := &handler.HttpContext{
		Request:  mockRequest,
		Response: &handler.HttpResponse{},
	}

	_, err := handleRequest(ctx, noOpNext)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(ctx.Response.Body) != "Message published" {
		t.Errorf("Expected 'Message published', got %s", string(ctx.Response.Body))
	}
}

func TestTriggerError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRequest := mock_handler.NewMockHttpRequest(ctrl)
	mockRequest.EXPECT().Data().Return([]byte(`{"name": "test"}`))

	mockTopic := mock_nitric.NewMockTopic(ctrl)
	mockTopic.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

	handleRequest := triggerHandler(mockTopic)

	ctx := &handler.HttpContext{
		Request:  mockRequest,
		Response: &handler.HttpResponse{},
	}

	ctx, err := handleRequest(ctx, noOpNext)

	if err != nil {
		t.Errorf("Expected http error, got %v", err)
	}

	if ctx.Response.Status != 500 {
		t.Errorf("Expected HTTP 500 status, got %d", ctx.Response.Status)
	}
}

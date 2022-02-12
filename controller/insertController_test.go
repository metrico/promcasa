package controllerv1

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/metrico/promcasa/service"
)

//for testing -> https://dev.to/koddr/go-fiber-by-examples-testing-the-application-1ldf
//run: go test ./...
func TestInsertController_PushStream(t *testing.T) {
	type fields struct {
		Controller    Controller
		InsertService *service.InsertService
	}
	type args struct {
		ctx *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &InsertController{
				Controller:    tt.fields.Controller,
				InsertService: tt.fields.InsertService,
			}
			if err := uc.PushStream(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("InsertController.PushStream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

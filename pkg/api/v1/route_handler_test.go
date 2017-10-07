package v1

import (
	"testing"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"

	"github.com/gin-gonic/gin"
)

func TestRouteHandler_Register(t *testing.T) {
	type fields struct {
		ingestionHandler ingestion.Handler
	}
	type args struct {
		router *gin.Engine
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &RouteHandler{
				ingestionHandler: tt.fields.ingestionHandler,
			}
			h.Register(tt.args.router)
		})
	}
}

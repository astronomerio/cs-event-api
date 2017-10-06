package v1

import (
	"testing"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/logger"
)

func TestRouteHandler_Register(t *testing.T) {
	type fields struct {
		ingestionHandler ingestion.IngestionHandler
		logger           logger.Logger
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
				logger:           tt.fields.logger,
			}
			h.Register(tt.args.router)
		})
	}
}

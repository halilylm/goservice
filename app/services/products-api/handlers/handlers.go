// Package handlers includes route mappings for debug endpoints.
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/halilylm/service/app/services/products-api/handlers/debug/checkgrp"
	"go.uber.org/zap"
)

// DebugRoutes registers debug applications routes for the service.
func DebugRoutes(app *fiber.App, build string, log *zap.SugaredLogger) {
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	app.Get("/debug/readiness", cgh.Readiness)
	app.Get("/debug/liveness", cgh.Liveness)
}

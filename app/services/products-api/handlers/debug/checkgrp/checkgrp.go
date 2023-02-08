// Package checkgrp maintains the group of handlers for health checking.
package checkgrp

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"os"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	Build string
	Log   *zap.SugaredLogger
}

func (h Handlers) Readiness(c *fiber.Ctx) error {
	data := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK

	h.Log.Infow("readiness", "statusCode", statusCode, "method", c.Method(), "path", c.Path(), "remoteaddr", c.IP())

	return c.Status(statusCode).JSON(data)
}

func (h Handlers) Liveness(c *fiber.Ctx) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status: "up",
		Build:  h.Build,
		Host:   host,
		Pod:    os.Getenv("KUBERNETES_PODNAME"),
		PodIP:  os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:   os.Getenv("KUBERNETES_NODENAME"),
	}

	statusCode := http.StatusOK

	h.Log.Infow("liveness", "statusCode", statusCode, "method", c.Method(), "path", c.Path(), "remoteaddr", c.IP())

	return c.Status(statusCode).JSON(data)
}

package routes

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"resumpro/internal/app/resumepro/server/handlers"
	"resumpro/internal/app/resumepro/server/middleware"
)

func GetRoutesList(logger *logrus.Logger, router *chi.Mux, lm *middleware.Log) *Routes {
	router.Use(lm.Handler)
	return &Routes{
		Routes: []*Route{
			{
				Pattern: "/",
				Method:  "GET",
				Handler: handlers.NewHello(logger).Handler,
			},
		},
		Router: router,
		Log:    logger,
	}
}

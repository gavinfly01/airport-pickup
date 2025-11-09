package http

import (
	"net/http"

	pkghttp "github.com/gavin/airport-pickup/pkg/http"
	"github.com/gin-gonic/gin"
)

// NewRouter wires all HTTP routes and returns an http.Handler (gin.Engine).
func NewRouter(orderApp OrderApp, settlementApp SettlementApp) http.Handler {
	r := gin.New()
	r.Use(pkghttp.CORS(), pkghttp.Logger(), pkghttp.Recovery())

	h := &Handler{orderApp: orderApp, settlementApp: settlementApp}

	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// Resources
	r.POST("/passengers", h.createPassenger)
	r.POST("/drivers", h.createDriver)
	r.POST("/pickup_requests", h.createPickupRequest)
	r.POST("/driver_offers", h.createDriverOffer)

	// bookings: GET list, POST complete (query id)
	r.GET("/bookings", h.listBookings)
	r.POST("/bookings", h.completeBooking)

	return r
}

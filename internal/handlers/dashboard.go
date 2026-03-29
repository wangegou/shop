package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ShowDashboard renders the admin dashboard page.
func ShowDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":      "Dashboard",
		"active_nav": "dashboard",
	})
}

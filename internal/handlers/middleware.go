package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthRequired is a middleware to ensure the user is authenticated.
// It redirects to the login page for HTML page requests and returns a JSON
// error for API requests.
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		// Check if the request is for an API endpoint or an HTML page
		if strings.HasPrefix(c.Request.URL.Path, "/admin/api") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.Redirect(http.StatusFound, "/admin/login")
		}
		c.Abort()
		return
	}

	c.Next()
}

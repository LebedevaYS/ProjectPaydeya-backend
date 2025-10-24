package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("userRole")

        if userRole != "admin" {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Access denied. Admin rights required",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
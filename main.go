package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Health check rout  
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"health": true,
		})
	})
	const PORT = ":5000"
	r.Run(PORT)
}

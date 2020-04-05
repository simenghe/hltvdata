package main

import (
	"hltvdata/scraper"
	"net/http"

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
	r.GET("/currentrankings", func(c *gin.Context) {
		teams := scraper.ScrapeHltvTeams()
		c.JSON(http.StatusOK, teams)
	})
	r.GET("/hltvtest", func(c *gin.Context) {
		teams := scraper.HltvTest()
		c.JSON(http.StatusOK, teams)
	})
	r.GET("/test", func(c *gin.Context) {
		go scraper.RankingTraverse()
		c.JSON(http.StatusOK, "loss")
	})
	r.GET("/testasync", func(c *gin.Context) {
		go scraper.RankingTraverseAsync()
		c.JSON(http.StatusOK, "asyncloss")
	})
	const PORT = ":5000"
	r.Run(PORT)
}

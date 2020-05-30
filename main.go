package main

import (
	"hltvdata/plot"
	"hltvdata/scraper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Routes defined here.
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
		c.JSON(http.StatusOK, scraper.RankingTraverseAsync())
	})
	// Adds a document, with new hltv urls
	r.GET("/updatehltvurls", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"duration": UpdateHLTVURLS().Seconds(),
		})
	})
	// Gets latest URL list
	r.GET("/gethltvurls", func(c *gin.Context) {
		urlObj, duration := GetHLTVURLS()
		c.JSON(http.StatusOK, gin.H{
			"duration":  duration,
			"urlList":   urlObj.URLS,
			"timestamp": urlObj.TimeStamp,
		})
	})
	r.GET("/updatehltvrankings", func(c *gin.Context) {

		c.JSON(http.StatusOK, UpdateHLTVRankings())
	})
	r.GET("/plot", func(c *gin.Context) {
		plot.TestPlot()
		c.JSON(http.StatusOK, gin.H{
			"Imagining": "Loss",
		})
	})
	const PORT = ":5000"
	r.Run(PORT)
}

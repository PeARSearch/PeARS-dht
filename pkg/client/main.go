package client

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/PeARSearch/cacophony-dht/pkg/peer"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Setup(cfg *peer.PeerConfig) *gin.Engine {
	r := gin.New()
	r.Use(CORSMiddleware())
	r.Use(recoveryLogger())
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK\n")
	})
	r.GET("/search/:word", func(c *gin.Context) {
		word := c.Param("word")
		val, err := cfg.GetData(c.Request.Context(), word)
		if err != nil {
			c.String(http.StatusFailedDependency, "NO OK\n")
		} else {
			c.String(http.StatusOK, fmt.Sprintf("%s\n", val))
		}
	})
	r.GET("/store/:word/:url", func(c *gin.Context) { // Get all jobs in a group
		word := c.Param("word")
		url := c.Param("url")
		err := cfg.PutData(c.Request.Context(), word, url)
		if err != nil {
			c.String(http.StatusFailedDependency, "NO OK\n")
		} else {
			c.String(http.StatusOK, "OK\n")

		}
	})

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+
			"Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept,"+
			" origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		c.Next()
	}
}

func recoveryLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"error":      err,
					"stacktrace": strings.Split(string(debug.Stack()), "\n"),
				}).Error("recovered panic")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

package core

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ExecuteApi(c *gin.Context) (interface{}, error) {
	param, err := ExtractParams(c)
	if err == nil {
		// result, ok := repository.Invoke("sql", param)
		return param, nil
	}
	return nil, err
}

func DBRRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//t := time.Now()

		if strings.Contains(strings.ToLower(c.Request.URL.Path), strings.ToLower("/api")) {
			result, err := ExecuteApi(c)
			if err != nil {
				c.JSON(200, gin.H{
					"statusCode":   500,
					"errorMessage": err.Error(),
					"data":         nil,
				})
			} else {
				c.JSON(200, gin.H{
					"statusCode":   200,
					"errorMessage": nil,
					"data":         result,
				})
			}
		} else {
			c.Next()
		}

		// latency := time.Since(t)
		// log.Print(latency)

		// // access the status we are sending
		// status := c.Writer.Status()
		// log.Println(status)
	}
}

func Start() {
	r := gin.Default()
	r.Use(DBRRequestHandler())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":  "Dbr",
			"time": time.Now(),
		})
	})

	r.GET("/dev", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":  "Dbr dev pages",
			"time": time.Now(),
		})
	})

	r.Run(":8080")
}

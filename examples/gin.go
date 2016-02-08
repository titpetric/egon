package examples

import (
	"github.com/SlinSo/goTemplateBenchmark/model"
	"github.com/gin-gonic/contrib/cache"
	"github.com/gin-gonic/gin"
	"time"
)

func mainGin() {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		u := &model.User{}
		c.Status(200)
		c.Writer.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}

		SimpleTemplate(c.Writer, u)
	})

	// caching can still happen in middleware
	store := cache.NewInMemoryStore(time.Second)
	r.GET("/hello_cache", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		u := &model.User{}
		c.Status(200)
		c.Writer.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}

		SimpleTemplate(c.Writer, u)
	}))
	r.Run()
}

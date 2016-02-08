package examples

import (
	"github.com/SlinSo/goTemplateBenchmark/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EgoRenderer struct {
	ego func(http.ResponseWriter)
}

func (r EgoRenderer) Render(w http.ResponseWriter) error {
	w.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}
	 r.ego(w)
	 return nil
}

func mainWrapper() {
	r := gin.Default()

	r.GET("/param", func(c *gin.Context) {
		
		c.Render(200, EgoRenderer{ego:func(w http.ResponseWriter) {
			u := &model.User{}
			SimpleTemplate(w, u)
		}})

	})

	r.GET("/simple", func(c *gin.Context) {
		c.Render(200, EgoRenderer{ego: func(w http.ResponseWriter){
			NoVarTemplate(w)
		}})
	})
	r.Run()
}

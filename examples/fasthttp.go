package examples

import (
	"github.com/SlinSo/goTemplateBenchmark/model"
	"github.com/valyala/fasthttp"
)

func mainFasthttp() {
	h := func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(200)
		ctx.SetContentType("text/html; charset=utf-8")
		u := &model.User{}
		
		SimpleTemplate(ctx, u)
	}
	
	if err := fasthttp.ListenAndServe(":8080", h); err != nil {
		panic(err)
	}
}

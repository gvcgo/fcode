package models

import (
	"fmt"
	"net/http"
	"os"

	"github.com/moqsien/fcode/cnf"
	"github.com/moqsien/fcode/models/cf"
	"github.com/moqsien/fcode/models/fitten"
	"github.com/moqsien/fcode/models/openai"

	"github.com/gin-gonic/gin"
)

var (
	sig = make(chan struct{})
)

func Serve() {
	if !cnf.DefaultConf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gin.Logger())

	r.POST("/v1/completions", func(c *gin.Context) {
		c.Set(cnf.ModelCtxKey, cnf.DefaultModel)
		c.Set(cnf.ProxyCtxKey, cnf.DefaultConf.Proxy)

		switch cnf.DefaultModel.Type {
		case "open_ai":
			openai.HandleAll(c)
		case "fitten":
			fitten.HandleAll(c)
		case "cf":
			cf.HandleAll(c)
		case "cf2":
			cf.HandleCFgptOss(c)
		default:
			fmt.Println("unspported model type")
			os.Exit(1)
		}
	})

	r.POST("/v1/choose/model", func(ctx *gin.Context) {
		name := ctx.Query("name")

		found := false
		for _, mm := range cnf.DefaultConf.AIModels {
			if mm.Name == name {
				found = true
				cnf.DefaultModel = mm
				dm := &cnf.DefaultM{}
				dm.Save(cnf.DefaultModel.Name)
				break
			}
		}
		if !found {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"err_msg": "model not found",
			})
		}
	})

	r.POST("/v1/stop", func(ctx *gin.Context) {
		sig <- struct{}{}
	})

	go func() {
		r.Run(cnf.DefaultConf.GetPort())
	}()
	<-sig
}

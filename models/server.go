package models

import (
	"fcode/cnf"
	"fcode/models/fitten"
	"fcode/models/openai"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Serve() {
	r := gin.Default()

	// r.Use(gin.Logger())

	r.POST("/v1/completions", func(c *gin.Context) {
		c.Set(cnf.ModelCtxKey, cnf.DefaultModel)

		if strings.HasPrefix(cnf.DefaultModel.Name, "fitten_code") {
			fitten.HandleAll(c)
		} else {
			openai.HandleAll(c)
		}
	})

	r.POST("/v1/choose/model", func(ctx *gin.Context) {
		name := ctx.Query("model")

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

	r.Run(cnf.DefaultConf.GetPort())
}

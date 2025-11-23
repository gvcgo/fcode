package models

import (
	"fcode/cnf"
	"fcode/models/fitten"
	"strings"

	"github.com/gin-gonic/gin"
)

func Serve() {
	r := gin.Default()

	// r.Use(gin.Logger())

	r.POST("/*path", func(c *gin.Context) {
		if strings.HasPrefix(cnf.DefaultModel.Name, "fitten_code") {
			c.Set(cnf.ModelCtxKey, cnf.DefaultModel)
			fitten.HandleAll(c)
		}
	})

	r.Run(cnf.DefaultConf.GetPort())
}

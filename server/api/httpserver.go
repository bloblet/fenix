package api

import (
	"github.com/bloblet/fenix/server/databases"
	"github.com/bloblet/fenix/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPApi struct {
	authdb databases.AuthenticationManager
}

func (api *HTTPApi) ServeHTTP() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.GET("/verify_email", api.verifyEmail)
	utils.Log().Infof("Serving HTTP on %v", config.API.HTTPHost)
	r.Run(config.API.HTTPHost)
}

func (api *HTTPApi) verifyEmail(c *gin.Context) {
	userID, _ := c.GetQuery("u")
	ott, _ := c.GetQuery("t")

	if userID == "" || ott == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	ok := api.authdb.Verify(ott, userID)
	if ok {
		c.String(http.StatusOK, "ok")

	} else {
		c.String(http.StatusBadRequest, "bad")
	}

}

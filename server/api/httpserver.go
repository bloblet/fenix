package api

import (
	"github.com/bloblet/fenix/server/databases"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net"
	"net/http"
)

type HTTPApi struct {
	authdb databases.AuthenticationManager
}

func (api *HTTPApi) Serve(l *net.Listener) {
	r := gin.Default()
	r.GET("/verify_email", api.verifyEmail)

	r.RunListener(*l)
}

func (api *HTTPApi) verifyEmail(c *gin.Context) {
	userID, _ := c.GetQuery("id")
	ott, _ := c.GetQuery("t")

	if userID == "" || ott == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	ok := api.authdb.Verify(ott, userID)
	if ok {
		c.Render(200, render.String{Data: []interface{}{"OK"}})
	} else {
		c.Render(404, render.String{Data: []interface{}{"Error"}})
	}

}
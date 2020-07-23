package api

import (
	"encoding/json"
	"fenix/databases"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "github.com/gorilla/websocket"
)

type create struct {
	Username string `json:"username"`
}

// API Fenix API
type API struct {
	userDatabase databases.UserDatabase
}

func (api *API) badRequest(w http.ResponseWriter) {
	api.error(w, "ERR_INVALIDREQUEST", "You probably sent invalid JSON.", http.StatusBadRequest)
}
func (api *API) internalError(w http.ResponseWriter) {
	api.error(w, "ERR_INTERNALERROR", "Something bad happened!", http.StatusInternalServerError)
}
func (api *API) error(w http.ResponseWriter, errcode, msg string, statusCode int) {
	output, err := json.Marshal(map[string]interface{}{"s": false, "e": errcode, "m": msg})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(statusCode)
	w.Header().Set("content-type", "application/json")

	w.Write(output)
}
func (api *API) create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	email, password, ok := r.BasicAuth()
	
	if !ok {
		api.badRequest(w)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		api.internalError(w)
		return
	}

	// Unmarshal
	var msg create
	err = json.Unmarshal(b, &msg)
	if err != nil {
		api.badRequest(w)
		return
	}

	if len(msg.Username) > 32 {
		api.error(w, "ERR_USERNAMETOOLONG", "Your username is above 32 characters!", http.StatusBadRequest)
	}
	
	user, err := api.userDatabase.CreateUser(email, password, msg.Username)

	if (err == databases.UserExists{}) {
		api.error(w, "ERR_USEREXISTS", "That user already exists!", http.StatusForbidden)
		return
	} else if err != nil {
		api.internalError(w)
		return
	}

	output, err := json.Marshal(map[string]interface{}{"s": true, "d": user})
	if err != nil {
		api.internalError(w)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	w.Header().Set("location", user.ID)
	w.Write(output)
}

// Serve starts the API
func (api *API) Serve() {
	api.userDatabase = databases.UserDatabase{}
	router := httprouter.New()
	router.POST("/6.0.1/create", api.create)

	log.Fatal(http.ListenAndServe(":8080", router))
}

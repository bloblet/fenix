package api

import (
	"encoding/json"
	"fenix/databases"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "github.com/gorilla/websocket"
)

type create struct {
	username, email, password string
}

// API Fenix API
type API struct {
	userDatabase databases.UserDatabase
}

func (api *API) create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg create
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user, err := api.userDatabase.CreateUser(msg.email, msg.password, msg.username)

	if (err == databases.UserExists{}) {
		http.Error(w, err.Error(), 403)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	output, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)
}

// Serve starts the API
func (api *API) Serve() {
	api.userDatabase = databases.UserDatabase{}
	router := httprouter.New()
	router.POST("/create", api.create)

	log.Fatal(http.ListenAndServe(":8080", router))
}

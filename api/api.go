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

// Error base type for the error channel
type Error struct {
	Error error
	Fatal bool
}

type create struct {
	Username string `json:"username"`
}

// NewAPI makes a new API object with bells and whistles
func NewAPI(username, password string) API {
	api := API{}
	api.username = username
	api.password = password
	return api
}

// API Fenix API
type API struct {
	username     string
	password     string
	userDatabase databases.UserDatabase
	isTesting    bool
	err          chan Error
}

func (api *API) badRequest(w http.ResponseWriter) {
	api.error(w, "ERR_INVALIDREQUEST", "You probably sent invalid JSON.", http.StatusBadRequest)
}

func (api *API) internalError(w http.ResponseWriter) {
	api.error(w, "ERR_INTERNALERROR", "Something bad happened!", http.StatusInternalServerError)
}

func (api *API) maybeError(err error) {
	if api.isTesting {
		e := Error{}
		e.Error = err
		e.Fatal = true
		api.err <- e
	}
}

func (api *API) error(w http.ResponseWriter, errcode, msg string, statusCode int) {
	output, err := json.Marshal(map[string]interface{}{"s": false, "e": errcode, "m": msg})

	// This only runs if the json marshal fails.  So this is reaaaaly bad.
	if err != nil {
		go api.maybeError(err)

		w.WriteHeader(503)
		w.Header().Set("content-type", "application/json")

		w.Write([]byte("\"s\": false, \"e\": \"ERR_INTERNALERROR\", \"m\": \"Something very bad has happened.\"}"))
		return
	}

	w.WriteHeader(500)
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
func (api *API) Serve(err chan Error, isTesting bool) {
	api.err = err
	api.isTesting = isTesting
	api.userDatabase = databases.NewUserDatabase(api.username, api.password, isTesting)
	router := httprouter.New()
	router.POST("/6.0.1/create", api.create)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

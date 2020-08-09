package api

import (
	"encoding/json"
	"fenix/databases"
	"io/ioutil"
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
func NewAPI(username, password, prefix string) API {
	api := API{}
	api.username = username
	api.password = password
	api.prefix = prefix
	return api
}

// API Fenix API
type API struct {
	prefix       string
	username     string
	password     string
	UserDatabase databases.UserDatabase
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
	output, err := json.Marshal(map[string]interface{}{"s": "f", "e": errcode, "m": msg})

	// This only runs if the json marshal fails.  So this is reaaaaly bad.
	if err != nil {
		go api.maybeError(err)

		w.WriteHeader(503)
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set("content-type", "application/json")

	w.Write(output)
}

func (api *API) create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	email, password, ok := r.BasicAuth()

	// Fail if the BasicAuth is invalid
	if !ok {
		api.badRequest(w)
		return
	}

	// Attempt to read the body.
	// TODO Make sure this can't be exploited
	b, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	// Fail if there was an error while reading the body.  This could be the user's fault, or it could be ours.
	// Until I run into this error, this will continue to be a 500 Internal Error
	if err != nil {
		api.internalError(w)
		return
	}

	// Unmarshal the request body
	var msg create
	err = json.Unmarshal(b, &msg)

	// Fail if the JSON body is invalid.
	if err != nil {
		api.badRequest(w)
		return
	}

	// Fail if the user's username is over the max length (32)
	if len(msg.Username) > 32 {
		api.error(w, "ERR_USERNAMETOOLONG", "Your username is above 32 characters!", http.StatusBadRequest)
		return
	}

	// Create the user
	user, err := api.UserDatabase.CreateUser(email, password, msg.Username)

	// Fail if the user already is 
	if (err == databases.UserExists{}) {
		api.error(w, "ERR_USEREXISTS", "That user already exists!", http.StatusForbidden)
		return
	} 

	// Fail if there's no more discriminators
	if (err == databases.NoMoreDiscriminators{}) {
		api.error(w, "ERR_NOMOREDISCRIMINATORS", "Too many users have that username!", 409)
		return
	}
	
	// If there was a error connecting to the database, its our fault.
	if err != nil {
		api.internalError(w)
		return
	}

	// Marshal the user into JSON
	output, err := json.Marshal(map[string]interface{}{"s": true, "d": user})
	if err != nil {
		api.internalError(w)
		return
	}
	
	// Respond
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	w.Header().Set("location", user.ID)
	w.Write(output)
}

func (api *API) Serve(isTesting bool) {
	api.isTesting = isTesting
	api.UserDatabase = databases.NewUserDatabase(api.username, api.password, isTesting, api.prefix)
	router := httprouter.New()
	router.POST("/6.0.1/create", api.create)

	go http.ListenAndServe("0.0.0.0:8080", router)
}

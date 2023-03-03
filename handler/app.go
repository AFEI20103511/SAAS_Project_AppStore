package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appstore/model"
	"appstore/service"
)


func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
    fmt.Println("Received one upload request")
    decoder := json.NewDecoder(r.Body)
    var app model.App
    if err := decoder.Decode(&app); err != nil {
        panic(err)
    }

    // call service level function to handle this request
    service.SaveApp(&app)

    fmt.Fprintf(w, "Upload request received: %s\n", app.Description)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one search request")
    

    // get param from request
    title := r.URL.Query().Get("title")
    description := r.URL.Query().Get("description")
 
    // call service to handle request
    var apps []model.App
    var err error
    apps, err = service.SearchApps(title, description)
    if err != nil {
        http.Error(w, "Failed to read Apps from backend", http.StatusInternalServerError)
        return
    }
 
    // construct response
    // tell front end that return type is json
    w.Header().Set("Content-Type", "application/json")
    js, err := json.Marshal(apps)
    if err != nil {
        http.Error(w, "Failed to parse Apps into JSON format", http.StatusInternalServerError)
        return
    }
    // write
    w.Write(js)
 }
 
 func checkoutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one checkout request")
    w.Header().Set("Content-Type", "text/plain")
 
    appID := r.FormValue("appID")
    s, err := service.CheckoutApp(r.Header.Get("Origin"), appID)
    if err != nil {
        fmt.Println("Checkout failed.")
        w.Write([]byte(err.Error()))
        return
    }
 
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(s.URL))
 
    fmt.Println("Checkout process started!")
 }
 
 
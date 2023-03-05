package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appstore/model"
	"appstore/service"
	"strconv"

	"github.com/form3tech-oss/jwt-go"
	"github.com/pborman/uuid"
)


func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
   fmt.Println("Received one upload request")

    // get data from request
    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]


    app := model.App{
        Id:          uuid.New(),
        User:        username.(string),
        Title:       r.FormValue("title"),
        Description: r.FormValue("description"),
    }


    price, err := strconv.Atoi(r.FormValue("price"))
    fmt.Printf("%v,%T", price, price)
    if err != nil {
        fmt.Println(err)
    }
    app.Price = price


    file, _, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }


    err = service.SaveApp(&app, file)
    if err != nil {
        http.Error(w, "Failed to save app to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save app to backend %v\n", err)
        return
    }


    fmt.Println("App is saved successfully.")
    fmt.Fprintf(w, "App is saved successfully: %s\n", app.Description)
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
 
 
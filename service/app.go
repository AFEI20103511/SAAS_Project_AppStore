package service

import (
	"errors"

	"fmt"
	"reflect"

	"appstore/backend"
	"appstore/constants"
	"appstore/model"

	"mime/multipart"

	"github.com/olivere/elastic/v7"
	"github.com/stripe/stripe-go/v74"
)


func SearchApps(title string, description string) ([]model.App, error) {
    // title is empty => searchByDes
    // description is emtpy => searchByTitle
    if title == "" {
        return SearchAppsByDescription(description)
    }
    if description == "" {
        return SearchAppsByTitle(title)
    }


    query1 := elastic.NewMatchQuery("title", title)
    query2 := elastic.NewMatchQuery("description", description)
    query := elastic.NewBoolQuery().Must(query1, query2)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
    if err != nil {
        return nil, err
    }


    return getAppFromSearchResult(searchResult), nil
}


func SearchAppsByTitle(title string) ([]model.App, error) {
    query := elastic.NewMatchQuery("title", title)
    query.Operator("AND")
    if title == "" {
        query.ZeroTermsQuery("all")
    }
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
    if err != nil {
        return nil, err
    }

    return getAppFromSearchResult(searchResult), nil
}


func SearchAppsByDescription(description string) ([]model.App, error) {
    query := elastic.NewMatchQuery("description", description)
    query.Operator("AND")
    if description == "" {
        query.ZeroTermsQuery("all")
    }
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
    if err != nil {
        return nil, err
    }


    return getAppFromSearchResult(searchResult), nil
}

func SearchAppsByID(appID string) (*model.App, error) {
    query := elastic.NewTermQuery("id", appID)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
    if err != nil {
        return nil, err
    }
    results := getAppFromSearchResult(searchResult)
    if len(results) == 1 {
        return &results[0], nil
    }
    return nil, nil
}


func getAppFromSearchResult(searchResult *elastic.SearchResult) []model.App {
    var ptype model.App
    var apps []model.App
    for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
        p := item.(model.App)
        apps = append(apps, p)
    }
    return apps
}

func SaveApp(app *model.App, file multipart.File ) error {
    // call Stripe to get PriceID and ProductID
    productID, priceID, err := backend.CreateProductWithPrice(app.Title, app.Description, int64(app.Price*100))
    if err != nil {
        fmt.Printf("Failed to create Product and Price using Stripe SDK %v\n", err)
        return err
    }
    
    app.ProductID = productID
    app.PriceID = priceID
    
    // call GCS to store file and get URL
    medialink, err := backend.GCSBackend.SaveToGCS(file, app.Id)
    if err != nil {
        return err
    }
    app.Url = medialink


    // save everything to ES
    err = backend.ESBackend.SaveToES(app, constants.APP_INDEX, app.Id)
    if err != nil {
        fmt.Printf("Failed to save app to elastic search with app index %v\n", err)
        return err
    }
    fmt.Println("App is saved successfully to ES app index.")
 

    // for test
    //fmt.Printf("Product %s with price %s is successfully created", productID, priceID)
    return nil
}
 
func CheckoutApp(domain string, appID string) (*stripe.CheckoutSession, error) {
    app, err := SearchAppsByID(appID)
    if err != nil {
        return nil, err
    }
    if app == nil {
        return nil, errors.New("unable to find app in elasticsearch")
    }
    return backend.CreateCheckoutSession(domain, app.PriceID)
 }
 
 
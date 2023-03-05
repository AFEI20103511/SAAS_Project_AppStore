package service

import (
    "fmt"
    "reflect"

    "appstore/backend"
    "appstore/constants"
    "appstore/model"

    "github.com/olivere/elastic/v7"
)

func CheckUser(username, password string) (bool, error) {
	// 1. Read ES using username, compare input password and password from ES
	// 2. Read es using username, 
	query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("username", username))
    query.Must(elastic.NewTermQuery("password", password))
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    var utype model.User
    for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
        u := item.(model.User)
        if u.Password == password {
            fmt.Printf("Login as %s\n", username)
            return true, nil
        }
    }
    return false, nil

}

func AddUser(user *model.User) (bool, error) {
	// Check whether user singed up before
	query := elastic.NewTermQuery("username", user.Username)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    if searchResult.TotalHits() > 0 {
        return false, nil
    }

    err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
    if err != nil {
        return false, err
    }
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil

}
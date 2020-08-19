package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Item struct {
	Id     int    `json:"Id"`
	Item   string `json:"Item"`
	Amount int    `json:"Amount"`
	Price  string `json:"Price"`
}

type ErrorMessage struct {
	Message string `json:"Message"`
}

var Items []Item

//GET
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hint: getAllItems woked.....")
	json.NewEncoder(w).Encode(Items)

}

//GetItemWithId
func GetItemWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	find := false
	str, _ := strconv.Atoi(vars["id"])
	for _, item := range Items {
		if item.Id == str {
			find = true
			json.NewEncoder(w).Encode(item)
		}
	}
	if !find {
		w.WriteHeader(http.StatusNotFound)
		var erM = ErrorMessage{Message: "Error : No one item exist "}
		json.NewEncoder(w).Encode(erM)
	}
}

//PostNewItem...
func PostNewItem(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var item Item
	json.Unmarshal(reqBody, &item)
	w.WriteHeader(http.StatusCreated)
	Items = append(Items, item)
	json.NewEncoder(w).Encode(item)
}

//DeleteItemWithId ...
func DeleteItemWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	find := false

	for i, item := range Items {
		if item.Id == id {
			find = true
			w.WriteHeader(http.StatusAccepted)
			Items = append(Items[:i], Items[i+1:]...)
		}
	}
	if !find {
		w.WriteHeader(http.StatusNotFound)
		var erM = ErrorMessage{Message: "Error"}
		json.NewEncoder(w).Encode(erM)
	}

}

//PutExistsItem ....
func PutExistsItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idKey, _ := strconv.Atoi(vars["id"])
	finded := false

	for i, item := range Items {
		if item.Id == idKey {
			finded = true
			reqBody, _ := ioutil.ReadAll(r.Body)
			w.WriteHeader(http.StatusAccepted)
			json.Unmarshal(reqBody, &Items[i])
		}
	}

	if !finded {
		w.WriteHeader(http.StatusNotFound)
		var erM = ErrorMessage{Message: "Error 404"}
		json.NewEncoder(w).Encode(erM)
	}

}

func main() {

	fmt.Println("REST API worked....")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/items", GetAllItems).Methods("GET")
	myRouter.HandleFunc("/item/{id}", GetItemWithId).Methods("GET")
	myRouter.HandleFunc("/item", PostNewItem).Methods("POST")
	myRouter.HandleFunc("/item/{id}", DeleteItemWithId).Methods("DELETE")
	myRouter.HandleFunc("/item/{id}", PutExistsItem).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

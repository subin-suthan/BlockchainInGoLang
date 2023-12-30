package main

import (
	"crypto/md5"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Block struct {
	Pos       int
	Data      ProductCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

type Product struct {
	Id      string `json:"id"`
	Desc    string `json:"desc"`
	Comapny string `json:"company"`
	DoM     string `json:"date_manufacturing"`
}

type ProductCheckout struct {
	ProductId    string `json:"product_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Blockchain struct {
	blocks []*Block
}

var Blockchain *Blockchain

func newProduct(w http.ResponseWriter, r *http.Request) {

	var prod Product
	if err := json.NewDecoder(r.Body).Decode((&prod)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("coould not create:%v", err)
		w.Write([]byte("could not create new Product"))
		return
	}

	h := md5.New()
	//io.WriteString(h,)

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newProduct).Methods("POST")

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))

}

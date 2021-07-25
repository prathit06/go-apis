package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-microservices/product-api/data"
	"github.com/gorilla/mux"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l: l}
}

func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall", http.StatusInternalServerError)
	}
}

func (p *Product) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall", http.StatusBadRequest)
	}

	p.l.Printf("prod : %#v", prod)
	data.AddProduct(prod)
}

func (p *Product) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id to int", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT/UPDATE for Product id :", id)
	prod := &data.Product{}
	err = prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not found", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, "Product Not found", http.StatusInternalServerError)
		return
	}
	return
}

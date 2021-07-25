package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-microservices/product-api/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l: l}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProducts(rw, r)
		return
	}
	if r.Method == http.MethodPut {

		pid := r.URL.Path
		rx := regexp.MustCompile(`/([0-9]+)`)
		g := rx.FindAllStringSubmatch(pid, -1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadGateway)
		}
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadGateway)
		}
		idString := g[0][1]
		idInt, _ := strconv.Atoi(idString)
		p.l.Println("Got ID", idInt)

		p.updateProducts(idInt, rw, r)
		return
	}
	// everything else
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall", http.StatusInternalServerError)
	}
}

func (p *Product) addProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall", http.StatusBadRequest)
	}

	p.l.Printf("prod : %#v", prod)
	data.AddProduct(prod)
}

func (p *Product) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT/UPDATE Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
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

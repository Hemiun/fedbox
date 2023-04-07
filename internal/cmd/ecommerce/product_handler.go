package ecommerce

import (
	"encoding/json"
	"io"
	"net/http"

	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/product"
	"github.com/go-chi/chi/v5"
)

type ProductService interface {
	CreateProduct(caller vocab.Actor, token string, p product.Product) (vocab.Item, error)
	GetProduct(caller vocab.Actor, token string, productID string) (product.Product, error)
	GetProducts(caller vocab.Actor, token string) ([]product.Product, error)
}

// postProductHandler handles product creation POST requests
func postProductHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(body) == 0 {
		err = errors.NewBadRequest(err, "not empty body expected")
		logger.Errorf("not empty body expected", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	//Parsing a product from request
	var dto product.Product
	err = json.Unmarshal(body, &dto)
	logger.Infof("ProductHandler. Product dto parsed. Product.Name=%s", dto.Name)

	if err != nil {
		err = errors.NewBadRequest(err, "can't process request body")
		logger.Errorf("can't process request body", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	//Getting oAuth token
	token := r.Header.Get("Authorization")

	//Getting current actor
	callerActor, ok := r.Context().Value(common.AuthActorKey{}).(vocab.Actor)
	if !ok {
		err = errors.NewBadRequest(err, "can't get actor from context")
		logger.Errorf("can't get actor from context", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	logger.Infof("ProductHandler. Current actor found. Actor.ID=%s", callerActor.ID)

	it, err := productService.CreateProduct(callerActor, token, dto)
	if err != nil {
		err = errors.NewBadRequest(err, "product creation error")
		logger.Errorf("product creation error", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	//Serializing result product dto to json
	var resultProductJson []byte
	resultProductJson, err = json.Marshal(it)
	if err != nil {
		err = errors.NewBadRequest(err, "result product dto serialization error")
		logger.Errorf("result product dto serialization error", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	//Writing response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", it.GetLink().String())
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resultProductJson)
	if err != nil {
		logger.Errorf("can't write response", err)
	}
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	//Getting oAuth token
	token := r.Header.Get("Authorization")

	//Getting current actor
	callerActor, ok := r.Context().Value(common.AuthActorKey{}).(vocab.Actor)
	if !ok {
		err = errors.NewBadRequest(err, "can't get actor from context")
		logger.Errorf("can't get actor from context", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	logger.Infof("ProductHandler. Current actor found. Actor.ID=%s", callerActor.ID)

	//Trying to find a product
	var products []product.Product
	products, err = productService.GetProducts(callerActor, token)
	if err != nil {
		err = errors.NewBadRequest(err, "product searching error")
		logger.Errorf("product searching error", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	//Writing response depending on searching result
	logger.Infof("ProductHandler. Product found. Products count = %s", len(products))
	var productData []byte
	productData, err = json.Marshal(products)
	if err != nil {
		err = errors.NewBadRequest(err, "product data marshaling error")
		logger.Errorf("product data marshaling error", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(productData)
	if err != nil {
		logger.Errorf("can't write response", err)
	}
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	productID := chi.URLParam(r, "productID")
	if productID == "" {
		err = errors.NewBadRequest(err, "productID not passed")
		logger.Errorf("productID not passed", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	//Getting oAuth token
	token := r.Header.Get("Authorization")

	//Getting current actor
	callerActor, ok := r.Context().Value(common.AuthActorKey{}).(vocab.Actor)
	if !ok {
		err = errors.NewBadRequest(err, "can't get actor from context")
		logger.Errorf("can't get actor from context", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	logger.Infof("ProductHandler. Current actor found. Actor.ID=%s", callerActor.ID)

	//Trying to find a product
	var p product.Product
	p, err = productService.GetProduct(callerActor, token, productID)
	if err != nil {
		err = errors.NewBadRequest(err, "product searching error")
		logger.Errorf("product searching error", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	//Writing response depending on searching result
	if p.Id == "" {
		logger.Infof("ProductHandler. Product not found.")
		w.WriteHeader(http.StatusNotFound)
	} else {
		logger.Infof("ProductHandler. Product found. Product.Name=%s", p.Name)
		var productData []byte
		productData, err = json.Marshal(p)
		if err != nil {
			err = errors.NewBadRequest(err, "product data marshaling error")
			logger.Errorf("product data marshaling error", err)
			errors.HandleError(err).ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(productData)
		if err != nil {
			logger.Errorf("can't write response", err)
		}
	}
}

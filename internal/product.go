package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/luthfiaed/owapp-be/internal/data"
)

func (app *application) getProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	products, err := app.products.GetByName(name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, products, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProductById(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("id")
	idInt, err := strconv.Atoi(productId)
	if err != nil || idInt < 1 {
		app.serverErrorResponse(w, r, err)
		return
	}

	product, err := app.products.GetById(idInt)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, product, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createReview(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Review    string `json:"review"`
		ProductID string `json:"productId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	productIdInt, err := strconv.Atoi(input.ProductID)
	if err != nil || productIdInt < 1 || input.Review == "" {
		app.badRequestResponse(w, r, errors.New("product id or review cannot be empty"))
		return
	}

	user := app.contextGetUser(r)
	var data data.Review
	data.ProductID = int64(productIdInt)
	data.Review = input.Review
	data.Username = user.Username

	err = app.products.Insert(data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, map[string]string{"message": "ok"}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateReview(w http.ResponseWriter, r *http.Request) {
	var input data.Review

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ID < 1 || input.ProductID < 1 || input.Review == "" {
		app.badRequestResponse(w, r, errors.New("review id, product id, or review cannot be empty"))
		return
	}

	/*
		VULNERABILITY POINT: we only check if the user is logged in
		we do not check whether the user who requested to update this data
		is the same user who created it
	*/
	user := app.contextGetUser(r)
	input.Username = user.Username

	err = app.products.Update(input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]string{"message": "ok"}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}

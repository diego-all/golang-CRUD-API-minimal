package main

import (
	models "golang-crud-api-minimal/internal"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type categoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type categoryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (app *application) CreateCategory(w http.ResponseWriter, r *http.Request) {

	var categoryReq categoryRequest
	var payload jsonResponse

	err := app.readJSON(w, r, &categoryReq)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}

	var category = models.Category{
		Name:        categoryReq.Name,
		Description: categoryReq.Description,
	}

	_, err = app.models.Category.Insert(category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload = jsonResponse{
		Error:   false,
		Message: "Category successfully created",
		Data:    envelope{"category": category.Name},
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *application) GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	category, err := app.models.Category.GetOneById(categoryID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var categoryResponse = categoryResponse{
		Name:        category.Name,
		Description: category.Description,
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Category successfully obtained",
		Data:    envelope{"category": categoryResponse},
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *application) UpdateCategory(w http.ResponseWriter, r *http.Request) {

	var categoryReq categoryRequest
	var payload jsonResponse

	err := app.readJSON(w, r, &categoryReq)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}

	categoryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, err = app.models.Category.GetOneById(categoryID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var category = models.Category{
		Name:        categoryReq.Name,
		Description: categoryReq.Description,
		UpdatedAt:   time.Now(),
		Id:          categoryID,
	}

	_, err = app.models.Category.Update(category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload = jsonResponse{
		Error:   false,
		Message: "Category successfully updated",
		Data:    envelope{"category": category.Name},
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *application) AllCategories(w http.ResponseWriter, r *http.Request) {
	var categories models.Category
	all, err := categories.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Data successfully obtained",
		Data:    envelope{"categories": all},
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.Category.DeleteByID(categoryID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Category successfully deleted",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

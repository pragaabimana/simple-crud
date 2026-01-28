package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "simple-crud/docs"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// =======================
// MODEL
// =======================

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// =======================
// STORAGE (fake DB)
// =======================

var (
	categories = map[int]*Category{}
	autoID     = 1
)

// =======================
// HANDLER
// =======================

// GetCategories godoc
// @Summary Get all categories
// @Tags Category
// @Produce json
// @Success 200 {array} Category
// @Router /categories [get]
func GetCategories(w http.ResponseWriter, r *http.Request) {
	result := []*Category{}
	for _, v := range categories {
		result = append(result, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// CreateCategory godoc
// @Summary Create category
// @Tags Category
// @Accept json
// @Produce json
// @Param body body Category true "Category"
// @Success 201 {object} Category
// @Router /categories [post]
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.ID = autoID
	autoID++
	categories[input.ID] = &input

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

// GetCategory godoc
// @Summary Get category detail
// @Tags Category
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 404 {string} string
// @Router /categories/{id} [get]
func GetCategory(w http.ResponseWriter, r *http.Request) {
	id := parseID(r.URL.Path)
	category, ok := categories[id]
	if !ok {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory godoc
// @Summary Update category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body Category true "Category"
// @Success 200 {object} Category
// @Failure 404 {string} string
// @Router /categories/{id} [put]
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := parseID(r.URL.Path)
	category, ok := categories[id]
	if !ok {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	var input Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category.Name = input.Name
	category.Description = input.Description

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory godoc
// @Summary Delete category
// @Tags Category
// @Param id path int true "Category ID"
// @Success 204
// @Failure 404 {string} string
// @Router /categories/{id} [delete]
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := parseID(r.URL.Path)
	if _, ok := categories[id]; !ok {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	delete(categories, id)
	w.WriteHeader(http.StatusNoContent)
}

// =======================
// ROUTER HELPER
// =======================

func parseID(path string) int {
	parts := strings.Split(path, "/")
	id, _ := strconv.Atoi(parts[len(parts)-1])
	return id
}

// =======================
// MAIN
// =======================

// @title Simple Category API
// @version 1.0
// @description Simple CRUD using net/http + Swagger
// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running"))
	})
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetCategories(w, r)
		case http.MethodPost:
			CreateCategory(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetCategory(w, r)
		case http.MethodPut:
			UpdateCategory(w, r)
		case http.MethodDelete:
			DeleteCategory(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("server running at :", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"gitub.com/ndyabagye/go-rest-demo/pkg/recipes"
	"net/http"
)

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

func main() {
	//create the Store and Recipe handler
	store := recipes.NewMemStore()
	router := mux.NewRouter()
	//define all the routes/the subrouter
	s := router.PathPrefix("/recipes").Subrouter()

	//Initialize the handlers
	home := homeHandler{}
	NewRecipesHandler(store, s)

	//create the router
	router.HandleFunc("/", home.ServeHTTP)

	//start the server
	err := http.ListenAndServe(":8010", router)
	if err != nil {
		return
	}
}

type RecipesHandler struct {
	store recipeStore
}

// NewRecipesHandler registers endpoints and returns a new RecipesHandler
func NewRecipesHandler(s recipeStore, router *mux.Router) *RecipesHandler {
	handler := &RecipesHandler{
		store: s,
	}

	router.HandleFunc("/", handler.ListRecipes).Methods("GET")
	router.HandleFunc("/", handler.CreateRecipe).Methods("POST")
	router.HandleFunc("/{id}", handler.GetRecipe).Methods("GET")
	router.HandleFunc("/{id}", handler.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/{id}", handler.DeleteRecipe).Methods("DELETE")

	return handler
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

func (h RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	//recipe object that will be populated from the JSON payload
	var recipe recipes.Recipe

	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	//Convert the name of the recipe into URL friendly string
	resourceID := slug.Make(recipe.Name)

	//call the store to add the name
	if err := h.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Set the status code to 200
	w.WriteHeader(http.StatusCreated)
}

func (h RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	//retrieve resources from the store
	resources, err := h.store.List()

	//convert list into json
	jsonBytes, err := json.Marshal(resources)

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	//extract resource id
	id := mux.Vars(r)["id"]

	//get the recipe from the store
	recipe, err := h.store.Get(id)

	if err != nil {
		//special case of not found error
		if err == recipes.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}
		// every other error
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(recipe)

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	//extract resource id
	id := mux.Vars(r)["id"]

	//recipe object that will be populated from JSON payload
	var recipe recipes.Recipe

	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Update(id, recipe); err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	//get the resource from the url
	id := mux.Vars(r)["id"]

	if err := h.store.Remove(id); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal server error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not found"))
}

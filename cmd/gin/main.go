package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gitub.com/ndyabagye/go-rest-demo/pkg/recipes"
	"net/http"
)

func main() {
	//create gin router
	router := gin.Default()
	//Initiate recipe handler and provide a data store implementation
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	//register routes
	router.GET("/", homePage)
	router.GET("/recipes", recipesHandler.ListRecipes)
	router.POST("/recipes", recipesHandler.CreateRecipe)
	router.GET("/recipes/:id", recipesHandler.GetRecipe)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipe)
	router.DELETE("/recipes/:id", recipesHandler.RemoveRecipe)

	//start the server
	router.Run()
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
}

type RecipesHandler struct {
	store recipeStore
}

// recipeStore is an interface for the data store
type recipeStore interface {
	Add(name string, recipes recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

// define handler function signatures
func (h RecipesHandler) CreateRecipe(c *gin.Context) {
	//get request body and convert it to recipes.Recipe
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//create a url friendly name
	id := slug.Make(recipe.Name)

	// add to the store
	if err := h.store.Add(id, recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	//return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h RecipesHandler) ListRecipes(c *gin.Context) {
	//call the store to list all the recipes
	r, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	//return the list, JSON encoding it into a list
	c.JSON(http.StatusOK, r)
}
func (h RecipesHandler) GetRecipe(c *gin.Context) {
	//retrieve url parameter
	id := c.Param("id")

	//get the recipe by ID from the store
	recipe, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	//return the recipe, JSON encoding is implicit
	c.JSON(http.StatusOK, recipe)
}

func (h RecipesHandler) UpdateRecipe(c *gin.Context) {
	// Get request body and convert it to recipes.Recipe
	var recipe recipes.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//retrieve the url parameter
	id := c.Param("id")

	//call store to update the recipe
	err := h.store.Update(id, recipe)
	if err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h RecipesHandler) RemoveRecipe(c *gin.Context) {
	//retrieve url param
	id := c.Param("id")

	//call store to delete recipe
	err := h.store.Remove(id)
	if err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// NewRecipesHandler is a constructor for RecipesHandler
func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}

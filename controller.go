package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewController() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basePathHandler)
	mux.HandleFunc("/api/v1/recipes", recipesHandler)
	mux.HandleFunc("/api/v1/recipe", recipeHandler)
	return mux
}

func basePathHandler(writer http.ResponseWriter, request *http.Request) {
	msg := "available endpoints: /api/v1/recipes, /api/v1/recipe"
	writer.Header().Set("Content-Type", "text/plain")
	writer.Header().Set("Content-Length", strconv.Itoa(len(msg)))
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Request-Method", "GET")
	writer.Header().Set("Access-Control-Request-Headers", "*")
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(int((time.Hour * 12).Seconds())))
	writer.Header().Set("Access-Control-Expose-Headers", "*")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(msg))
}

/*
GET /api/v1/recipes?ingredients=food-1,food-2,…[&sort_criteria=[time_to_cook|popularity|num_of_ingredients|...]]

Response
[{
	“recipeName”:<string>,
	"recipeId":<int>,
	“thumbnail”:<base64 string>
}]
*/
func recipesHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ingredients := request.URL.Query()["ingredients"]
	if len(ingredients) == 0 || containsEmptyArgs(ingredients) || containsInvalidArgs(ingredients) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("empty \"ingredients\" query"))
		return
	}

	type recipe struct {
		Name      string `json:"strMeal"`
		Thumbnail string `json:"strMealThumb"`
		Id        string `json:"idMeal"`
	}
	type recipes struct {
		Meals []recipe `json:"meals"`
	}

	var allIngredients []string
	for _, ingredientsQuery := range ingredients {
		allIngredients = append(allIngredients, strings.Split(ingredientsQuery, ",")...)
	}

	url := "https://www.themealdb.com/api/json/v1/1/filter.php?i="
	recipesResult := make(map[string]recipe, 30)

	for _, ingredient := range allIngredients {
		bytes, err := executeGetRequest(url + ingredient)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}

		recipes := recipes{}
		err = json.Unmarshal(bytes, &recipes)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}

		for _, recipe := range recipes.Meals {
			recipesResult[recipe.Name] = recipe
		}
	}

	type resultJson struct {
		RecipeName string `json:"recipeName"`
		RecipeId   int    `json:"recipeId"`
		Thumbnail  string `json:"thumbnail"`
	}

	var result []resultJson

	for name, recipe := range recipesResult {
		thumbnail, err := getThumbnail(request.Header["X-Recipe-Thumbnail"], recipe.Thumbnail)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
		id, _ := strconv.Atoi(recipe.Id)

		result = append(result, resultJson{
			RecipeName: name,
			RecipeId:   id,
			Thumbnail:  thumbnail,
		})
	}

	res, err := json.Marshal(&result)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Length", strconv.Itoa(len(res)))
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Request-Method", "GET")
	writer.Header().Set("Access-Control-Request-Headers", "*")
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(int((time.Hour * 12).Seconds())))
	writer.Header().Set("Access-Control-Expose-Headers", "*")
	writer.WriteHeader(http.StatusOK)
	writer.Write(res)
}

/*
GET /api/v1/recipe/id

Response
{
	“recipeName”:<string>,
	“thumbnail”:<base64 string>,
	“instructions”:<string>,
	“ingredients”:[{
		“name”:<string>,
		“quantity”:<string>
	}],
 	“timeToCook”:<float>
}
*/
func recipeHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlParts := strings.Split(request.URL.Path, "/")
	recipeId := urlParts[len(urlParts) - 1]

	if recipeId == "recipe" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("no recipe id provided"))
		return
	}

	url := "https://www.themealdb.com/api/json/v1/1/lookup.php?i=" + recipeId

	bytes, err := executeGetRequest(url)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	type recipes struct {
		Meals []Recipe `json:"meals"`
	}

	recipesJson := recipes{}
	err = json.Unmarshal(bytes, &recipesJson)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	recipe := recipesJson.Meals[0]

	recipeType := reflect.ValueOf(recipe)

	var ingredients []Ingredient

	for i := 1; i <= 20; i++ {
		ingredientField := recipeType.FieldByName("Ingredient"+strconv.Itoa(i))
		quantityField := recipeType.FieldByName("Measure"+strconv.Itoa(i))

		if ingredientField.String() != "" {
			ingredients = append(ingredients, Ingredient{
				Name:     ingredientField.String(),
				Quantity: quantityField.String(),
			})
		}
	}

	thumbnail, err := getThumbnail(request.Header["X-Recipe-Thumbnail"], recipe.Thumbnail)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	res := Result{
		RecipeName:   recipe.Name,
		Thumbnail:    thumbnail,
		Instructions: recipe.Instructions,
		Ingredients:  ingredients,
		TimeToCook:   0,
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Length", strconv.Itoa(len(resBytes)))
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Request-Method", "GET")
	writer.Header().Set("Access-Control-Request-Headers", "*")
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(int((time.Hour * 12).Seconds())))
	writer.Header().Set("Access-Control-Expose-Headers", "*")
	writer.WriteHeader(http.StatusOK)
	writer.Write(resBytes)
}

func getThumbnail(headerVal []string, thumbnailUrl string) (string, error) {
	if headerVal != nil && headerVal[0] == "BASE64" {
		bytes, err := executeGetRequest(thumbnailUrl)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(bytes), nil
	}
	return strings.ReplaceAll(thumbnailUrl, "\\/", "/"), nil
}

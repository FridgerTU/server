package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	Url "net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
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
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Header().Set("Content-Length", strconv.Itoa(len(msg)))
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := writer.Write([]byte(msg)); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

/*
GET /api/v1/recipes?ingredients=food-1,food-2,…[&sort_criteria=[time_to_cook|popularity|num_of_ingredients|...]]

Response
[{
	“recipeName”:<string>,
	“thumbnail”:<base64 string>
}]
*/
//TODO refactor
func recipesHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ingredients := request.URL.Query()["ingredients"]
	if len(ingredients) == 0 || containsEmptyArgs(ingredients) || containsInvalidArgs(ingredients) {
		writer.WriteHeader(http.StatusBadRequest)
		if _, err := writer.Write([]byte("empty \"ingredients\" query")); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
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
			if _, err := writer.Write([]byte(err.Error())); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}

		recipes := recipes{}
		err = json.Unmarshal(bytes, &recipes)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			if _, err := writer.Write([]byte(err.Error())); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}

		for _, recipe := range recipes.Meals {
			recipesResult[recipe.Name] = recipe
		}
	}

	type resultJson struct {
		RecipeName string `json:"recipeName"`
		Thumbnail  string `json:"thumbnail"`
	}

	var result []resultJson

	for name, recipe := range recipesResult {
		thumbnail, err := getThumbnail(request.Header["X-Recipe-Thumbnail"], recipe.Thumbnail)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			if _, err := writer.Write([]byte(err.Error())); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}

		result = append(result, resultJson{
			RecipeName: name,
			Thumbnail:  thumbnail,
		})
	}

	res, err := json.Marshal(&result)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Length", strconv.Itoa(len(res)))
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := writer.Write(res); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

/*
GET /api/v1/recipe?name=...

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
//TODO refactor
func recipeHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	recipeName := request.URL.Query()["name"]
	if len(recipeName) == 0 || len(recipeName) > 1 || containsEmptyArgs(recipeName) || containsInvalidArgs(recipeName) {
		writer.WriteHeader(http.StatusBadRequest)
		if _, err := writer.Write([]byte("invalid \"name\" query")); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	url := "https://www.themealdb.com/api/json/v1/1/search.php?s=" + Url.QueryEscape(recipeName[0])

	bytes, err := executeGetRequest(url)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	type recipes struct {
		Meals []Recipe `json:"meals"`
	}

	recipesJson := recipes{}
	err = json.Unmarshal(bytes, &recipesJson)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	recipe := recipesJson.Meals[0]

	recipeType := reflect.ValueOf(recipe)

	var ingredients []Ingredient

	for i := 1; i <= 20; i++ {
		ingredientField := recipeType.FieldByName("Ingredient"+strconv.FormatInt(int64(i), 10))
		quantityField := recipeType.FieldByName("Measure"+strconv.FormatInt(int64(i), 10))

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
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
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
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Length", strconv.Itoa(len(resBytes)))
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := writer.Write(resBytes); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getThumbnail(headerVal []string, thumbnailUrl string) (string, error) {
	if headerVal != nil && headerVal[0] == "BASE64" {
		bytes, err := executeGetRequest(thumbnailUrl + "/preview")
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(bytes), nil
	}
	return strings.ReplaceAll(thumbnailUrl, "\\/", "/") + "/preview", nil
}

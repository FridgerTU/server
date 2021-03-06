package main

type Recipe struct {
	Id               string `json:"idMeal"`
	Name             string `json:"strMeal"`
	DrinkAlternative string `json:"strDrinkAlternate"`
	Category         string `json:"strCategory"`
	Area             string `json:"strArea"`
	Instructions     string `json:"strInstructions"`
	Thumbnail        string `json:"strMealThumb"`
	Tags             string `json:"strTags"`
	YouTubeLink      string `json:"strYoutube"`

	Ingredient1  string `json:"strIngredient1"`
	Ingredient2  string `json:"strIngredient2"`
	Ingredient3  string `json:"strIngredient3"`
	Ingredient4  string `json:"strIngredient4"`
	Ingredient5  string `json:"strIngredient5"`
	Ingredient6  string `json:"strIngredient6"`
	Ingredient7  string `json:"strIngredient7"`
	Ingredient8  string `json:"strIngredient8"`
	Ingredient9  string `json:"strIngredient9"`
	Ingredient10 string `json:"strIngredient10"`
	Ingredient11 string `json:"strIngredient11"`
	Ingredient12 string `json:"strIngredient12"`
	Ingredient13 string `json:"strIngredient13"`
	Ingredient14 string `json:"strIngredient14"`
	Ingredient15 string `json:"strIngredient15"`
	Ingredient16 string `json:"strIngredient16"`
	Ingredient17 string `json:"strIngredient17"`
	Ingredient18 string `json:"strIngredient18"`
	Ingredient19 string `json:"strIngredient19"`
	Ingredient20 string `json:"strIngredient20"`

	Measure1  string `json:"strMeasure1"`
	Measure2  string `json:"strMeasure2"`
	Measure3  string `json:"strMeasure3"`
	Measure4  string `json:"strMeasure4"`
	Measure5  string `json:"strMeasure5"`
	Measure6  string `json:"strMeasure6"`
	Measure7  string `json:"strMeasure7"`
	Measure8  string `json:"strMeasure8"`
	Measure9  string `json:"strMeasure9"`
	Measure10 string `json:"strMeasure10"`
	Measure11 string `json:"strMeasure11"`
	Measure12 string `json:"strMeasure12"`
	Measure13 string `json:"strMeasure13"`
	Measure14 string `json:"strMeasure14"`
	Measure15 string `json:"strMeasure15"`
	Measure16 string `json:"strMeasure16"`
	Measure17 string `json:"strMeasure17"`
	Measure18 string `json:"strMeasure18"`
	Measure19 string `json:"strMeasure19"`
	Measure20 string `json:"strMeasure20"`

	Source       string `json:"strSource"`
	DateModified string `json:"dateModified"`
}

type Ingredient struct {
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
}

type Result struct {
	RecipeName   string       `json:"recipeName"`
	Thumbnail    string       `json:"thumbnail"`
	Instructions string       `json:"instructions"`
	Ingredients  []Ingredient `json:"ingredients"`
	TimeToCook   float32      `json:"timeToCook"`
}

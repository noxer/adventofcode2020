package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	products, err := readIngredients("input.txt")
	if err != nil {
		fmt.Printf("Konnte die Zutatenliste nicht lesen: %s\n", err)
		os.Exit(1)
	}

	allergens := make(map[string]Ingredients)
	for _, product := range products {
		for _, allergen := range product.Allergens {
			if ingedients, ok := allergens[allergen]; ok {
				allergens[allergen] = ingedients.Overlap(product.Ingredients)
			} else {
				allergens[allergen] = product.Ingredients
			}
		}
	}

	changed := true
	for changed {
		changed = false

		for _, ingredients := range allergens {
			if len(ingredients) == 1 {
				changed = changed || removeFromMap(allergens, ingredients.First())
			}
		}
	}

	shortAllergens := make(map[string]string)
	shortAllergensInv := make(map[string]string)
	allergenNames := make([]string, 0)
	for allergen, ingredients := range allergens {
		shortAllergens[ingredients.First()] = allergen
		shortAllergensInv[allergen] = ingredients.First()

		allergenNames = append(allergenNames, allergen)
	}

	sort.Strings(allergenNames)

	//	for allergen, ingredients := range allergens {
	//		fmt.Printf("%s: %v\n", allergen, ingredients)
	//	}

	count := 0
	for _, product := range products {
		for ingredient := range product.Ingredients {
			if _, ok := shortAllergens[ingredient]; !ok {
				count++
			}
		}
	}

	fmt.Printf("%d Zutaten ohne Allergen\n", count)

	var allergicIng []string
	for _, allergen := range allergenNames {
		allergicIng = append(allergicIng, shortAllergensInv[allergen])
	}

	fmt.Printf("Allergene: %s\n", strings.Join(allergicIng, ","))
}

func removeFromMap(allergens map[string]Ingredients, ingredient string) bool {
	removed := false
	for _, ingredients := range allergens {
		if len(ingredients) > 1 {
			delete(ingredients, ingredient)
			removed = true
		}
	}
	return removed
}

// Ingredients ...
type Ingredients map[string]struct{}

// Overlap ...
func (i Ingredients) Overlap(o Ingredients) Ingredients {
	result := make(Ingredients)

	for ingredient := range i {
		if _, ok := o[ingredient]; ok {
			result[ingredient] = struct{}{}
		}
	}

	return result
}

// First ...
func (i Ingredients) First() string {
	for ingredient := range i {
		return ingredient
	}

	return ""
}

// Product ...
type Product struct {
	Ingredients Ingredients
	Allergens   []string
}

func readIngredients(name string) ([]Product, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var products []Product
	for s.Scan() {
		line := s.Text()

		parts := strings.SplitN(line, " (contains ", 2)
		ingredients := strings.Split(parts[0], " ")
		allergens := strings.Split(strings.TrimSuffix(parts[1], ")"), ", ")

		ingredientSet := make(Ingredients)
		for _, ingredient := range ingredients {
			ingredientSet[ingredient] = struct{}{}
		}

		products = append(products, Product{Ingredients: ingredientSet, Allergens: allergens})
	}

	return products, nil
}

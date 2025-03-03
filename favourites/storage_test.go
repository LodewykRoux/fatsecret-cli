package favourites

import (
	"os"
	"testing"
)

var tmpFileName = "favourites_test_*.csv"

func SeedData(fileName string) {
	SaveFavourite(Favourite{Name: "Apple", Uom: "1 piece", Calories: "95", Fat: "0.3", Carbs: "25", Protein: "0.5"}, fileName)
	SaveFavourite(Favourite{Name: "Banana", Uom: "1 piece", Calories: "105", Fat: "0.3", Carbs: "27", Protein: "1.3"}, fileName)
}

func SetupTempFile(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", tmpFileName)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	return tmpFile.Name()
}

func TestFormatFavourite(t *testing.T) {
	fav := Favourite{Name: "Apple", Uom: "1 piece", Calories: "95", Fat: "0.3", Carbs: "25", Protein: "0.5"}
	expected := "    🍽  Food: Apple\n" +
		"    📏 Serving Size: 1 piece\n" +
		"    🔥 Calories: 95 kcal\n" +
		"    🥑 Fat: 0.3 g\n" +
		"    🍞 Carbs: 25 g\n" +
		"    🍗 Protein: 0.5 g\n"

	result := FormatFavourite(fav)

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestSaveAndRetrieveFavourites(t *testing.T) {
	fName := SetupTempFile(t)
	fav := Favourite{Name: "Apple", Uom: "1 piece", Calories: "95", Fat: "0.3", Carbs: "25", Protein: "0.5"}

	if err := SaveFavourite(fav, fName); err != nil {
		t.Fatalf("Failed to save favourite: %v", err)
	}

	favs, err := RetrieveFavourites(fName)
	if err != nil {
		t.Fatalf("Failed to retrieve favourites: %v", err)
	}

	if len(favs) != 1 || favs[0].Name != "Apple" {
		t.Errorf("Expected Apple, got %+v", favs)
	}
}

func TestDeleteFavourite(t *testing.T) {
	fName := SetupTempFile(t)

	SeedData(fName)

	if err := DeleteFavourite(Favourite{Name: "Apple"}, fName); err != nil {
		t.Fatalf("Failed to delete favourite: %v", err)
	}

	favs, _ := RetrieveFavourites(fName)
	if len(favs) != 1 || favs[0].Name != "Banana" {
		t.Errorf("Expected Banana, got %+v", favs)
	}
}

func TestFilterFavourites(t *testing.T) {
	fName := SetupTempFile(t)

	SeedData(fName)

	favs, _ := RetrieveFavourites(fName)
	searchTerm := "Apple"

	filtered, err := FilterFavourites(favs, searchTerm)
	if err != nil {
		t.Fatalf("FilterFavourites returned an error: %v", err)
	}

	if len(filtered) != 1 || filtered[0].Name != searchTerm {
		t.Errorf("Expected Apple, got %+v", filtered)
	}
}

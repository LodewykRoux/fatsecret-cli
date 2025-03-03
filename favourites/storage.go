package favourites

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func GetFavouriteFile() string {
	return "favourites.csv"
}

type Favourite struct {
	Name     string
	Uom      string
	Calories string
	Fat      string
	Carbs    string
	Protein  string
}

func FormatFavourite(favourite Favourite) string {
	return fmt.Sprintf(
		"    🍽  Food: %s\n"+
			"    📏 Serving Size: %s\n"+
			"    🔥 Calories: %s kcal\n"+
			"    🥑 Fat: %s g\n"+
			"    🍞 Carbs: %s g\n"+
			"    🍗 Protein: %s g\n",
		favourite.Name,
		favourite.Uom,
		favourite.Calories,
		favourite.Fat,
		favourite.Carbs,
		favourite.Protein,
	)
}

func SaveFavourite(favourite Favourite, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{
		favourite.Name,
		favourite.Uom,
		favourite.Calories,
		favourite.Fat,
		favourite.Carbs,
		favourite.Protein,
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteFavourite(favourite Favourite, fileName string) error {
	favourites, err := RetrieveFavourites(fileName)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Rewrite all favourites except the one to be deleted
	for _, f := range favourites {
		if f.Name == favourite.Name {
			continue
		}
		writer.Write([]string{
			f.Name,
			f.Uom,
			f.Calories,
			f.Fat,
			f.Carbs,
			f.Protein,
		})
	}

	return nil
}

func RetrieveFavourites(fileName string) ([]Favourite, error) {
	var favourites []Favourite

	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return favourites, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	for _, record := range records {
		favourites = append(favourites, Favourite{
			Name:     record[0],
			Uom:      record[1],
			Calories: record[2],
			Fat:      record[3],
			Carbs:    record[4],
			Protein:  record[5],
		})
	}

	return favourites, nil
}

func FilterFavourites(favourites []Favourite, searchTerm string) ([]Favourite, error) {
	var favsFiltered []Favourite
	for _, val := range favourites {
		if strings.Contains(strings.ToLower(val.Name), strings.ToLower(searchTerm)) {
			favsFiltered = append(favsFiltered, val)
		}
	}

	return favsFiltered, nil
}

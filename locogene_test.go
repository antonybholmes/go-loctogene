package loctogene

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/antonybholmes/go-dna"
)

func TestWithin(t *testing.T) {
	fmt.Println("Within")

	file := fmt.Sprintf("../data/loctogene/%s.db", "grch38")
	db, err := GetDB(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	location := dna.Location{Chr: "chr3", Start: 187721370, End: 187733550}

	records, err := GetGenesWithin(db, &location, 1)

	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(&records)
	fmt.Printf("%s", string(b))
}

func TestClosest(t *testing.T) {
	fmt.Println("Closest")

	file := fmt.Sprintf("../data/loctogene/%s.db", "grch38")

	db, err := GetDB(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	location := dna.Location{Chr: "chr3", Start: 187721377, End: 187745725}

	records, err := ClosestGenes(db, &location, 10, 1)

	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(&records)
	fmt.Printf("%s", string(b))
}

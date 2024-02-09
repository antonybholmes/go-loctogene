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
	db, err := NewLoctogeneDB(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	location := dna.NewLocation("chr3", 187721370, 187733550)

	records, err := db.WithinGenes(location, Gene)

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

	db, err := NewLoctogeneDB(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	location := dna.NewLocation("chr3", 187721377, 187745725)

<<<<<<< HEAD
	records, err := ClosestGenes(db, &location, 10, 1)
=======
	records, err := db.ClosestGenes(location, 10, 1)
>>>>>>> dev

	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(&records)
	fmt.Printf("%s", string(b))
}

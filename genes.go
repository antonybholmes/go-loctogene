package loctogene

import (
	"database/sql"
	"fmt"

	"github.com/antonybholmes/go-dna"
)

const WITHIN_GENE_SQL = "SELECT id, chr, start, end, strand, gene_id, gene_symbol, start - ? " +
	"FROM genes " +
	"WHERE level=? AND chr=? AND ((start <= ? AND end >= ?) OR (start <= ? AND end >= ?)) " +
	"ORDER BY start ASC"

const CLOSEST_GENE_SQL = "SELECT id, chr, start, end, strand, gene_id, gene_symbol, stranded_start - ? " +
	"FROM genes " +
	"WHERE level=? AND chr=? " +
	"ORDER BY ABS(stranded_start - ?) " +
	"LIMIT ?"

type FeatureRecord struct {
	Id         int    `json:"id"`
	Chr        string `json:"chr"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
	Strand     string `json:"strand"`
	GeneId     string `json:"gene_id"`
	GeneSymbol string `json:"gene_symbol"`
	Dist       int    `json:"d"`
}

type Features struct {
	Loc      string          `json:"loc"`
	Level    string          `json:"level"`
	Features []FeatureRecord `json:"features"`
}

func GetLevel(level string) int {
	switch level {
	case "transcript", "2":
		return 2
	case "exon", "3":
		return 3
	default:
		return 1
	}
}

func GetLevelType(level int) string {
	switch level {
	case 2:
		return "transcript"
	case 3:
		return "exon"
	default:
		return "gene"
	}
}

func GetGenesWithin(db *sql.DB, location *dna.Location, level int) (*Features, error) {
	mid := (location.Start + location.End) / 2

	rows, err := db.Query(WITHIN_GENE_SQL,
		mid,
		level,
		location.Chr,
		location.Start,
		location.Start,
		location.End,
		location.End)

	if err != nil {
		return nil, err //fmt.Errorf("there was an error with the database query")
	}

	return RowsToRecords(location, rows, level)
}

func GetClosestGenes(db *sql.DB, location *dna.Location, n int, level int) (*Features, error) {
	mid := (location.Start + location.End) / 2

	rows, err := db.Query(CLOSEST_GENE_SQL,
		mid,
		level,
		location.Chr,
		mid,
		n)

	if err != nil {
		return nil, err //fmt.Errorf("there was an error with the database query")
	}

	return RowsToRecords(location, rows, level)
}

func RowsToRecords(location *dna.Location, rows *sql.Rows, level int) (*Features, error) {
	defer rows.Close()

	var id int
	var chr string
	var start int
	var end int
	var strand string
	var geneId string
	var geneSymbol string
	var d int
	t := GetLevelType(level)

	var records = []FeatureRecord{}

	for rows.Next() {
		err := rows.Scan(&id, &chr, &start, &end, &strand, &geneId, &geneSymbol, &d)

		if err != nil {
			return nil, err //fmt.Errorf("there was an error with the database records")
		}

		if strand == "-" {
			t := start
			start = end
			end = t
		}

		records = append(records, FeatureRecord{Id: id, Chr: chr, Start: start, End: end, Strand: strand, GeneId: geneId, GeneSymbol: geneSymbol, Dist: d})
	}

	return &Features{Loc: fmt.Sprintf("%s:%d-%d", location.Chr, location.Start, location.End), Level: t, Features: records}, nil
}

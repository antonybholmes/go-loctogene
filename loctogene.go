package loctogene

import (
	"database/sql"

	"github.com/antonybholmes/go-dna"
)

const WITHIN_GENE_SQL = `SELECT id, chr, start, end, strand, gene_id, gene_symbol, start - ?
 FROM genes
 WHERE level=? AND chr=? AND ((start <= ? AND end >= ?) OR (start <= ? AND end >= ?))
 ORDER BY start ASC`

const CLOSEST_GENE_SQL = `SELECT id, chr, start, end, strand, gene_id, gene_symbol, stranded_start - ?
 FROM genes
 WHERE level=? AND chr=?
 ORDER BY ABS(stranded_start - ?)
 LIMIT ?`

const WITHIN_GENE_AND_PROMOTER_SQL = `SELECT id, chr, start, end, strand, gene_id, gene_symbol, start - ? 
 FROM genes 
 WHERE level = ? AND chr = ? AND ((start - ? <= ? AND end + ? >= ?) OR (start - ? <= ? AND end + ? >= ?)) 
 ORDER BY start ASC`

const IN_EXON_SQL = `SELECT id, chr, start, end, strand, gene_id, gene_symbol, start - ? 
 FROM genes 
 WHERE level=3 AND gene_id=? AND chr=? AND ((start <= ? AND end >= ?) OR (start <= ? AND end >= ?)) 
 ORDER BY start ASC`

const IN_PROMOTER_SQL = `SELECT id, chr, start, end, strand, gene_id, gene_symbol, start - ? 
 FROM genes 
 WHERE level=2 AND gene_id=? AND chr=? AND ? >= stranded_start - ? AND ? <= stranded_start + ? 
 ORDER BY start ASC`

type GenomicFeature struct {
	Id         int    `json:"id"`
	Chr        string `json:"chr"`
	Start      uint   `json:"start"`
	End        uint   `json:"end"`
	Strand     string `json:"strand"`
	GeneId     string `json:"gene_id"`
	GeneSymbol string `json:"gene_symbol"`
	Dist       int    `json:"d"`
}

type GenomicFeatures struct {
	Location string           `json:"location"`
	Level    string           `json:"level"`
	Features []GenomicFeature `json:"features"`
}

var ERROR_FEATURES = GenomicFeatures{Location: "", Level: "", Features: []GenomicFeature{}}

type Level int

const (
	Gene Level = iota
	Transcript
	Exon
)

func ParseLevel(level string) Level {
	switch level {
	case "transcript", "2":
		return Transcript
	case "exon", "3":
		return Exon
	default:
		return Gene
	}
}

func (level Level) String() string {
	switch level {
	case Exon:
		return "Exon"
	case Transcript:
		return "Transcript"
	default:
		return "Gene"
	}
}

type LoctogeneDB struct {
	db                    *sql.DB
	withinGeneStmt        *sql.Stmt
	withinGeneAndPromStmt *sql.Stmt
	inExonStmt            *sql.Stmt
	closestGeneStmt       *sql.Stmt
}

func NewLoctogeneDB(file string) (*LoctogeneDB, error) {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	withinGeneStmt, err := db.Prepare(WITHIN_GENE_SQL)

	if err != nil {
		return nil, err
	}

	withinGeneAndPromStmt, err := db.Prepare(WITHIN_GENE_AND_PROMOTER_SQL)

	if err != nil {
		return nil, err
	}

	inExonStmt, err := db.Prepare(IN_EXON_SQL)

	if err != nil {
		return nil, err
	}

	closestGeneStmt, err := db.Prepare(CLOSEST_GENE_SQL)

	if err != nil {
		return nil, err
	}

	return &LoctogeneDB{db, withinGeneStmt, withinGeneAndPromStmt, inExonStmt, closestGeneStmt}, nil
}

func (loctogenedb *LoctogeneDB) Close() {
	loctogenedb.db.Close()
}

func (loctogenedb *LoctogeneDB) WithinGenes(location *dna.Location, level Level) (*GenomicFeatures, error) {
	mid := (location.Start + location.End) / 2

	rows, err := loctogenedb.withinGeneStmt.Query(
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

	return rowsToRecords(location, rows, level)
}

func (loctogenedb *LoctogeneDB) WithinGenesAndPromoter(location *dna.Location, level Level, pad uint) (*GenomicFeatures, error) {
	mid := (location.Start + location.End) / 2

	rows, err := loctogenedb.withinGeneAndPromStmt.Query(
		mid,
		level,
		location.Chr,
		pad,
		location.Start,
		pad,
		location.Start,
		pad,
		location.End,
		pad,
		location.End)

	if err != nil {
		return nil, err //fmt.Errorf("there was an error with the database query")
	}

	return rowsToRecords(location, rows, level)
}

func (loctogenedb *LoctogeneDB) InExon(location *dna.Location, geneId string) (*GenomicFeatures, error) {
	mid := (location.Start + location.End) / 2

	rows, err := loctogenedb.inExonStmt.Query(
		mid,
		geneId,
		location.Chr,
		location.Start,
		location.Start,
		location.End,
		location.End)

	if err != nil {
		return nil, err //fmt.Errorf("there was an error with the database query")
	}

	return rowsToRecords(location, rows, Exon)
}

func (loctogenedb *LoctogeneDB) ClosestGenes(location *dna.Location, n uint16, level Level) (*GenomicFeatures, error) {
	mid := (location.Start + location.End) / 2

	rows, err := loctogenedb.closestGeneStmt.Query(mid,
		level,
		location.Chr,
		mid,
		n)

	if err != nil {
		return nil, err //fmt.Errorf("there was an error with the database query")
	}

	return rowsToRecords(location, rows, level)
}

func rowsToRecords(location *dna.Location, rows *sql.Rows, level Level) (*GenomicFeatures, error) {
	defer rows.Close()

	var id int
	var chr string
	var start uint
	var end uint
	var strand string
	var geneId string
	var geneSymbol string
	var d int

	var features = []GenomicFeature{}

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

		features = append(features, GenomicFeature{Id: id, Chr: chr, Start: start, End: end, Strand: strand, GeneId: geneId, GeneSymbol: geneSymbol, Dist: d})
	}

	return &GenomicFeatures{Location: location.String(), Level: level.String(), Features: features}, nil
}

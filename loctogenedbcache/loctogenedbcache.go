package loctogenedbcache

import (
	"github.com/antonybholmes/go-loctogene"
)

var Cache = loctogene.NewLoctogeneDbCache()

func Dir(dir string) {
	Cache.Dir(dir)

}
func Db(assembly string) (*loctogene.LoctogeneDb, error) {
	return Cache.Db(assembly)
}

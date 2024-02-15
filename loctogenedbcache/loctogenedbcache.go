package loctogenedbcache

import (
	"github.com/antonybholmes/go-loctogene"
)

var Cache = loctogene.NewLoctogeneDbCache()

func Dir(dir string) *loctogene.LoctogeneDbCache {
	Cache.Dir(dir)
	return Cache
}
func Db(assembly string) (*loctogene.LoctogeneDb, error) {
	return Cache.Db(assembly)
}

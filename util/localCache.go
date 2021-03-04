package util

import (
	"context"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type GraphCache struct {
	Cache      *cache.Cache
	DocsByType map[string][]string
}

func NewCache() *GraphCache {
	gob.Register(eos.Name(""))
	gob.Register(eos.Asset{})
	gob.Register(eos.TimePoint(0))
	gob.Register(eos.Checksum256{})
	gob.Register(docgraph.Document{})
	gob.Register(docgraph.Edge{})
	gob.Register(map[string][]string{})
	gob.Register(map[string]string{})

	gc := GraphCache{}
	gc.Cache = cache.New(60*time.Minute, 60*time.Minute)
	gc.DocsByType = make(map[string][]string)
	return &gc
}

func GetCache(ctx context.Context, api *eos.API, contract eos.AccountName) (*GraphCache, error) {

	gCache := NewCache()
	err := gCache.Cache.LoadFile(".graph.cache")
	if err != nil {
		zap.S().Debugf("Unable to load cache: %v", err)
		zap.S().Debug("Cache file not found, building fresh one")
		return FreshCache(ctx, api, contract)
	}
	var found bool
	dbt, found := gCache.Cache.Get("DocsByType")
	if !found {
		zap.S().Debug("DocsByType was not found in the cache, assume it is expired and building a freshie")
		return FreshCache(ctx, api, contract)
	}
	gCache.DocsByType = dbt.(map[string][]string)
	zap.S().Debug("Cache file found, loading into memory")
	return gCache, nil
}

func FreshCache(ctx context.Context, api *eos.API, contract eos.AccountName) (*GraphCache, error) {

	gCache := NewCache()

	documents, err := docgraph.GetAllDocuments(ctx, api, contract)
	if err != nil {
		return nil, fmt.Errorf("cannot get all documents: %v", err)
	}

	for _, document := range documents {

		docType, err := document.GetType()
		if err != nil {
			return nil, fmt.Errorf("document with invalid or missing type: %v %v", document.Hash.String(), err)
		}
		gCache.DocsByType[string(docType)] = append(gCache.DocsByType[string(docType)], document.Hash.String())
		gCache.Cache.Set(document.Hash.String()[:5], document.Hash.String(), cache.DefaultExpiration)
		gCache.Cache.Set(document.Hash.String(), document, cache.DefaultExpiration)
	}
	gCache.Cache.Set("DocsByType", gCache.DocsByType, cache.DefaultExpiration)

	edges, err := docgraph.GetAllEdges(ctx, api, contract)
	if err != nil {
		return nil, fmt.Errorf("cannot get all edges: %v", err)
	}

	for _, edge := range edges {
		gCache.Cache.Set(strconv.Itoa(int(edge.ID)), edge, cache.DefaultExpiration)
	}

	err = gCache.Cache.SaveFile(".graph.cache")
	if err != nil {
		return nil, fmt.Errorf("cannot save cache to file .graph.cache: %v", err)
	}

	return gCache, nil
}

func Get(ctx context.Context, api *eos.API, contract eos.AccountName, hash string) (docgraph.Document, error) {

	gc, err := GetCache(ctx, api, contract)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("cannot get cache: %v", err)
	}

	cachedItem, found := gc.Cache.Get(hash)
	if !found {
		zap.S().Debugf("Document is not found in cache: %v . Loading from blockchain", hash)

		loadedDoc, err := docgraph.LoadDocument(ctx, api, contract, hash)
		if err != nil {
			return docgraph.Document{}, fmt.Errorf("Unable to load document directly from blockchain: %v %v", hash, err)
		}
		return loadedDoc, nil
	} else {
		zap.S().Debugf("Pass through key found, reading 2nd level Document from cache: %v", hash)
		switch x := cachedItem.(type) {
		case docgraph.Document:
			return x, nil
		case string:
			cachedDocument, found := gc.Cache.Get(x)
			if !found {
				zap.S().Debugf("2nd level Document is not found in cache: %v; loading from blockchain", x)

				loadedDoc, err := docgraph.LoadDocument(ctx, api, contract, x)
				if err != nil {
					return docgraph.Document{}, fmt.Errorf("unable to load 2nd level document directly from blockchain: %v %v", hash, err)
				}
				return loadedDoc, nil
			}
			return cachedDocument.(docgraph.Document), nil
		default:
			return docgraph.Document{}, fmt.Errorf("item in cache is neither a Document nor a string: %v %v", hash, err)
		}
	}

}

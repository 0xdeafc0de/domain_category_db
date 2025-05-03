package db

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/armon/go-radix"
	lru "github.com/hashicorp/golang-lru"
)

const (
	DBStoreDir     = "dbstore"
	BaseCategoryID = 101 //Category IDs starts from this ID
)

type CategoryDB struct {
	Cache              *lru.Cache
	FullDB             map[string]uint8
	fullDBHashed       map[uint64]uint8 // When using hashes
	radixTree          *radix.Tree      // Optional: Radix tree for lookups
	useRadix           bool             // Flag to use radix instead of map
	useHash            bool
	CategoriesIdToName map[uint8]string
	CategoriesNameToId map[string]uint8
}

func NewCategoryDB(UseRadix, UseHash bool) *CategoryDB {
	cache, _ := lru.New(1_000_000) // 1M LRU cache
	return &CategoryDB{
		Cache:              cache,
		useRadix:           UseRadix,
		useHash:            UseHash,
		radixTree:          radix.New(),
		FullDB:             make(map[string]uint8),
		fullDBHashed:       make(map[uint64]uint8),
		CategoriesIdToName: make(map[uint8]string),
		CategoriesNameToId: make(map[string]uint8),
	}
}

func (db *CategoryDB) AddCategory(id uint8, name string) {
	db.CategoriesIdToName[id] = name
	db.CategoriesNameToId[name] = id
}

func (db *CategoryDB) LoadDomainsFromURL(dbStorePath, url, category string) (error, int) {
	id, ok := db.CategoriesNameToId[category]
	if !ok {
		id = BaseCategoryID + uint8(len(db.CategoriesNameToId))
		db.AddCategory(id, category)
	}

	dirPath := filepath.Join(dbStorePath, category)
	os.MkdirAll(dirPath, os.ModePerm)
	fileName := filepath.Base(url)
	filePath := filepath.Join(dirPath, fileName)

	var reader io.Reader
	local := false
	file, err := os.Open(filePath)
	sz := 0

	if err == nil {
		defer file.Close()
		reader = file
		local = true
	} else {
		fmt.Println("*** Downloading from URL", url)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download from %s: %v", url, err), sz
		}
		defer resp.Body.Close()

		f, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", filePath, err), sz
		}
		defer f.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save to file %s: %v", filePath, err), sz
		}
		f.Seek(0, io.SeekStart)
		reader = f
	}

	scanner := bufio.NewScanner(reader)
	count := 0
	for scanner.Scan() {
		//domain := strings.TrimSpace(scanner.Text())
		domain := normalizeDomain(scanner.Text())
		if domain == "" || strings.HasPrefix(domain, "#") {
			continue
		}
		if db.useRadix {
			db.radixTree.Insert(domain, id)
		} else if db.useHash {
			hash := hashDomain(domain)
			db.fullDBHashed[hash] = id
		} else {
			db.FullDB[domain] = id
		}
		count++
	}

	// Get the total DB Size
	if db.useRadix {
		sz = db.radixTree.Len()
	} else if db.useHash {
		sz = len(db.fullDBHashed)
	} else {
		sz = len(db.FullDB)
	}

	if local {
		fmt.Printf("Loaded %d domains from %s into category %s\n", count, dirPath, category)
	} else {
		fmt.Printf("Loaded %d domains from %s into category %s\n", count, url, category)
	}
	return scanner.Err(), sz
}

func (db *CategoryDB) Lookup(domain string) (string, bool, int64) {
	start := time.Now()
	normDomain := normalizeDomain(domain)
	if val, ok := db.Cache.Get(normDomain); ok {
		if catId, ok2 := val.(uint8); ok2 {
			return db.CategoriesIdToName[catId], true, time.Since(start).Nanoseconds()
		}
		// Fallback if the type is not as expected
	}

	var catId uint8
	var found bool
	if db.useRadix {
		val, exists := db.radixTree.Get(domain)
		if exists {
			catId = val.(uint8)
			found = true
		}
	} else if db.useHash {
		hash := hashDomain(normDomain)
		catId, found = db.fullDBHashed[hash]
	} else {
		catId, found = db.FullDB[normDomain]
	}

	if found {
		db.Cache.Add(normDomain, catId)
	}

	return db.CategoriesIdToName[catId], found, time.Since(start).Nanoseconds()
}

func normalizeDomain(domain string) string {
	domain = strings.ToLower(domain)
	domain = strings.TrimSpace(domain)
	return domain
}

func hashDomain(domain string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(domain))
	return h.Sum64()
}

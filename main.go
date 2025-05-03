package main

import (
	"github.com/0xdeafc0de/domain-category-db/config"
	"github.com/0xdeafc0de/domain-category-db/db"
	"github.com/0xdeafc0de/domain-category-db/rest"
	"log"
)

const UseRadix = false
const UseHasdDB = false
const MB = 1024 * 1024

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbInstance := db.NewCategoryDB(UseRadix, UseHasdDB)
	log.Println("DBStore Path = ", cfg.DBStorePath)
	n := 0
	totalCnt := 0
	totalSz := 0

	for _, src := range cfg.Categories {
		//log.Printf("Loading category '%s' from %s", src.Category, src.URL)
		err, dbCount, dbSize := dbInstance.LoadDomainsFromURL(cfg.DBStorePath, src.URL, src.Category)
		if err != nil {
			log.Printf("Error in loading URL %s. Err = %v", src.URL, err)
			continue
		}
		n++
		totalCnt += dbCount
		totalSz += dbSize
	}
	log.Printf("Total %d categories loaded. Total DB Count = %d, Size ~%.2f MB\n", n, totalCnt, float64(totalSz)/MB)

	rest.StartServer(dbInstance)
	log.Println("Exiting")
}

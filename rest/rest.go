package rest

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/0xdeafc0de/domain-category-db/db"
	"sort"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

func StartServer(categoryDB *db.CategoryDB) {
	r := mux.NewRouter()

	r.HandleFunc("/lookup", func(w http.ResponseWriter, r *http.Request) {
		domain := r.URL.Query().Get("domain")
		if domain == "" {
			http.Error(w, "Missing domain parameter", http.StatusBadRequest)
			return
		}
		cat, ok, nanos := categoryDB.Lookup(domain)
		if !ok {
			http.Error(w, "Domain not found", http.StatusNotFound)
			fmt.Fprintf(w, "Lookup Time: %d ns\n", nanos)
			return
		}
		fmt.Fprintf(w, "Category: %s\nLookup Time: %d ns\n", cat, nanos)
	})

	r.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		totalEntries := len(categoryDB.FullDB)
		approxBytes := 0
		for domain := range categoryDB.FullDB {
			approxBytes += len(domain) + 1 // domain string + 1 byte for category ID
		}
		fmt.Fprintf(w, "=== Domain Category DB Info ===\n")
		fmt.Fprintf(w, "Total DB entries         : %d\n", totalEntries)
		fmt.Fprintf(w, "Approximate DB size      : %d bytes\n", approxBytes)
		fmt.Fprintf(w, "Cached entries (LRU)     : %d\n", categoryDB.Cache.Len())
		fmt.Fprintf(w, "Categories (ID → Name)   :\n")

		var keys []int
		// Extract and sort the keys
		for id := range categoryDB.CategoriesIdToName {
			keys = append(keys, int(id))
		}
		sort.Ints(keys)

		// Print the map in sorted order by key
		for _, k := range keys {
			fmt.Fprintf(w, "  %d → %s\n", k, categoryDB.CategoriesIdToName[uint8(k)])
		}
	})

	go func() {
		fmt.Println("Starting pprof server on :8082")
		fmt.Println(http.ListenAndServe("localhost:8082", nil))
	}()

	fmt.Println("Starting server on port 8081")
	http.ListenAndServe(":8081", r)
	fmt.Println("Server exiting")
}

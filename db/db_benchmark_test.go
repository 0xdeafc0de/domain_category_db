package db

import (
	"bufio"
	"fmt"
	"github.com/0xdeafc0de/domain_category_db/config"
	dm "github.com/0xdeafc0de/domain_category_db/domains"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var dbstorePath = "../dbstore"

// setupBenchmarkDB loads a small fixed set of domains into the DB and cache
func setupBenchmarkDB(b *testing.B) *CategoryDB {
	db := NewCategoryDB(false, false)
	db.AddCategory(0, "test")

	// Adjust these according to your config
	category := "phishing"        // Example category name
	fileName := "phishing-nl.txt" // Last part of the URL
	path := filepath.Join("..", "dbstore", category, fileName)

	//fmt.Println("Opening file", path)
	file, err := os.Open(path)
	if err != nil {
		b.Fatalf("Benchmark setup error: missing file %s. Run the main program to populate it first.", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() && count < 1000 {
		domain := scanner.Text()
		db.FullDB[domain] = 0 // Simulate category ID 0
		db.Cache.Add(domain, uint8(0))
		count++
	}

	if err := scanner.Err(); err != nil {
		panic("failed to read domain file: " + err.Error())
	}

	return db
}

func BenchmarkLookup_Cached(b *testing.B) {
	db := setupBenchmarkDB(b)

	//domain := "ozzon-mobi-age.com"
	domains := make([]string, 0, 1000)
	for domain := range db.FullDB {
		domains = append(domains, domain)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		domain := domains[i%len(domains)]
		_, _, _ = db.Lookup(domain)
	}
}

func BenchmarkLookup_FullDB(b *testing.B) {
	db := setupBenchmarkDB(b)
	domain := "domain999999.com" // Not in cache

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = db.Lookup(domain)
	}
}

// SelectDomains selects N domains form the file with equal distribution from begingin, mid and end
func SelectDomains(N int, filename string) ([]string, error) {
	fmt.Printf("Selecting %d domains from %s\n", N, filename)
	n := N / 3

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		domains = append(domains, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	total := len(domains)
	if total < N {
		return nil, fmt.Errorf("file must contain at least %d domains, got %d", N, total)
	}

	start := domains[:n]
	midStart := (total / 2) - 16
	mid := domains[midStart : midStart+n]
	end := domains[total-n:]

	selected := append(start, mid...)
	selected = append(selected, end...)

	return selected, nil
}

type Source struct {
	URL      string
	Category string
}

func TestDomainList(t *testing.T) {
	category := "phishing"        // Example category name
	fileName := "phishing-nl.txt" // Last part of the URL
	path := filepath.Join(dbstorePath, category, fileName)

	domainListSz := 190222
	domains, err := SelectDomains(domainListSz, path)
	if err != nil {
		t.Errorf("Error geting domains. err = %v", err)
		return
	}
	fmt.Println("List Size - ", len(domains))
	//fmt.Println("List     - ", domains)

	dbInstance := NewCategoryDB(false, false)

	sources := []Source{
		{
			URL:      "https://blocklistproject.github.io/Lists/alt-version/phishing-nl.txt",
			Category: "phishing",
		},
	}

	for _, src := range sources {
		fmt.Printf("Loading category '%s' from %s\n", src.Category, src.URL)
		err := dbInstance.LoadDomainsFromURL(dbstorePath, src.URL, src.Category)
		if err != nil {
			fmt.Printf("Error loading %s: %v\n", src.URL, err)
		}
	}

	// Now do lookup
	totalTimeNanos := int64(0)
	for _, domain := range domains {
		//fmt.Println("Domain - ", domain)
		cat, ok, nanos := dbInstance.Lookup(domain)
		if !ok {
			fmt.Println("Domain not found", domain)
		}
		totalTimeNanos += nanos
		_ = cat
		//fmt.Printf("Category: %s\tLookup Time: %d ns\n", cat, nanos)
	}
	avgNanos := totalTimeNanos / int64(len(domains))
	fmt.Printf("Average Nanos = %d\n", avgNanos)

	// Do one more time to use cache
	totalTimeNanos = int64(0)
	for _, domain := range domains {
		cat, ok, nanos := dbInstance.Lookup(domain)
		if !ok {
			fmt.Println("Domain not found", domain)
		}
		totalTimeNanos += nanos
		_ = cat
	}
	avgNanos = totalTimeNanos / int64(len(domains))
	fmt.Printf("Average Nanos = %d\n", avgNanos)
}

func TestGetTLD1(t *testing.T) {
	category := "facebook"
	fileName := "facebook-nl.txt"
	path := filepath.Join(dbstorePath, category, fileName)

	domainListSz := 22458
	domains, err := SelectDomains(domainListSz, path)
	if err != nil {
		t.Errorf("Error geting domains. err = %v", err)
		return
	}
	fmt.Println("List Size - ", len(domains))
	tldStore := dm.NewTLDStore()

	for _, domain := range domains {
		/*name, err := dm.GetTLD1(domain)
		if err != nil {
			fmt.Println("Error -", err)
			continue
		}
		fmt.Println(domain, name) */
		if err := tldStore.Add(domain); err != nil {
			fmt.Println("Err adding Domain ", domain, "err = ", err)
		}
	}

	list := tldStore.List()
	fmt.Println("TLD1s = ", list)
}

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("../config.json")
	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
	}

	dbstorePath = cfg.DBStorePath
	fmt.Println("DBStore Path = ", dbstorePath)

	e := m.Run()
	os.Exit(e)
}

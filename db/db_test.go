package db

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLookupRadixNormalizesDomain(t *testing.T) {
	db := NewCategoryDB(true, false)
	db.AddCategory(101, "phishing")
	db.radixTree.Insert("example.com", uint8(101))

	cat, found, _ := db.Lookup("  EXAMPLE.COM  ")
	if !found {
		t.Fatalf("expected radix lookup to find normalized domain")
	}
	if cat != "phishing" {
		t.Fatalf("expected category phishing, got %q", cat)
	}

	if cached, ok := db.Cache.Get("example.com"); !ok {
		t.Fatalf("expected normalized domain to be cached")
	} else if cached.(uint8) != 101 {
		t.Fatalf("expected cached category 101, got %v", cached)
	}
}

func TestLoadDomainsFromURLRewindsDownloadedFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Example.com\n#comment\n\nOther.org\n"))
	}))
	defer server.Close()

	storeDir := t.TempDir()
	db := NewCategoryDB(false, false)

	err, count, size := db.LoadDomainsFromURL(storeDir, server.URL+"/domains.txt", "phishing")
	if err != nil {
		t.Fatalf("LoadDomainsFromURL returned error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 domains to load, got %d", count)
	}
	if size != 2 {
		t.Fatalf("expected database size 2, got %d", size)
	}
	if got := db.FullDB["example.com"]; got != BaseCategoryID {
		t.Fatalf("expected normalized domain example.com to be stored, got category %d", got)
	}
	if got := db.FullDB["other.org"]; got != BaseCategoryID {
		t.Fatalf("expected normalized domain other.org to be stored, got category %d", got)
	}

	filePath := filepath.Join(storeDir, "phishing", "domains.txt")
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected downloaded file to exist: %v", err)
	}
}
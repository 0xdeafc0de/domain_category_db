package domains

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	//"golang.org/x/net/publicsuffix"
)

type TLDStore struct {
	mu   sync.RWMutex
	data map[string]struct{}
}

func NewTLDStore() *TLDStore {
	return &TLDStore{
		data: make(map[string]struct{}),
	}
}

func (store *TLDStore) Add(domain string) error {
	tld1, err := GetTLD1(domain)
	if err != nil {
		return err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[tld1] = struct{}{}
	return nil
}

func (store *TLDStore) List() []string {
	store.mu.RLock()
	defer store.mu.RUnlock()
	result := make([]string, 0, len(store.data))
	for tld := range store.data {
		result = append(result, tld)
	}
	sort.Strings(result)
	return result
}

func GetTLD1(domain string) (string, error) {
	// Normalize domain
	//domain = strings.ToLower(domain)

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid domain: %s", domain)
	}
	// Return the second last component (SLD)

	//eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(domain)
	//if err != nil {
	//	return "", err
	//}

	// Split and return the TLD-1 (e.g., facebook from facebook.com)
	//parts := strings.Split(eTLDPlusOne, ".")
	//if len(parts) > 1 {
	//	return parts[0], nil
	//}
	//return "", nil

	tld1 := parts[len(parts)-2]
	if tld1 == "m" {
		fmt.Println("Got m for ", domain, tld1)
	}
	// Return the second last component (SLD)
	return tld1, nil
}

package tld

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

var defaultTLDs = []string{
	// Common
	"com", "net", "org", "info", "biz", "jp", "co",
	// Tech & Startups
	"io", "dev", "app", "tech", "ai",
	// Individual
	"me", "work", "site", "link",
	// Other
	"moe", "xyz", "blue", "now", "cloud", "space", "online", "world",
}

type CheckRequest struct {
	Domain string `json:"domain"`
}

type CheckResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Result struct {
		Available bool   `json:"available"`
		Domain    string `json:"domain"`
	} `json:"result"`
}

func CheckCommonTLDs(baseDomain string) map[string]bool {
	results := make(map[string]bool)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, tld := range defaultTLDs {
		fullDomain := fmt.Sprintf("%s.%s", baseDomain, tld)
		url := fmt.Sprintf("https://rdap.org/domain/%s", fullDomain)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			slog.Error("Failed to create request", "domain", fullDomain, "error", err)
			continue
		}
		req.Header.Set("Accept", "application/rdap+json")

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("HTTP request failed", "domain", fullDomain, "error", err)
			continue
		}

		isAvailable := false
		if resp.StatusCode == http.StatusNotFound {
			isAvailable = true
		} else if resp.StatusCode != http.StatusOK {
			slog.Warn("Unexpected RDAP status", "domain", fullDomain, "status", resp.StatusCode)
		}
		resp.Body.Close()

		results[fullDomain] = isAvailable
		time.Sleep(1100 * time.Millisecond)
	}

	return results
}

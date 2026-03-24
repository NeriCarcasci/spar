package friends

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type PublicProfile struct {
	Username        string         `json:"username"`
	Rank            string         `json:"rank"`
	Division        int            `json:"division"`
	TotalSP         int            `json:"total_sp"`
	Streak          int            `json:"streak"`
	Solves          map[string]int `json:"solves"`
	TotalSolved     int            `json:"total_solved"`
	TotalChallenges int            `json:"total_challenges"`
	Languages       map[string]int `json:"languages"`
	TrackMedals     map[string]string `json:"track_medals"`
	LastUpdated     time.Time      `json:"last_updated"`
	SparVersion     string         `json:"spar_version"`
}

type SyncResult struct {
	Friend  Friend
	Profile *PublicProfile
	Error   error
	Status  string
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

func RawContentURL(f Friend) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/profile/profile.json", f.Username, f.RepoName)
}

func FetchProfile(f Friend) (PublicProfile, error) {
	rawURL := RawContentURL(f)
	resp, err := httpClient.Get(rawURL)
	if err != nil {
		return PublicProfile{}, fmt.Errorf("fetching profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return PublicProfile{}, fmt.Errorf("not_found")
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return PublicProfile{}, fmt.Errorf("rate_limited")
	}
	if resp.StatusCode != http.StatusOK {
		return PublicProfile{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return PublicProfile{}, fmt.Errorf("reading response: %w", err)
	}

	var p PublicProfile
	if err := json.Unmarshal(body, &p); err != nil {
		return PublicProfile{}, fmt.Errorf("parsing profile: %w", err)
	}
	return p, nil
}

func SyncAll(friends []Friend) []SyncResult {
	if len(friends) == 0 {
		return nil
	}

	results := make([]SyncResult, len(friends))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	var rateLimited bool
	var mu sync.Mutex

	for i, f := range friends {
		wg.Add(1)
		go func(idx int, fr Friend) {
			defer wg.Done()

			mu.Lock()
			if rateLimited {
				mu.Unlock()
				results[idx] = SyncResult{Friend: fr, Status: "rate_limited"}
				return
			}
			mu.Unlock()

			sem <- struct{}{}
			defer func() { <-sem }()

			profile, err := FetchProfile(fr)
			r := SyncResult{Friend: fr}
			if err != nil {
				errMsg := err.Error()
				switch errMsg {
				case "not_found":
					r.Status = "not_found"
				case "rate_limited":
					r.Status = "rate_limited"
					mu.Lock()
					rateLimited = true
					mu.Unlock()
				default:
					r.Status = "error"
				}
				r.Error = err
			} else {
				r.Profile = &profile
				r.Status = "ok"
			}
			results[idx] = r
		}(i, f)
	}

	wg.Wait()
	return results
}

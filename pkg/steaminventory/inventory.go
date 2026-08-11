package inventory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type UserInventory struct {
	Items             []Assets       `json:"assets"`
	ItemsDescriptions []Descriptions `json:"descriptions"`
	TotalItems        int            `json:"total_inventory_count"`
}

type Assets struct {
	ClassID string `json:"classid"`
}

type Descriptions struct {
	ClassID        string `json:"classid"`
	MarketHashName string `json:"market_hash_name"`
}

func createHTTPClient() (*http.Client, error) {

	proxyEnv := os.Getenv("PROXY_URL")

	if proxyEnv == "" {
		// fmt.Println("no proxy detected. Using default http client")
		return &http.Client{
			Timeout: 15 * time.Second,
		}, nil
	}

	proxyURL, err := url.Parse(proxyEnv)
	if err != nil {
		return nil, fmt.Errorf("proxy url invalid: %w", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	return &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}, nil
}

func GetUserInventory(steamID string, game, contextID int) (map[string]int, error) {

	fmt.Println("getting user inventory...")

	userInv := fmt.Sprintf("https://steamcommunity.com/inventory/%s/%d/%d?l=english&count=1000", steamID, game, contextID)

	client, err := createHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("error creating http client: %w", err)
	}

	res, err := client.Get(userInv)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error making request: %s", res.Status)
	}

	var assets UserInventory
	err = json.NewDecoder(res.Body).Decode(&assets)
	if err != nil {
		return nil, fmt.Errorf("error decoding json response: %w", err)
	}

	if assets.TotalItems == 0 {
		return nil, fmt.Errorf("probably 429 Too Many Requests!")
	}

	var assetMap = map[string]int{}
	for _, asset := range assets.Items {
		assetMap[asset.ClassID]++
	}

	var userItems = map[string]int{}
	for _, items := range assets.ItemsDescriptions {
		userItems[items.MarketHashName] = assetMap[items.ClassID]
	}

	return userItems, nil
}

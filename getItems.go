package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetUserInventory(steamID, game, contextID int) (map[string]int, error) {

	var repeatedSkins = map[string]int{}

	fmt.Println("Getting user inventory...")

	userInv := fmt.Sprintf("https://steamcommunity.com/inventory/%d/%d/%d?l=english&count=1000", steamID, game, contextID)

	res, err := http.Get(userInv)
	if err != nil {
		return repeatedSkins, err
	}
	defer res.Body.Close()

	// fmt.Printf("Request status: %s\n", res.Status)
	if res.StatusCode <= 199 || res.StatusCode >= 300 {
		return repeatedSkins, fmt.Errorf("Error making request: %s", res.Status)
	}

	// // Uncomment to see response as string
	// body, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	return repeatedSkins, err
	// }
	// fmt.Println(string(body))

	var assets UserInventory
	err = json.NewDecoder(res.Body).Decode(&assets)
	if err != nil {
		return repeatedSkins, err
	}

	// if total_inventory_count = 0 even though got status 200 you are rate limit
	if res.StatusCode == 200 && assets.TotalItems == 0 {
		return repeatedSkins, fmt.Errorf("Probably 429 Too Many Requests!")
	}

	var assetMap = map[string]int{}
	for _, asset := range assets.Items {

		_, ok := assetMap[asset.ClassID]
		if !ok {
			assetMap[asset.ClassID] = 1
		} else {
			assetMap[asset.ClassID]++
		}
	}

	for _, items := range assets.ItemsDescriptions {
		repeatedSkins[items.MarketHashName] = assetMap[items.ClassID]
	}

	return repeatedSkins, nil
}

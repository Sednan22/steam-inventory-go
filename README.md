# Steam Inventory Fetcher (Go)

A lightweight and efficient Go tool (and library) to fetch public Steam user inventories and aggregate item/skin counts by their market name (MarketHashName).

## 🚀 Features

- Fetches inventories for any Steam game (CS2, TF2, Rust, etc.) using a user's SteamID64.
- Maps and cross-references item IDs (assets) with their respective details (descriptions).
- Returns a map containing the market name of each item and its total quantity in the user's inventory.
- Modular project structure (pkg/ for the reusable library and cmd/ for the CLI tool).

## 🛠️ Installation & Building

Ensure you have Go (1.20+) and Make installed on your system.

Clone the repository:
git clone https://github.com/Sednan22/steam-inventory-go
cd steam-inventory-go

Build the Command Line Interface (CLI) executable using the Makefile:
make build

## 📖 Usage (CLI)

Run the compiled binary by passing the three required arguments:

./getUserInventory

### Arguments:
- SteamID64: The 64-bit numerical ID of the Steam user (profile inventory must be set to Public).
- AppID: The Steam Game ID (e.g., 730 for CS2, 440 for TF2).
- ContextID: The inventory context ID (typically 2 for most tradeable game items).

### Example:
./getUserInventory 76561198000000000 730 2

## ⚠️ Known Limitations (Steam API)

The public Steam Community Inventory API ([steamcommunity.com/inventory/](https://steamcommunity.com/inventory/)...) enforces strict rate limits:
- Making frequent requests in a short timeframe will lead to temporary IP blocks (HTTP 429 or HTTP 200 responses returning total_inventory_count: 0).
- For production or high-volume usage, implementing retry strategies with exponential backoff, proxy rotation, or response caching is strongly recommended.
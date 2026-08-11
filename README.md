# Steam Inventory Fetcher (Go)

A lightweight and efficient Go tool (and library) to fetch public Steam user inventories, aggregate item counts by market name, and export data in multiple formats.

## 🚀 Features

- Fetches inventories for any Steam game (CS2, TF2, Rust, etc.) using a user's SteamID64.
- Maps and cross-references item IDs (assets) with their respective details (descriptions).
- Aggregates items by market name (`market_hash_name`) and calculates total quantities.
- **Multiple Output Formats:** Formats data into **JSON**, **CSV**, or formatted plain text.
- **CLI Flags & File Export:** Output directly to the terminal or save results to a file with automatic format detection based on file extension.
- Modular project structure (`pkg/` for the reusable library and `cmd/` for the CLI tool).

## 🛠️ Installation & Building

Ensure you have Go (1.20+) and Make installed on your system.

Clone the repository:
```bash
git clone https://github.com/Sednan22/steam-inventory-go
cd steam-inventory-go
```

Build the Command Line Interface (CLI) executable using the Makefile:
```bash
make build
```

## 📖 Usage (CLI)

Run the compiled binary by passing the required positional arguments and optional flags:

`./getUserInventory [flags] <steamID64> <appID> <contextID>`

### Positional Arguments:
- **SteamID64:** The 64-bit numerical ID of the Steam user (profile inventory must be set to Public).
- **AppID:** The Steam Game ID (e.g., `730` for CS2, `440` for TF2).
- **ContextID:** The inventory context ID (typically `2` for most tradeable game items).

### Available Flags:
- `-format`: Explicitly specify the output format (`json`, `csv`, or `txt`).
- `-output`: Path to the output file where data should be saved.

---

### Examples

#### 1. Basic Terminal Output (Default Text Format)
```bash
./getUserInventory 76561198000000000 730 2
```

#### 2. Formatted Terminal Output
Print the inventory directly to the console formatted as JSON or CSV:
```bash
./getUserInventory -format=json 76561198000000000 730 2
./getUserInventory -format=csv 76561198000000000 440 2
```

#### 3. Export to File with Automatic Format Detection
When using `-output`, the format is automatically inferred from the file extension (`.json`, `.csv`, or `.txt`):
```bash
./getUserInventory -output=inventory.json 76561198000000000 730 2
./getUserInventory -output=inventory.csv 76561198000000000 730 2
```

#### 4. Export to File with Explicit Format
Override or specify the format when writing to custom file names:
```bash
./getUserInventory -format=json -output=my_inventory_backup 76561198000000000 730 2
```

## ⚠️ Known Limitations (Steam API)

The public Steam Community Inventory API (`steamcommunity.com/inventory/...`) enforces strict rate limits:
- **Rate Limits & IP Blocks:** Making frequent requests in a short timeframe will lead to temporary IP blocks (HTTP 429 or HTTP 200 responses returning `total_inventory_count: 0`). For production or high-volume usage, implementing retry strategies with exponential backoff, proxy rotation, or response caching is strongly recommended.
- **Context ID & Item Visibility Restrictions:**
  - Using `contextID = 2` only fetches **tradable items** accessible via the public endpoint without authentication.
  - Fetching non-tradable items, items on trade hold, or private context data (`contextID = 16`) requires authenticated requests with valid session cookies (`sessionid` and `steamLoginSecure`).
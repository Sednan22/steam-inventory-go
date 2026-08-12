# Steam Inventory Fetcher (Go)

A lightweight and efficient Go tool (and library) to fetch public Steam user inventories, aggregate item counts by market name, and export data in multiple formats.

## 🚀 Features

- Fetches inventories for any Steam game (CS2, TF2, Rust, etc.) using a user's SteamID64.
- Maps and cross-references item IDs (assets) with their respective details (descriptions).
- Aggregates items by market name (`market_hash_name`) and calculates total quantities.
- **Concurrent Multi-Fetch:** Fetches inventories for multiple users in parallel using **Goroutines** and **Channels**.
- **Built-in Rate Limiting:** Enforces request pacing via `time.Ticker` to mitigate API blocks during batch fetching.
- **Input Sanitization:** Automatically trims spaces, trailing control characters (`\r\n`), and ignores blank lines when reading user list files.
- **Multiple Output Formats:** Formats data into **JSON**, **CSV** (with dedicated `UserID` columns), or formatted plain text.
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

Run the compiled binary by passing positional arguments and optional flags. The tool supports two operational modes: **Single User** and **Multi-Fetch**.

### Available Flags:
- `-multi`: Enable batch fetching for multiple users listed in `usersList.txt` (default: `false`).
- `-format`: Explicitly specify the output format (`json`, `csv`, or `txt`).
- `-output`: Path to the output file where data should be saved.

---

### 1. Single User Mode

Fetch inventory data for a single Steam user.

**Usage:**
`./getUserInventory [flags] <steamID64> <appID> <contextID>`

#### Examples:
```bash
# Basic terminal output (Default Text Format)
./getUserInventory 76561198000000000 730 2

# Export directly to JSON or CSV terminal output
./getUserInventory -format=json 76561198000000000 730 2
./getUserInventory -format=csv 76561198000000000 440 2

# Export to file with automatic format detection based on extension
./getUserInventory -output=inventory.json 76561198000000000 730 2
./getUserInventory -output=inventory.csv 76561198000000000 730 2
```

---

### 2. Multi-Fetch Mode (Batch Processing)

Fetch inventories for multiple users listed in a `usersList.txt` file placed in the working directory. 

Each line in `usersList.txt` must contain a single `SteamID64`:
```text
76561198000000000
76561198111111111
```

**Usage:**
`./getUserInventory -multi=true [flags] <appID> <contextID>`

#### Examples:
```bash
# Fetch inventories for all IDs in usersList.txt and output to terminal
./getUserInventory -multi=true 730 2

# Batch fetch and export aggregated results into a single CSV file
./getUserInventory -multi=true -output=all_inventories.csv 730 2

# Batch fetch and export aggregated results into a formatted JSON file
./getUserInventory -multi=true -output=all_inventories.json 730 2
```

---

## 📊 Export Formats

| Format | Description | Structure Example |
| :--- | :--- | :--- |
| **JSON** | Aggregated nested object mapping UserIDs to item names and counts. | `{"76561198...": {"AK-47 \| Redline": 1}}` |
| **CSV** | Tabular format with dedicated columns for `UserID`, `Item`, and `Quantity`. | `UserID,Item,Quantity\n76561198...,"AK-47 \| Redline",1` |
| **TXT** *(Default)* | Human-readable text block grouped by user. | `76561198... inventory:\n - AK-47 -> 1` |

---

## ⚠️ Known Limitations (Steam API)

The public Steam Community Inventory API (`[steamcommunity.com/inventory/](https://steamcommunity.com/inventory/)...`) enforces strict rate limits:
- **Rate Limits & IP Blocks:** Making frequent requests in a short timeframe will lead to temporary IP blocks (HTTP 429 or HTTP 200 responses returning `total_inventory_count: 0`). Multi-fetch mode uses an internal rate-limiting ticker to space out requests, but large batches may still hit API thresholds.
- **Context ID & Item Visibility Restrictions:**
  - Using `contextID = 2` only fetches **tradable items** accessible via the public endpoint without authentication.
  - Fetching non-tradable items, items on trade hold, or private context data (`contextID = 16`) requires authenticated requests with valid session cookies (`sessionid` and `steamLoginSecure`).
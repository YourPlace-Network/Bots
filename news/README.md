# News Bot

Fetches RSS/Atom feeds and posts articles as on-chain transactions on Base.

## Setup

### config.json

Create a `config.json` in the `news/` directory:

```json
{
  "avatar": "https://example.com/avatar.png",
  "banner": "https://example.com/banner.png",
  "description": "A bot that posts news articles",
  "feeds": [
    "https://example.com/rss",
    "https://other-site.com/feed.xml"
  ],
  "maxPostLength": 500,
  "pollIntervalSeconds": 300,
  "rpcUrl": "https://mainnet.base.org",
  "username": "NewsBot",
  "vertical": "news"
}
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `avatar` | No | | Profile avatar URL (set via `-metadata` flag) |
| `banner` | No | | Profile banner URL (set via `-metadata` flag) |
| `description` | No | | Profile description (set via `-metadata` flag) |
| `feeds` | Yes | | List of RSS/Atom feed URLs |
| `maxPostLength` | No | `500` | Maximum character length per post |
| `pollIntervalSeconds` | No | `300` | Seconds between feed polling cycles |
| `rpcUrl` | Yes | | Base RPC endpoint URL |
| `username` | No | | Profile display name (set via `-metadata` flag) |
| `vertical` | No | | Profile category (set via `-metadata` flag) |

### wallet.json

The wallet file is stored at `Bots/wallet.json` (parent directory) and is shared across all bots.

On first run, if no `wallet.json` exists, the bot automatically generates a new Ethereum keypair and saves it. The bot prints the new wallet address to the console.

Fund this address with ETH on Base to pay for transaction gas fees.

The wallet file contains:

```json
{
  "address": "0x...",
  "privateKey": "0x...",
  "publicKey": "0x..."
}
```

## Usage

```sh
make              # continuous polling
make run_metadata # set profile metadata onchain and exit
make run_single   # post only the latest article from each feed, then exit
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `config.json` | Path to config file |
| `-data` | `data` | Path to data directory (stores dedup database) |
| `-metadata` | `false` | Set profile metadata onchain from config values, then exit |
| `-once` | `false` | Process all unposted articles once, then exit |
| `-single` | `false` | Post only the latest article from each feed, then exit |

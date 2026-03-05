# rtcheck

A command-line tool to check ABA routing numbers against real-time payment networks — [RTP](https://www.theclearinghouse.org/payment-systems/rtp) (The Clearing House) and [FedNow](https://www.frbservices.org/financial-services/fednow) (Federal Reserve).

`rtcheck` fetches the latest participant lists, caches them locally, and lets you look up individual routing numbers, bulk-analyze files, or browse the full directory of participants.

## Install

### Go

Requires Go 1.21+:

```sh
go install github.com/seb-chavez/rtcheck@latest
```

### Pre-built binaries

Download a binary for your platform from [GitHub Releases](https://github.com/seb-chavez/rtcheck/releases). Binaries are available for macOS (Intel & Apple Silicon), Linux, and Windows.

### Verify installation

```sh
rtcheck --version
```

## Quick start

```sh
# Check a single routing number
rtcheck lookup 021000021

# Bulk-analyze a file of routing numbers
rtcheck analyze payments.csv

# Browse all participants
rtcheck directory
```

## Commands

### `rtcheck lookup <routing-number>`

Check if a single routing number participates in RTP and/or FedNow.

```
$ rtcheck lookup 021000021
┌────────────────┬────────────────┐
│     FIELD      │     VALUE      │
├────────────────┼────────────────┤
│ Routing Number │ 021000021      │
│ Institution    │ JPMORGAN CHASE │
│ RTP            │ Yes            │
│ FedNow         │ Yes            │
└────────────────┴────────────────┘
```

JSON output:

```
$ rtcheck lookup 021000021 --format json
{
  "routing_number": "021000021",
  "institution": "JPMORGAN CHASE",
  "rtp": true,
  "fednow": true
}
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format: `table`, `json` | `table` |

**Exit codes:**

| Code | Meaning |
|------|---------|
| `0` | Routing number participates in at least one real-time network |
| `1` | Invalid routing number or error |
| `2` | Valid routing number but not on any real-time network |

The exit code behavior makes `rtcheck` scriptable — you can use it in shell pipelines or CI checks:

```sh
if rtcheck lookup 021000021 > /dev/null 2>&1; then
  echo "Supports real-time payments"
fi
```

### `rtcheck analyze <file>`

Bulk-check a file of routing numbers against both networks. Deduplicates routing numbers and provides aggregate statistics.

```
$ rtcheck analyze payments.csv
```

By default, `analyze` prints a summary table. Use `--no-summary` to output the per-routing-number detail rows instead. Use `-o` to write detailed results to a CSV file alongside the summary.

```sh
# Summary only (default)
rtcheck analyze payments.csv

# Detail rows only
rtcheck analyze payments.csv --no-summary

# Summary on screen, details to file
rtcheck analyze payments.csv -o results.csv

# Machine-readable output
rtcheck analyze payments.csv --format json
rtcheck analyze payments.csv --format csv
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format: `table`, `json`, `csv` | `table` |
| `--no-summary` | Skip summary, output detail rows only | `false` |
| `-o`, `--output` | Write detailed per-RTN results to a CSV file | — |

**Supported file formats:**

| Extension | Format |
|-----------|--------|
| `.csv` | Comma-separated values |
| `.xlsx`, `.xls` | Excel spreadsheets (reads the first sheet) |
| `.txt`, `.tsv`, or other | Plain text (one routing number per line) or auto-detected CSV/TSV |

The routing number column is auto-detected by matching common header names (`routing_number`, `rtn`, `aba`, `routing`, `transit_number`, etc.). If no header matches, `rtcheck` falls back to finding the column where >50% of values are valid 9-digit routing numbers.

### `rtcheck directory`

Browse and search all routing numbers participating in RTP and/or FedNow. Results are paginated interactively in table mode.

```sh
# Browse all participants (paginated)
rtcheck directory

# Search by institution name
rtcheck directory --search "chase"

# Search by routing number prefix
rtcheck directory --search "0210"

# Filter by network
rtcheck directory --network rtp
rtcheck directory --network fednow
rtcheck directory --network both

# Export full directory to CSV
rtcheck directory --format csv > participants.csv

# Combine filters
rtcheck directory --network fednow --search "wells" --format json
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format: `table`, `json`, `csv` | `table` |
| `--search` | Filter by institution name (case-insensitive) or routing number prefix | — |
| `--network` | Filter by network: `rtp`, `fednow`, `both` | — (show all) |

In table mode, results are displayed 50 at a time. Use `n`/`p`/`q` to navigate pages.

## Global flags

These flags apply to all commands:

| Flag | Description | Default |
|------|-------------|---------|
| `--refresh` | Force re-download of participant data (ignore cache) | `false` |
| `--cache-dir` | Override cache directory | `~/.rtcheck/data/` |
| `-v`, `--version` | Print version | — |
| `-h`, `--help` | Print help | — |

## Caching

`rtcheck` caches downloaded participant data locally at `~/.rtcheck/data/` with a **24-hour TTL**. This avoids hitting external sources on every invocation.

- Data is automatically refreshed when the cache expires
- Use `--refresh` to force a re-download at any time
- Use `--cache-dir` to store cache files in a custom location

## Data sources

| Network | Source | Approx. size |
|---------|--------|--------------|
| **RTP** | [The Clearing House](https://www.theclearinghouse.org/payment-systems/rtp/rtn) | ~2,100 routing numbers |
| **FedNow** | [Federal Reserve](https://www.frbservices.org/financial-services/fednow/organizations) | ~1,980 routing numbers |
| **Institution names** | [moov-io/fed](https://github.com/moov-io/fed) FedACH dictionary | — |

## License

MIT

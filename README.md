# rtcheck

Check your routing numbers against real-time payment networks (RTP & FedNow).

`rtcheck` fetches the latest participant lists from [The Clearing House](https://www.theclearinghouse.org/payment-systems/rtp/rtn) (RTP) and the [Federal Reserve](https://www.frbservices.org/financial-services/fednow/organizations) (FedNow), then checks your routing numbers against them.

## Install

**Homebrew:**

    brew install seb-chavez/tap/rtcheck

**Go:**

    go install github.com/seb-chavez/rtcheck@latest

**Binary:** Download from [GitHub Releases](https://github.com/seb-chavez/rtcheck/releases).

## Usage

### Look up a single routing number

    rtcheck lookup 021000021

### Analyze a file of routing numbers

    rtcheck analyze payments.csv

Supports CSV, TSV, TXT, and Excel (.xlsx) files. Auto-detects the routing number column.

### Browse all participants

    rtcheck directory
    rtcheck directory --search "chase"
    rtcheck directory --network rtp

### Output formats

All commands support `--format json` and `--format csv` for machine-readable output:

    rtcheck lookup 021000021 --format json
    rtcheck analyze payments.csv --format json
    rtcheck directory --format csv > all-participants.csv

## Data Sources

- **RTP:** [The Clearing House](https://www.theclearinghouse.org/payment-systems/rtp/rtn) (~2,100 routing numbers)
- **FedNow:** [Federal Reserve](https://www.frbservices.org/financial-services/fednow/organizations) (~1,980 routing numbers)
- **Institution Names:** [moov-io/fed](https://github.com/moov-io/fed) FedACH dictionary

Data is cached locally at `~/.rtcheck/data/` with a 24-hour TTL. Use `--refresh` to force re-download.

## License

MIT

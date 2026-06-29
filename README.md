<img width="1200" alt="Invoice" src="https://github.com/maaslalani/nap/assets/42545625/16dae9d9-390c-49b6-aedd-3f882b17f57b">

# Invoice

Generate invoices from the command line.

## Command Line Interface

```bash
invoice generate --from "Dream, Inc." --to "Imagine, Inc." \
    --item "Rubber Duck" --quantity 2 --rate 25 \
    --tax 0.13 --discount 0.15 \
    --note "For debugging purposes."
```

<img src="https://vhs.charm.sh/vhs-66CMd4UQuXkuxX9djHUnGX.gif" width="600" />

View the generated PDF at `invoice.pdf`, you can customize the output location
with `--output`.

```bash
open invoice.pdf
```

<img width="574" alt="Example invoice" src="https://github.com/maaslalani/nap/assets/42545625/13153de2-dfa1-41e6-a18e-4d3a5cea5b74">

### Environment

Save repeated information with environment variables:

```bash
export INVOICE_LOGO=/path/to/image.png
export INVOICE_FROM="Dream, Inc."
export INVOICE_TO="Imagine, Inc."
export INVOICE_TAX=0.13
export INVOICE_RATE=25
```

Generate new invoice:

```bash
invoice generate \
    --item "Yellow Rubber Duck" --quantity 5 \
    --item "Special Edition Plaid Rubber Duck" --quantity 1 \
    --note "For debugging purposes." \
    --output duck-invoice.pdf
```

### Configuration File

Or, save repeated information with JSON / YAML:

```json
{
    "logo": "/path/to/image.png",
    "from": "Dream, Inc.",
    "to": "Imagine, Inc.",
    "tax": 0.13,
    "items": ["Yellow Rubber Duck", "Special Edition Plaid Rubber Duck"],
    "quantities": [5, 1],
    "rates": [25, 25],
}
```

Generate new invoice by importing the configuration file:

```bash
invoice generate --import path/to/data.json \
    --output duck-invoice.pdf
```

### Metadata

Add structured key-value data to your invoice with `--metadata` (or `-m`). Each
entry uses `key=value` format:

```bash
invoice generate --from "Dream, Inc." --to "Imagine, Inc." \
    --item "Rubber Duck" --quantity 2 --rate 25 \
    --metadata "Bank=HSBC" \
    --metadata "IBAN=DE89 3704 0044 0532 0130 00" \
    --metadata "Tax ID=DE123456789"
```

Metadata entries are rendered in a **DETAILS** section on the PDF, sorted
alphabetically by key.

In JSON/YAML config files, metadata is a plain object:

```json
{
    "from": "Dream, Inc.",
    "to": "Imagine, Inc.",
    "metadata": {
        "Bank": "HSBC",
        "IBAN": "DE89 3704 0044 0532 0130 00",
        "Tax ID": "DE123456789"
    }
}
```

#### Note templates

Metadata values can be referenced in `--note` using `${Key}` placeholders.
This lets you keep values in one place and reuse them in the note text:

```bash
invoice generate --from "Dream, Inc." --to "Imagine, Inc." \
    --item "Rubber Duck" --quantity 2 --rate 25 \
    --metadata "Bank=HSBC" \
    --metadata "IBAN=DE89 3704 0044 0532 0130 00" \
    --note "Please transfer to ${Bank}, IBAN: ${IBAN}"
```

The note renders as: *Please transfer to HSBC, IBAN: DE89 3704 0044 0532 0130 00*

### Custom Templates

If you would like a custom invoice template for your business or company, please
reach out via:

* [Email](mailto:maas@lalani.dev)
* [Twitter](https://twitter.com/maaslalani)

## Installation

### Homebrew (macOS / Linux)

```sh
brew install jpdarago/tap/invoice
```

No Go toolchain required — this installs a prebuilt binary.

### From source

```sh
go install github.com/maaslalani/invoice@main
```

Or download a binary from the [releases](https://github.com/jpdarago/invoice/releases).

### Verifying a release

Each release tag is GPG-signed, so you can confirm it came from the maintainer:

```sh
git clone https://github.com/jpdarago/invoice
cd invoice
git verify-tag v1.0.0    # replace with the release tag
```

GitHub also displays such tags as "Verified". Homebrew independently verifies
the SHA-256 of every downloaded archive against the formula over HTTPS.

### Docker

Build the image:

```bash
docker build -t invoice .
```

Run with Docker:

```bash
docker run --rm -v "$PWD:/out" invoice generate \
    --from "Dream, Inc." --to "Imagine, Inc." \
    --item "Rubber Duck" --quantity 2 --rate 25 \
    --output /out/invoice.pdf
```

The `-v "$PWD:/out"` mount makes the generated PDF available on your host.

## Development

### Test

```bash
go test ./...
```

### Build

```bash
go build
```

### Lint

```bash
golangci-lint run
```

## License

[MIT](https://github.com/maaslalani/invoice/blob/master/LICENSE)

## Feedback

I'd love to hear your feedback on improving `invoice`.

Feel free to reach out via:
* [Email](mailto:maas@lalani.dev)
* [Twitter](https://twitter.com/maaslalani)
* [GitHub issues](https://github.com/maaslalani/invoice/issues/new)

---

<sub><sub>z</sub></sub><sub>z</sub>z

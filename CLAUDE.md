# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CLI tool that generates PDF invoices from command-line flags, environment variables, or JSON/YAML config files. Written in Go (~470 lines across 4 source files), all in the `main` package.

## Build & Run

```bash
go build                              # build binary
go install github.com/maaslalani/invoice@main  # install
./invoice generate --from "Acme" --to "Client" --item "Widget" --quantity 2 --rate 50
```

## Test

```bash
go test ./...
```

## Lint

```bash
golangci-lint run
```

## Architecture

| File | Purpose |
|------|---------|
| `main.go` | Entry point, `Invoice` struct, Cobra command/flag definitions, PDF generation orchestration |
| `pdf.go` | All PDF rendering functions (layout, fonts, rows, totals, footer) using `gopdf` |
| `import.go` | JSON/YAML file import with CLI flag override merging |
| `currency.go` | Currency-to-symbol mapping (`currencySymbols` map) |
| `Inter/` | Embedded TTF fonts (Inter regular + bold) via `go:embed` |

**Data flow:** CLI flags/env vars (Viper) → optional JSON/YAML import → `Invoice` struct → PDF rendering → `.pdf` file output.

**Key dependencies:** `cobra` (CLI), `viper` (config/env), `gopdf` (PDF generation), `yaml.v3` (YAML parsing).

**Configuration layering:** CLI flags override imported file values. Environment variables use `INVOICE_` prefix (e.g., `INVOICE_FROM`, `INVOICE_TAX`).

**Multiline fields:** `from`, `to`, and `note` support `\n` literal strings for line breaks in the PDF output.

**PDF layout:** Fixed column offsets defined as constants in `pdf.go` (quantity at 360, rate at 405, amount at 480). Page size is A4 with 40px margins. Totals and notes are positioned at Y=600, footer at Y=800.

# bragg-xrd

bragg-xrd is a Go Bragg diffraction calculator. It solves `n*lambda =
2*d*sin(theta)` for given wavelength, spacing (or cubic lattice constant plus
Miller indices), and order, returning theta and 2-theta. For cubic lattices it
uses `d = a/sqrt(h^2+k^2+l^2)`, applies body-centered or face-centered
systematic absences, and can enumerate a powder pattern of allowed peaks. The
service is available through HTTP JSON endpoints and CLI subcommands with no
web page.

## Usage

Run the HTTP server:

```bash
go run . serve -addr :8080
```

Evaluate from the command line:

```bash
go run . bragg -lambda 1.5406 -d 2.087 -n 1
go run . bragg -lambda 1.5406 -a 3.615 -hkl 111 -lattice fcc
go run . powder -lambda 1.5406 -a 3.615 -lattice fcc -max 4
```

Run the copper example:

```bash
go run . example -file example/fcc-cu-cu.json
```

Cu K-alpha (`1.5406 A`) on Cu FCC `(111)` gives `2-theta` near 43.3 degrees.

## HTTP API

```text
POST /api/bragg   {"lambda":1.5406,"d":2.087,"n":1}
POST /api/bragg   {"lambda":1.5406,"a":3.615,"hkl":[1,1,1]}
POST /api/powder  {"lambda":1.5406,"a":3.615,"lattice":"fcc"}
GET  /health
```

Invalid wavelength, spacing, cell constant, order, or lattice type return an
error body with HTTP 400. Orders with `sin(theta) > 1` return
`"possible": false` instead of a fake angle.

## Systematic Absences

- Primitive/simple: all nonzero hkl allowed
- BCC: `h+k+l` odd forbidden
- FCC: h,k,l must be all even or all odd

## Code Layout

```text
internal/xrd       Bragg equation, cubic spacing, structure factors, powder
internal/patterns  normalized peak tables, CSV and summaries
internal/verify    sin-theta gate, forbidden filtering, trend checks
internal/server    HTTP handlers and JSON responses
internal/cli       subcommand parsing and terminal output
example/           offline scenario JSON files
```

## Build and Test

```bash
export GOTOOLCHAIN=local CGO_ENABLED=0
go build ./...
go test ./...
```

The Dockerfile builds the server binary and starts it on port 8080.

package cli

import (
	"fmt"
	"os"

	"bragg-xrd/internal/xrd"
)

func Run(args []string) int {
	if len(args) == 0 {
		return runServe([]string{})
	}
	switch args[0] {
	case "bragg":
		return runBragg(args[1:])
	case "powder":
		return runPowder(args[1:])
	case "example":
		return runExample(args[1:])
	case "serve":
		return runServe(args[1:])
	case "help", "-h", "--help":
		printHelp()
		return 0
	case "version":
		fmt.Println("bragg-xrd 1.0.0")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Println(`bragg-xrd: Bragg diffraction angle calculator

Usage:
  bragg-xrd                          start HTTP server on :8080
  bragg-xrd bragg -lambda 1.5406 -d 2.087 -n 1
  bragg-xrd bragg -lambda 1.5406 -a 3.615 -hkl 111 -lattice fcc
  bragg-xrd powder -lambda 1.5406 -a 3.615 -lattice fcc -max 4
  bragg-xrd example -file example/fcc-cu-cu.json
  bragg-xrd serve -addr :8080

HTTP:
  POST /api/bragg   {"lambda":1.5406,"d":2.087,"n":1}
  POST /api/powder  {"lambda":1.5406,"a":3.615,"lattice":"fcc"}`)
}

func fail(err error) int {
	fmt.Fprintln(os.Stderr, "error:", err)
	return 1
}

func runBragg(args []string) int {
	fs := flagSet("bragg")
	lambda := fs.Float64("lambda", 1.5406, "wavelength")
	d := fs.Float64("d", 0, "spacing")
	a := fs.Float64("a", 0, "lattice constant")
	n := fs.Int("n", 1, "order")
	hklText := fs.String("hkl", "111", "miller indices")
	lattice := fs.String("lattice", "primitive", "lattice type")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if err := xrd.ValidateLattice(*lattice); err != nil {
		return fail(err)
	}
	if *d <= 0 && *a > 0 {
		hkl, err := xrd.MillersFromText(*hklText)
		if err != nil {
			return fail(err)
		}
		dValue, err := xrd.LatticeSpacing(*a, hkl)
		if err != nil {
			return fail(err)
		}
		d = &dValue
	}
	result, err := xrd.BraggAngle(*lambda, *d, *n)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runPowder(args []string) int {
	fs := flagSet("powder")
	lambda := fs.Float64("lambda", 1.5406, "wavelength")
	a := fs.Float64("a", 3.615, "lattice constant")
	lattice := fs.String("lattice", "fcc", "lattice type")
	maxIndex := fs.Int("max", 4, "max miller index")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	result, err := xrd.Powder(*lambda, *a, *lattice, *maxIndex)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runExample(args []string) int {
	fs := flagSet("example")
	file := fs.String("file", "example/fcc-cu-cu.json", "scenario JSON file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	scenario, err := xrd.LoadScenario(*file)
	if err != nil {
		return fail(err)
	}
	result, err := xrd.RunScenario(scenario)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func loadExample(path string) (xrd.Scenario, error) {
	return xrd.LoadScenario(path)
}

func ExamplePaths() []string {
	return xrd.ScenarioPaths()
}

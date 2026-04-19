package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ccvass/swarmex/swarmex-pack"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	packFile := "swarmex-pack.yml"
	overrides := make(map[string]string)

	// Parse --set and --pack-file from args
	args := os.Args[2:]
	var positional []string
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--set" && i+1 < len(args):
			i++
			parts := strings.SplitN(args[i], "=", 2)
			if len(parts) == 2 {
				overrides[parts[0]] = parts[1]
			}
		case args[i] == "--pack-file" && i+1 < len(args):
			i++
			packFile = args[i]
		default:
			positional = append(positional, args[i])
		}
	}

	switch os.Args[1] {
	case "install", "upgrade":
		if len(positional) < 1 {
			fmt.Fprintln(os.Stderr, "usage: swarmex-pack install <name> [--set key=value]")
			os.Exit(1)
		}
		cfg, err := pack.LoadPack(packFile)
		fatal(err)
		rendered, err := pack.Render(cfg, overrides)
		fatal(err)
		fatal(pack.Install(positional[0], rendered))
		fmt.Printf("✓ %s installed\n", positional[0])

	case "uninstall":
		if len(positional) < 1 {
			fmt.Fprintln(os.Stderr, "usage: swarmex-pack uninstall <name>")
			os.Exit(1)
		}
		fatal(pack.Uninstall(positional[0]))
		fmt.Printf("✓ %s removed\n", positional[0])

	case "render":
		cfg, err := pack.LoadPack(packFile)
		fatal(err)
		rendered, err := pack.Render(cfg, overrides)
		fatal(err)
		os.Stdout.Write(rendered)

	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `swarmex-pack — Helm-like packaging for Docker Swarm

Commands:
  install <name>   Deploy a stack from pack template
  upgrade <name>   Update an existing stack (same as install)
  uninstall <name> Remove a stack
  render           Print rendered YAML to stdout

Flags:
  --set key=value  Override a template value
  --pack-file      Path to pack file (default: swarmex-pack.yml)`)
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

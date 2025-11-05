package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	count      bool
	duplicate  bool
	unique     bool
	ignoreCase bool
	numFields  int
	numChars   int
	inputFile  string
	outputFile string
}

type LineInfo struct {
	original string
	key      string
	count    int
}

func main() {
	config, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]\n")
		os.Exit(1)
	}

	if err := runUniq(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() (*Config, error) {
	config := &Config{}

	flag.BoolVar(&config.count, "c", false, "")
	flag.BoolVar(&config.duplicate, "d", false, "")
	flag.BoolVar(&config.unique, "u", false, "")
	flag.BoolVar(&config.ignoreCase, "i", false, "")
	flag.IntVar(&config.numFields, "f", 0, "")
	flag.IntVar(&config.numChars, "s", 0, "")

	flag.Parse()

	mutualExclusive := 0
	if config.count {
		mutualExclusive++
	}
	if config.duplicate {
		mutualExclusive++
	}
	if config.unique {
		mutualExclusive++
	}
	if mutualExclusive > 1 {
		return nil, fmt.Errorf("parameters -c, -d, -u are mutually exclusive")
	}

	args := flag.Args()
	switch len(args) {
	case 0:
	case 1:
		config.inputFile = args[0]
	case 2:
		config.inputFile = args[0]
		config.outputFile = args[1]
	default:
		return nil, fmt.Errorf("too many arguments")
	}

	if config.numFields < 0 {
		return nil, fmt.Errorf("num_fields must be non-negative")
	}
	if config.numChars < 0 {
		return nil, fmt.Errorf("num_chars must be non-negative")
	}

	return config, nil
}

func runUniq(config *Config) error {
	var input io.Reader
	if config.inputFile != "" {
		file, err := os.Open(config.inputFile)
		if err != nil {
			return fmt.Errorf("open input file: %w", err)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	var output io.Writer
	if config.outputFile != "" {
		file, err := os.Create(config.outputFile)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer file.Close()
		output = file
	} else {
		output = os.Stdout
	}

	lines, err := readAndProcessLines(input, config)
	if err != nil {
		return fmt.Errorf("read lines: %w", err)
	}

	return writeResults(output, lines, config)
}

func readAndProcessLines(input io.Reader, config *Config) ([]*LineInfo, error) {
	scanner := bufio.NewScanner(input)
	lineMap := make(map[string]*LineInfo)
	var order []string

	for scanner.Scan() {
		original := scanner.Text()
		key := processLine(original, config)

		if info, exists := lineMap[key]; exists {
			info.count++
		} else {
			lineMap[key] = &LineInfo{
				original: original,
				key:      key,
				count:    1,
			}
			order = append(order, key)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	result := make([]*LineInfo, 0, len(order))
	for _, key := range order {
		result = append(result, lineMap[key])
	}

	return result, nil
}

func processLine(line string, config *Config) string {
	processed := line

	if config.numFields > 0 {
		fields := strings.Fields(processed)
		if len(fields) > config.numFields {
			processed = strings.Join(fields[config.numFields:], " ")
		} else {
			processed = ""
		}
	}

	if config.numChars > 0 {
		if len(processed) > config.numChars {
			processed = processed[config.numChars:]
		} else {
			processed = ""
		}
	}

	if config.ignoreCase {
		processed = strings.ToLower(processed)
	}

	return processed
}

func writeResults(output io.Writer, lines []*LineInfo, config *Config) error {
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for _, info := range lines {
		switch {
		case config.count:
			if _, err := fmt.Fprintf(writer, "%d %s\n", info.count, info.original); err != nil {
				return err
			}
		case config.duplicate:
			if info.count > 1 {
				if _, err := fmt.Fprintf(writer, "%s\n", info.original); err != nil {
					return err
				}
			}
		case config.unique:
			if info.count == 1 {
				if _, err := fmt.Fprintf(writer, "%s\n", info.original); err != nil {
					return err
				}
			}
		default:
			if _, err := fmt.Fprintf(writer, "%s\n", info.original); err != nil {
				return err
			}
		}
	}

	return nil
}

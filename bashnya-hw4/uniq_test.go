package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantConfig    *Config
		wantError     bool
		errorContains string
	}{
		{
			name:       "no flags",
			args:       []string{},
			wantConfig: &Config{},
		},
		{
			name:       "count flag",
			args:       []string{"-c"},
			wantConfig: &Config{count: true},
		},
		{
			name:       "duplicate flag",
			args:       []string{"-d"},
			wantConfig: &Config{duplicate: true},
		},
		{
			name:       "unique flag",
			args:       []string{"-u"},
			wantConfig: &Config{unique: true},
		},
		{
			name:       "ignore case flag",
			args:       []string{"-i"},
			wantConfig: &Config{ignoreCase: true},
		},
		{
			name:       "num fields flag",
			args:       []string{"-f", "2"},
			wantConfig: &Config{numFields: 2},
		},
		{
			name:       "num chars flag",
			args:       []string{"-s", "5"},
			wantConfig: &Config{numChars: 5},
		},
		{
			name:       "input file only",
			args:       []string{"input.txt"},
			wantConfig: &Config{inputFile: "input.txt"},
		},
		{
			name:       "input and output files",
			args:       []string{"input.txt", "output.txt"},
			wantConfig: &Config{inputFile: "input.txt", outputFile: "output.txt"},
		},
		{
			name: "multiple flags",
			args: []string{"-i", "-f", "2", "-s", "3", "input.txt"},
			wantConfig: &Config{
				ignoreCase: true,
				numFields:  2,
				numChars:   3,
				inputFile:  "input.txt",
			},
		},
		{
			name:          "mutually exclusive flags cd",
			args:          []string{"-c", "-d"},
			wantError:     true,
			errorContains: "mutually exclusive",
		},
		{
			name:          "mutually exclusive flags cu",
			args:          []string{"-c", "-u"},
			wantError:     true,
			errorContains: "mutually exclusive",
		},
		{
			name:          "mutually exclusive flags du",
			args:          []string{"-d", "-u"},
			wantError:     true,
			errorContains: "mutually exclusive",
		},
		{
			name:          "negative num fields",
			args:          []string{"-f", "-1"},
			wantError:     true,
			errorContains: "non-negative",
		},
		{
			name:          "negative num chars",
			args:          []string{"-s", "-5"},
			wantError:     true,
			errorContains: "non-negative",
		},
		{
			name:          "too many arguments",
			args:          []string{"in.txt", "out.txt", "extra.txt"},
			wantError:     true,
			errorContains: "too many arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = append([]string{"uniq"}, tt.args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			config, err := parseFlags()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseFlags() expected error, got nil")
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("parseFlags() error = %v, should contain %v", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseFlags() unexpected error: %v", err)
				return
			}

			if config.count != tt.wantConfig.count ||
				config.duplicate != tt.wantConfig.duplicate ||
				config.unique != tt.wantConfig.unique ||
				config.ignoreCase != tt.wantConfig.ignoreCase ||
				config.numFields != tt.wantConfig.numFields ||
				config.numChars != tt.wantConfig.numChars ||
				config.inputFile != tt.wantConfig.inputFile ||
				config.outputFile != tt.wantConfig.outputFile {
				t.Errorf("parseFlags() = %+v, want %+v", config, tt.wantConfig)
			}
		})
	}
}

func TestProcessLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		config   *Config
		expected string
	}{
		{
			name:     "no processing",
			line:     "Hello World",
			config:   &Config{},
			expected: "Hello World",
		},
		{
			name:     "ignore case",
			line:     "Hello World",
			config:   &Config{ignoreCase: true},
			expected: "hello world",
		},
		{
			name:     "ignore fields within range",
			line:     "field1 field2 field3 field4",
			config:   &Config{numFields: 2},
			expected: "field3 field4",
		},
		{
			name:     "ignore fields beyond range",
			line:     "field1 field2",
			config:   &Config{numFields: 3},
			expected: "",
		},
		{
			name:     "ignore chars within range",
			line:     "Hello World",
			config:   &Config{numChars: 6},
			expected: "World",
		},
		{
			name:     "ignore chars beyond range",
			line:     "Hello",
			config:   &Config{numChars: 10},
			expected: "",
		},
		{
			name:     "fields then chars",
			line:     "a b c d e f",
			config:   &Config{numFields: 2, numChars: 2},
			expected: "c d e f",
		},
		{
			name:     "case with fields",
			line:     "A B C D",
			config:   &Config{ignoreCase: true, numFields: 1},
			expected: "b c d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processLine(tt.line, tt.config)
			if result != tt.expected {
				t.Errorf("processLine(%q) = %q, want %q", tt.line, result, tt.expected)
			}
		})
	}
}

func TestUniqIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		config   *Config
		expected string
	}{
		{
			name:     "basic unique",
			input:    "a\nb\na\nc\nb\na\n",
			config:   &Config{},
			expected: "a\nb\nc\n",
		},
		{
			name:     "count mode",
			input:    "a\nb\na\nc\nb\na\n",
			config:   &Config{count: true},
			expected: "3 a\n2 b\n1 c\n",
		},
		{
			name:     "duplicate mode",
			input:    "a\nb\na\nc\nb\na\nd\n",
			config:   &Config{duplicate: true},
			expected: "a\nb\n",
		},
		{
			name:     "unique mode",
			input:    "a\nb\na\nc\nb\na\nd\n",
			config:   &Config{unique: true},
			expected: "c\nd\n",
		},
		{
			name:     "ignore case",
			input:    "A\nb\na\nB\nc\n",
			config:   &Config{ignoreCase: true},
			expected: "A\nb\nc\n",
		},
		{
			name:     "ignore fields",
			input:    "1 a\n2 a\n3 b\n1 c\n",
			config:   &Config{numFields: 1},
			expected: "1 a\n3 b\n1 c\n",
		},
		{
			name:     "ignore chars",
			input:    "abc\ndef\nadc\n",
			config:   &Config{numChars: 1},
			expected: "abc\ndef\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			var output bytes.Buffer

			lines, err := readAndProcessLines(input, tt.config)
			if err != nil {
				t.Fatalf("readAndProcessLines failed: %v", err)
			}

			err = writeResults(&output, lines, tt.config)
			if err != nil {
				t.Fatalf("writeResults failed: %v", err)
			}

			if output.String() != tt.expected {
				t.Errorf("Output = %q, want %q", output.String(), tt.expected)
			}
		})
	}
}

func TestEmptyInput(t *testing.T) {
	input := strings.NewReader("")
	var output bytes.Buffer

	config := &Config{}
	lines, err := readAndProcessLines(input, config)
	if err != nil {
		t.Fatalf("readAndProcessLines failed: %v", err)
	}

	err = writeResults(&output, lines, config)
	if err != nil {
		t.Fatalf("writeResults failed: %v", err)
	}

	if output.String() != "" {
		t.Errorf("Output should be empty, got %q", output.String())
	}
}

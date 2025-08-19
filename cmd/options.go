package main

import (
	"errors"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
)

type Options struct {
	NatsURL         string
	CredentialsFile string
	PostgresURL     string
	SchemaDir       string
}

var (
	ErrBadNatsURL         = errors.New("Invalid NATS URL")
	ErrBadCredentialsFile = errors.New("Invalid credentials file path")
	ErrBadPostgresURL     = errors.New("Invalid Postgres URL")
	ErrBadSchemaDir       = errors.New("Invalid schema directory path")
)

// Parses command line arguments from os.Args[1:] and returns an Options struct
func GetOptions() (Options, error) {
	return ParseOptions(os.Args[1:])
}

// ParseOptions parses the provided arguments and returns an Options struct
func ParseOptions(args []string) (Options, error) {
	// Create a new FlagSet for parsing the args
	flagset := flag.NewFlagSet("beaker", flag.ContinueOnError)

	// Set default options
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Options{}, err
	}
	workingDir, err := os.Getwd()
	if err != nil {
		return Options{}, err
	}
	options := Options{
		CredentialsFile: filepath.Join(homeDir, "NATS_CREDS_APP.creds"),
		PostgresURL:     "postgres://postgres:password@localhost:5432/beaker_development?sslmode=disable",
		NatsURL:         "tls://connect.ngs.global",
		SchemaDir:       filepath.Join(workingDir, "schemas"),
	}

	// Use flags to parse command line arguments
	flagset.StringVar(&options.CredentialsFile, "credentials", options.CredentialsFile, "Path to the NATS credentials file. See https://docs.nats.io/nats-concepts/security/ for details")
	flagset.StringVar(&options.PostgresURL, "postgres", options.PostgresURL, "Postgres connection URL. See https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING for details")
	flagset.StringVar(&options.NatsURL, "nats", options.NatsURL, "NATS server URL. See https://docs.nats.io/nats-concepts/nats-server/ for details")
	flagset.StringVar(&options.SchemaDir, "schema", options.SchemaDir, "Path to the JSON schema directory")
	// Add help flag
	flagset.Bool("help", false, "Show help message")

	// Parse the provided arguments
	if err := flagset.Parse(args); err != nil {
		return Options{}, err
	}

	// If help flag is set, print usage and exit
	if flagset.Parsed() && flagset.Lookup("help").Value.String() == "true" {
		flagset.Usage()
		os.Exit(0)
	}

	// Validate the credentials file path
	if options.CredentialsFile == "" {
		return Options{}, ErrBadCredentialsFile
	}
	if _, err := os.Stat(options.CredentialsFile); errors.Is(err, fs.ErrNotExist) {
		return Options{}, err
	}
	// Validate the Postgres URL
	if options.PostgresURL == "" {
		return Options{}, ErrBadPostgresURL
	}

	// Validate the Schema directory
	if options.SchemaDir == "" {
		return Options{}, ErrBadSchemaDir
	}
	if _, err := os.Stat(options.SchemaDir); errors.Is(err, fs.ErrNotExist) {
		return Options{}, err
	}

	return options, nil
}

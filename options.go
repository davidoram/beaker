package main

import (
	"errors"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
)

type Options struct {
	CredentialsFile string
	PostgresURL     string
}

var (
	ErrBadCredentialsFile = errors.New("Invalid credentials file path")
	ErrBadPostgresURL     = errors.New("Invalid Postgres URL")
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
	options := Options{
		CredentialsFile: filepath.Join(homeDir, "credentials.txt"),
		PostgresURL:     "postgres://postgres:password@localhost:5432/beaker_development?sslmode=disable",
	}

	// Use flags to parse command line arguments
	flagset.StringVar(&options.CredentialsFile, "credentials", options.CredentialsFile, "Path to the NATS credentials file. See https://docs.nats.io/nats-concepts/security/ for details")
	flagset.StringVar(&options.PostgresURL, "postgres", options.PostgresURL, "Postgres connection URL. See https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING for details")

	// Parse the provided arguments
	if err := flagset.Parse(args); err != nil {
		return Options{}, err
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

	return options, nil
}

package config

import (
	myfeeder "github.com/go-mods/tagsvar/modules/config/feeder"
	configLoader "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strings"

	_ "github.com/go-mods/tagsvar/modules/logger"
)

// C is the global config
// It is initialized in the init() function
var C *AppConfig

// Version is the version of the application
// It is set at build time
var Version = "development"

// BuildDate is the date of the build
// It is set at build time
var BuildDate = "not defined"

// AppConfig holds the configuration of the application
// It is loaded from environment variables or used default values.
type AppConfig struct {
	// Version is the version of the application
	Version string
	// BuildDate is the date of the build
	BuildDate string
	// Prefix is the prefix of the generated files
	Prefix string `env:"TAGSVAR_PREFIX" default:""`
	// Suffix is the suffix of the generated files
	Suffix string `env:"TAGSVAR_SUFFIx" default:".tags"`
	// Verbose enables verbose output
	Verbose bool `env:"TAGSVAR_VERBOSE" default:"false"`
	// Silent disables output
	Silent bool `env:"TAGSVAR_SILENT" default:"false"`
}

// init initializes the config
func init() {
	// Create an instance of AppConfig
	C = &AppConfig{}

	// Load configuration access
	C.load()
}

// load loads the configuration from default values and environment variables
func (c *AppConfig) load() {

	// Create the config loader
	loader := configLoader.New()

	// Add feeder to load from default values
	loader.AddFeeder(myfeeder.Default{})

	// Add feeder to load from environment variables
	loader.AddFeeder(feeder.Env{})

	// Read config access
	err := loader.AddStruct(c).Feed()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load environment variables")
	}
}

// Setup : this function is called while the config is loaded by golobby/config
func (c *AppConfig) Setup() {
	c.Version = Version
	c.BuildDate = BuildDate
}

// IsGoFile checks if the file is a .go file
func (c *AppConfig) IsGoFile(fileName string) bool {
	return filepath.Ext(fileName) == ".go"
}

// IsGeneratedFile checks if the file is a generated file
// It must be a .go file
// It must start with the prefix and end with the suffix
func (c *AppConfig) IsGeneratedFile(fileName string) bool {
	// Get the file name only
	fileName = filepath.Base(fileName)
	// Check if the file is a .go file
	if !c.IsGoFile(fileName) {
		return false
	}
	// Remove the extension
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	// Check if the file starts with the prefix
	if !strings.HasPrefix(filepath.Base(fileName), c.Prefix) {
		return false
	}
	// Check if the file ends with the suffix
	if !strings.HasSuffix(filepath.Base(fileName), c.Suffix) {
		return false
	}
	return true
}

// IsProjectFile checks if the file is a project file (not a generated file, neither a test file)
func (c *AppConfig) IsProjectFile(fileName string) bool {
	// Get the file name only
	fileName = filepath.Base(fileName)
	// Check if the file is a .go file
	if !c.IsGoFile(fileName) {
		return false
	}
	// Check if the file is a generated file
	if c.IsGeneratedFile(fileName) {
		return false
	}
	// Check if the file is a test file
	if strings.HasSuffix(fileName, "_test.go") {
		return false
	}
	return true
}

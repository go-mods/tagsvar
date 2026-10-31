package config

import (
	"runtime/debug"

	configLoader "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/rs/zerolog/log"

	myfeeder "github.com/go-mods/tagsvar/modules/config/feeder"
	_ "github.com/go-mods/tagsvar/modules/logger"
)

// C is the global config
// It is initialized in the init() function
var C *AppConfig

// Version is the version of the application
// It is set at build time (task go:build). Otherwise, the version of the module
// recorded by go install is used, then "development".
var Version = ""

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
	Suffix string `env:"TAGSVAR_SUFFIX" default:".vars"`
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
	c.Version = appVersion()
	c.BuildDate = BuildDate
}

// appVersion returns the version set at build time, or the version of the module
// recorded by go install (github.com/go-mods/tagsvar@v1.2.3 gives v1.2.3)
func appVersion() string {
	if Version != "" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "development"
}

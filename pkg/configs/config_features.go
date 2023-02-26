package configs

import "github.com/kelseyhightower/envconfig"

// Features configuration
type Features struct {
	Swagger bool `default:"true"`
}

// FeaturesConfig processes env
// to feature flags configuration
func FeaturesConfig() Features {
	var featureFlags Features
	envconfig.MustProcess("FEATURE", &featureFlags)

	return featureFlags
}

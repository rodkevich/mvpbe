package features

// Config configuration
type Config struct {
	Swagger bool `envconfig:"FEATURE_SWAGGER" default:"false"`
}

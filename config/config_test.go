package config

//
// func TestNewConfig(t *testing.T) {
// 	t.Parallel()
// 	err := godotenv.Load("../examples/example.env")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	cfg := &sample.Config{
// 		Http:     HTTPConfig(),
// 		Database: DataSourceConfig(),
// 		Cache:    CacheConfig(),
// 		Features: FeaturesConfig(),
// 	}
// 	assert.Equal(t, "mvp_db", cfg.Database.DBName)
// 	assert.Equal(t, false, cfg.Cache.Enable)
// }

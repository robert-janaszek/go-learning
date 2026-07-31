package config

type AppConfig struct {
	Port int
}

func Load() {
	parseEnv()
}

func parseEnv() {}

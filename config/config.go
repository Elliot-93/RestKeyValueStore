package config

type Config struct {
	Port     int
	LruDepth int
}

var config Config

func Populate(port, lruDepth int) {
	config = Config{Port: port, LruDepth: lruDepth}
}

func Read() Config {
	return config
}

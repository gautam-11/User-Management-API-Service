package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

type Constants struct {
	PORT  string
	Mongo struct {
		URL    string
		DBName string
	}
}

type Config struct {
	Constants
	Session  *mgo.Session
	Database *mgo.Database
}

func EnvPort() (*Config, error) {
	config := Config{}
	constants, err := initViper()
	if err != nil {
		return nil, err
	}
	config.Constants = constants
	return &config, err

}

// Connect() is used to establish a database connection to mongodb
func Connect() (*Config, error) {
	config := Config{}
	constants, err := initViper()
	config.Constants = constants
	if err != nil {
		return &config, err
	}
	dbSession, err := mgo.Dial(config.Constants.Mongo.URL)
	if err != nil {
		return &config, err
	}
	config.Session = dbSession
	config.Database = dbSession.DB(config.Constants.Mongo.DBName)

	fmt.Println("Connected to database!!")
	return &config, err
}

func initViper() (Constants, error) {
	viper.SetConfigName("user-api.config") // Configuration fileName without the .TOML or .YAML extension
	viper.AddConfigPath(".")               // Search the root directory for the configuration file
	err := viper.ReadInConfig()            // Find and read the config file
	if err != nil {                        // Handle errors reading the config file
		return Constants{}, err
	}
	viper.WatchConfig() // Watch for changes to the configuration file and recompile
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.SetDefault("PORT", "8080")
	if err = viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	var constants Constants
	err = viper.Unmarshal(&constants)
	return constants, err
}

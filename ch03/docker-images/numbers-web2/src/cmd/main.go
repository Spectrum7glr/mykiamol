package main

import (
	"numbersweb2"
)

// func initializeConfig() {
// 	viper.SetDefault("api_endpoint", "http://numbers-api/rnd")
// 	viper.SetDefault("port", "80")
// 	viper.SetDefault("timeout", 2)
// 	viper.AddConfigPath(".")
// 	viper.AddConfigPath("./config")
// 	viper.SetConfigName("ch04")
// 	viper.SetConfigType("env")
// 	dir, _ := os.Getwd()
// 	fmt.Println(dir)
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// 	viper.AutomaticEnv()
// }

func main() {
	// initializeConfig()
	numbersweb2.InitEnvConfigs()
	numbersweb2.RunServer()
}

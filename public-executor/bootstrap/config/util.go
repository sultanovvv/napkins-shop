package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

// convertStringToSlice преобразует строку вида 323232,888 в слайс целых чисел
func convertStringToSlice[T Int](input string) []T {
	if input == "" {
		return []T{}
	}
	strNumbers := strings.Split(input, ",")

	var result []T
	for _, strNum := range strNumbers {
		num, err := strconv.ParseInt(strNum, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("convert to number: %v", err))
		}
		result = append(result, T(num))
	}

	return result
}

func InitViperByEnv(configPath string) {
	if configPath == "" {
		zap.S().Info("Config path not provided, using environment variables")
		viper.AutomaticEnv()
		return
	}

	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			zap.S().Errorf("Config file does not exist: %s", configPath)
		} else {
			zap.S().Errorf("Error checking config file: %v", err)
		}
		viper.AutomaticEnv()
		return
	}

	viper.SetConfigFile(configPath)
	readInConfig()
}

func readInConfig() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			zap.S().Errorf("Config file not found, using defaults settings: %v", err)
		} else {
			zap.S().Fatalf("Error reading config file: %v", err)
		}
	}
}

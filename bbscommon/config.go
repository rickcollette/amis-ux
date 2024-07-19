package bbscommon

import (
	"encoding/json"
	"os"
)

type Config struct {
	BBSName           string `json:"bbs_name"`
	SysopName         string `json:"sysop_name"`
	AllowNewUsers     bool   `json:"allow_new_users"`
	AsciiFolder       string `json:"ascii_folder"`
	AtasciiFolder     string `json:"atascii_folder"`
	AnsiFolder        string `json:"ansi_folder"`
	MenusFolder       string `json:"menus_folder"`
	ExecutablesFolder string `json:"executables_folder"`
	SysopPassword     string `json:"sysop_password"`
	PortNumber        int    `json:"port_number"`
}

func LoadConfig(configPath string) (Config, error) {
	var config Config
	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}
	return config, nil
}

func SaveConfig(configPath string, config Config) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

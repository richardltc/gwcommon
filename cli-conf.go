package gwcommon

import (
	"log"

	"github.com/spf13/viper"
)

const (
	// CCLIConfFile - To be used only by GoDeploy
	CCLIConfFile    string = "cli"
	CCLIConfFileExt string = ".yaml"
)

// CLIConfStruct - The CLI application config struct
type CLIConfStruct struct {
	BinFolder                 string      // The folder that contains the coin binary files
	Currency                  string      // USD, GBP
	FirstTimeRun              bool        // Is this the first time the server has run? If so, we need to store the BinFolder
	ProjectType               ProjectType // The project type
	Port                      string      // The port that the server should run on
	RefreshTimer              int         // Refresh interval
	RPCuser                   string      // The rpcuser
	RPCpassword               string      // The rpc password
	ServerIP                  string      // The IP address of the coin daemon server
	Token                     string      // Stored after generation and is checked to be equal with the clients
	UserConfirmedSeedRecovery bool        // Whether or not the user has said they've stored their recovery seed has been stored
}

func GetCLIConfStruct() (CLIConfStruct, error) {

	viper.SetConfigName(CCLIConfFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	var cs CLIConfStruct

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cs)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return cs, nil
}

// SetCLIConfStruct - Save the CLI config struct via viper
func SetCLIConfStruct(cs CLIConfStruct) error {

	viper.SetConfigName(CCLIConfFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.Set("BinFolder", cs.BinFolder)
	viper.Set("Currency", cs.Currency)
	viper.Set("FirstTimeRun", cs.FirstTimeRun)
	viper.Set("ProjectType", cs.ProjectType)
	viper.Set("RefreshTimer", cs.RefreshTimer)
	viper.Set("rpcuser", cs.RPCuser)
	viper.Set("rpcpassword", cs.RPCpassword)
	viper.Set("ServerIP", cs.ServerIP)
	viper.Set("Port", cs.Port)
	viper.Set("Token", cs.Token)
	viper.Set("UserConfirmedSeedRecovery", cs.UserConfirmedSeedRecovery)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

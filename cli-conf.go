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
	FirstTimeRun              bool        // Is this the first time the server has run? If so, we need to store the BinFolder
	ProjectType               ProjectType // The project type
	Port                      string      // The port that the server should run on
	RPCuser                   string      // The rpcuser
	RPCpassword               string      // The rpc password
	ServerIP                  string      // The IP address of the coin daemon server
	Token                     string      // Stored after generation and is checked to be equal with the clients
	UserConfirmedSeedRecovery bool        // Whether or not the user has said they've stored their recovery seed has been stored
}

// CreateDefaultCLIConfFile - Only to be used by GoDeploy
// func CreateDefaultCLIConfFile(confDir string, pt ProjectType) error {
// 	conf, err := newCLIConfStruct(pt)
// 	if err != nil {
// 		return err
// 	}

// 	jssb, err := json.MarshalIndent(conf, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	f, err := os.Create(confDir + CCLIConfFile)
// 	if err != nil {
// 		return err
// 	}

// 	log.Println("Creating default CLI config file " + f.Name())
// 	_, err = f.WriteString(string(jssb))
// 	err = f.Close()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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

// // GetCLIConfigStruct - Retrieve the application config struct
// func GetCLIConfigStruct(refreshFields bool) (CLIConfStruct, error) {

// 	// We can't do the below, because we don't know what project we currently are, as that's dictated by GoDeploy

// 	// Create the file if it doesn't already exist
// 	// dir := AddTrailingSlash(confDir)
// 	// if _, err := os.Stat(dir + cConfFile); os.IsNotExist(err) {
// 	// 	createDefaultConfFile(confDir, pt)
// 	// }

// 	// Get the config file
// 	dir, err := GetRunningDir()
// 	if err != nil {
// 		return CLIConfStruct{}, fmt.Errorf("Unable to GetRunningDir - %v", err)
// 	}
// 	file, err := ioutil.ReadFile(dir + CCLIConfFile)
// 	if err != nil {
// 		return CLIConfStruct{}, err
// 	}

// 	cs := CLIConfStruct{}

// 	err = json.Unmarshal([]byte(file), &cs)
// 	if err != nil {
// 		return CLIConfStruct{}, err
// 	}

// 	// Now, let's write the file back because it may have some new fields
// 	if refreshFields {
// 		SetCLIConfigStruct(dir, cs)
// 	}

// 	return cs, nil
// }

// func newCLIConfStruct(pt ProjectType) (CLIConfStruct, error) {
// 	cnf := CLIConfStruct{}
// 	var err error

// 	switch pt {
// 	case PTDivi:
// 		cnf.ProjectType = PTDivi
// 	case PTPhore:
// 		cnf.ProjectType = PTPhore
// 	case PTPIVX:
// 		cnf.ProjectType = PTPIVX
// 	case PTTrezarcoin:
// 		cnf.ProjectType = PTTrezarcoin
// 	default:
// 		err = errors.New("Unable to determine ProjectType")
// 	}

// 	cnf.Port = "4000"
// 	cnf.ServerIP = "127.0.0.1"

// 	if err != nil {
// 		return cnf, err
// 	}

// 	return cnf, nil
// }

// SetCLIConfStruct - Save the CLI config struct via viper
func SetCLIConfStruct(cs CLIConfStruct) error {

	viper.SetConfigName(CCLIConfFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.Set("BinFolder", cs.BinFolder)
	viper.Set("FirstTimeRun", cs.FirstTimeRun)
	viper.Set("ProjectType", cs.ProjectType)
	viper.Set("rpcuser", cs.RPCuser)
	viper.Set("rpcpassword", cs.RPCpassword)
	viper.Set("ServerIP", cs.ServerIP)
	viper.Set("Port", cs.Port)
	viper.Set("Token", cs.Token)
	viper.Set("UserConfirmedSeedRecovery", cs.UserConfirmedSeedRecovery)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// // SetCLIConfigStruct - Save the application config struct
// func SetCLIConfigStruct(dir string, cs CLIConfStruct) error {
// 	jssb, _ := json.MarshalIndent(cs, "", "  ")
// 	dir = AddTrailingSlash(dir)
// 	sFile := dir + CCLIConfFile

// 	f, err := os.Create(sFile)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = f.WriteString(string(jssb))
// 	err = f.Close()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

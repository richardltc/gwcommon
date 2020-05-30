package gwcommon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	// CCLIConfFile - To be used only by GoDeploy
	CCLIConfFile string = "cli-config.json"
)

// CLIConfStruct - The CLI application config struct
type CLIConfStruct struct {
	AppName                   string
	ProjectType               ProjectType
	ServerIP                  string
	Port                      string
	UserConfirmedSeedRecovery bool
}

// CreateDefaultCLIConfFile - Only to be used by GoDeploy
func CreateDefaultCLIConfFile(confDir string, pt ProjectType) error {
	conf, err := newCLIConfStruct(pt)
	if err != nil {
		return err
	}

	jssb, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.Create(confDir + CServerConfFile)
	if err != nil {
		return err
	}

	log.Println("Creating default CLI config file " + f.Name())
	_, err = f.WriteString(string(jssb))
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

// GetCLIConfigStruct - Retrieve the application config struct
func GetCLIConfigStruct(refreshFields bool) (CLIConfStruct, error) {

	// We can't do the below, because we don't know what project we currently are, as that's dictated by GoDeploy

	// Create the file if it doesn't already exist
	// dir := AddTrailingSlash(confDir)
	// if _, err := os.Stat(dir + cConfFile); os.IsNotExist(err) {
	// 	createDefaultConfFile(confDir, pt)
	// }

	// Get the config file
	dir, err := GetRunningDir()
	if err != nil {
		return CLIConfStruct{}, fmt.Errorf("Unable to GetRunningDir - %v", err)
	}
	file, err := ioutil.ReadFile(dir + CCLIConfFile)
	if err != nil {
		return CLIConfStruct{}, err
	}

	cs := CLIConfStruct{}

	err = json.Unmarshal([]byte(file), &cs)
	if err != nil {
		return CLIConfStruct{}, err
	}

	// Now, let's write the file back because it may have some new fields
	if refreshFields {
		SetCLIConfigStruct(dir, cs)
	}

	return cs, nil
}

func newCLIConfStruct(pt ProjectType) (CLIConfStruct, error) {
	cnf := CLIConfStruct{}
	var err error

	switch pt {
	case PTDivi:
		cnf.AppName = CAppNameGoDivi
		cnf.ProjectType = PTDivi
	case PTPhore:
		cnf.AppName = CAppNameGoPhore
		cnf.ProjectType = PTPhore
	case PTPIVX:
		cnf.AppName = CAppNameGoPIVX
		cnf.ProjectType = PTPIVX
	case PTTrezarcoin:
		cnf.AppName = CAppNameGoTrezarcoin
		cnf.ProjectType = PTTrezarcoin
	default:
		err = errors.New("Unable to determine ProjectType")
	}

	cnf.UserConfirmedSeedRecovery = false
	cnf.Port = "4000"
	cnf.ServerIP = "127.0.0.1"

	if err != nil {
		return cnf, err
	}

	return cnf, nil
}

// SetCLIConfigStruct - Save the application config struct
func SetCLIConfigStruct(dir string, cs CLIConfStruct) error {
	jssb, _ := json.MarshalIndent(cs, "", "  ")
	dir = AddTrailingSlash(dir)
	sFile := dir + CServerConfFile

	f, err := os.Create(sFile)
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(jssb))
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

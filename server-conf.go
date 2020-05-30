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
	// CServerConfFile - To be used only by GoDeploy
	CServerConfFile string = "server-config.json"
)

// ServerConfStruct - The server application config struct
type ServerConfStruct struct {
	AppName                   string
	ProjectType               ProjectType
	ServerIP                  string
	Port                      string
	UserConfirmedSeedRecovery bool
}

// CreateDefaultServerConfFile - Only to be used by GoDeploy
func CreateDefaultServerConfFile(confDir string, pt ProjectType) error {
	conf, err := newServerConfStruct(pt)
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

	log.Println("Creating default server config file " + f.Name())
	_, err = f.WriteString(string(jssb))
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

// GetServerConfigStruct - Retrieve the application config struct
func GetServerConfigStruct(refreshFields bool) (ServerConfStruct, error) {

	// We can't do the below, because we don't know what project we currently are, as that's dictated by GoDeploy

	// Create the file if it doesn't already exist
	// dir := AddTrailingSlash(confDir)
	// if _, err := os.Stat(dir + cConfFile); os.IsNotExist(err) {
	// 	createDefaultConfFile(confDir, pt)
	// }

	// Get the config file
	dir, err := GetRunningDir()
	if err != nil {
		return ServerConfStruct{}, fmt.Errorf("Unable to GetRunningDir - %v", err)
	}
	file, err := ioutil.ReadFile(dir + CServerConfFile)
	if err != nil {
		return ServerConfStruct{}, err
	}

	cs := ServerConfStruct{}

	err = json.Unmarshal([]byte(file), &cs)
	if err != nil {
		return ServerConfStruct{}, err
	}

	// Now, let's write the file back because it may have some new fields
	if refreshFields {
		SetServerConfigStruct(dir, cs)
	}

	return cs, nil
}

func newServerConfStruct(pt ProjectType) (ServerConfStruct, error) {
	cnf := ServerConfStruct{}
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

	cnf.Port = "4000"
	cnf.ServerIP = "127.0.0.1"

	cnf.UserConfirmedSeedRecovery = false

	if err != nil {
		return cnf, err
	}

	return cnf, nil
}

// SetServerConfigStruct - Save the application config struct
func SetServerConfigStruct(dir string, cs ServerConfStruct) error {
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

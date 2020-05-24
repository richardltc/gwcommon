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
	// CConfFile - To be used only by GoDeploy
	CConfFile string = "config.json"
)

// ConfStruct - The global application config struct
type ConfStruct struct {
	AppName                   string
	ProjectType               ProjectType
	UserConfirmedSeedRecovery bool
}

// CreateDefaultConfFile - Only to be used by GoDeploy
func CreateDefaultConfFile(confDir string, pt ProjectType) error {
	conf, err := newConfStruct(pt)
	if err != nil {
		return err
	}

	jssb, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.Create(confDir + CConfFile)
	if err != nil {
		return err
	}

	log.Println("Creating default config file " + f.Name())
	_, err = f.WriteString(string(jssb))
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

// GetConfigStruct - Retrieve the application config struct
func GetConfigStruct(refreshFields bool) (ConfStruct, error) {

	// We can't do the below, because we don't know what project we currently are, as that's dictated by GoDeploy

	// Create the file if it doesn't already exist
	// dir := AddTrailingSlash(confDir)
	// if _, err := os.Stat(dir + cConfFile); os.IsNotExist(err) {
	// 	createDefaultConfFile(confDir, pt)
	// }

	// Get the config file
	dir, err := GetRunningDir()
	if err != nil {
		return ConfStruct{}, fmt.Errorf("Unable to GetRunningDir - %v", err)
	}
	file, err := ioutil.ReadFile(dir + CConfFile)
	if err != nil {
		return ConfStruct{}, err
	}

	cs := ConfStruct{}

	err = json.Unmarshal([]byte(file), &cs)
	if err != nil {
		return ConfStruct{}, err
	}

	// Now, let's write the file back because it may have some new fields
	if refreshFields {
		SetConfigStruct(dir, cs)
	}

	return cs, nil
}

func newConfStruct(pt ProjectType) (ConfStruct, error) {
	cnf := ConfStruct{}
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

	if err != nil {
		return cnf, err
	}

	return cnf, nil
}

// SetConfigStruct - Save the application config struct
func SetConfigStruct(dir string, cs ConfStruct) error {
	jssb, _ := json.MarshalIndent(cs, "", "  ")
	dir = AddTrailingSlash(dir)
	sFile := dir + CConfFile

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

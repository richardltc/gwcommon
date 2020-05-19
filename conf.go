package gwcommon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	cConfFile string = "config.json"
)

type ConfStruct struct {
	AppName     string
	ProjectType ProjectType
}

func CreateDefaultConfFile(confDir string, pt ProjectType) error {
	var conf = newConfStruct(pt)

	jssb, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.Create(confDir + cConfFile)
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

func GetConfigStruct(confDir string, refreshFields bool) (ConfStruct, error) {

	// We can't do the below, because we don't know what project we currently are, as that's dictated by GoDeploy
	// Create the file if it doesn't already exist
	// dir := AddTrailingSlash(confDir)
	// if _, err := os.Stat(dir + cConfFile); os.IsNotExist(err) {
	// 	createDefaultConfFile(confDir, pt)
	// }

	// Get the config file
	dir := AddTrailingSlash(confDir)
	file, err := ioutil.ReadFile(dir + cConfFile)
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

func newConfStruct(pt ProjectType) ConfStruct {
	cnf := ConfStruct{}

	switch pt {
	case PTDivi:
		cnf.AppName = CAppNameGoDivi
	case PTPIVX:
		//TODO Populate with PIVX and Trezar
	}

	return cnf
}

func SetConfigStruct(dir string, cs ConfStruct) error {
	jssb, _ := json.MarshalIndent(cs, "", "  ")
	dir = AddTrailingSlash(dir)
	sFile := dir + cConfFile

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

package gwcommon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	CConfFile string = "config.json"
)

type ConfStruct struct {
	AppName string
}

func CreateDefaultConfFile(confDir string, pt ProjectType) error {
	var conf = newConfStruct()

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

func getConfigStruct(confDir string, refreshFields bool, pt ProjectType) (ConfStruct, error) {

	// Create the file if it doesn't already exist
	dir := AddTrailingSlash(confDir)
	if _, err := os.Stat(dir + CConfFile); os.IsNotExist(err) {
		CreateDefaultConfFile(confDir, pt)
	}

	// Get the config file
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
		setConfigStruct(dir, cs)
	}

	return cs, nil
}

func newConfStruct(pt ProjectType) ConfStruct {
	cnf := ConfStruct{}
	cnf.AppName = "Enter SendGrid Key here"

	return cnf
}

func setConfigStruct(dir string, cs ConfStruct) error {
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

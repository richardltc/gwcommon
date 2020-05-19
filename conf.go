package gdcommon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	gdc "richardmace.co.uk/godivi/gdcommon"
)

const (
	cConfFile string = "config.json"
)

type confStruct struct {
	AppName string
}

func CreateDefaultConfFile(confDir string, pt gdc.ProjectType) error {
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

func getConfigStruct(confDir string, refreshFields bool) (confStruct, error) {

	// Create the file if it doesn't already exist
	dir := AddTrailingSlash(confDir)
	if _, err := os.Stat(dir + cConfFile); os.IsNotExist(err) {
		createDefaultConfFile(confDir)
	}

	// Get the config file
	file, err := ioutil.ReadFile(dir + cConfFile)
	if err != nil {
		return confStruct{}, err
	}

	cs := confStruct{}

	err = json.Unmarshal([]byte(file), &cs)
	if err != nil {
		return confStruct{}, err
	}

	// Now, let's write the file back because it may have some new fields
	if refreshFields {
		setConfigStruct(dir, cs)
	}

	return cs, nil
}

func newConfStruct(pt gdc.ProjectType) confStruct {
	cnf := confStruct{}
	cnf.AppName = "Enter SendGrid Key here"

	return cnf
}

func setConfigStruct(dir string, cs confStruct) error {
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

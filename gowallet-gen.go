package gwcommon

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/inconshreveable/go-update"
)

const (
	// CAppVersion - The app version of the suite of apps
	CAppVersion string = "0.21.3" // All of the individual apps will have the same version to make it easier for the user
	cUnknown    string = "Unknown"
	// CDownloadURLGD - The download file lotcation for GoDivi
	CDownloadURLGD string = "https://bitbucket.org/rmace/godivi/downloads/"

	// CAppCLIFileCompiled - Should only be used by GoDeploy
	CAppCLIFileCompiled string = "cli"
	// CAppCLIFileCompiledWin - Should only be used by GoDeploy
	CAppCLIFileCompiledWin string = "cli.exe"
	// CAppServerFileCompiled - Should only be used by GoDeploy
	CAppServerFileCompiled string = "web"
	// CAppServerFileCompiledWin - Should only be used by GoDeploy
	CAppServerFileCompiledWin string = "web.exe"
	// CAppUpdaterFileCompiled - Should only be used by GoDeploy
	CAppUpdaterFileCompiled string = "updater"
	// CAppUpdaterFileCompiledWin - Should only be used by GoDeploy
	CAppUpdaterFileCompiledWin string = "updater.exe"

	cWalletSeedFileGoDivi string = "unsecure-divi-seed.txt"

	// Divid Responses
	cDiviDNotRunningError     string = "error: couldn't connect to server"
	cDiviDDIVIServerStarting  string = "DIVI server starting"
	cDividRespWalletEncrypted string = "wallet encrypted"

	cGoDiviExportPath         string = "export PATH=$PATH:"
	CUninstallConfirmationStr string = "Confirm"
	CSeedStoredSafelyStr      string = "Confirm"

	// CMinRequiredMemoryMB - Needed by install command
	CMinRequiredMemoryMB int = 920
	CMinRequiredSwapMB   int = 2048

	// Wallet Security Statuses - Should be types?
	CWalletStatusLocked      string = "locked"
	CWalletStatusUnlocked    string = "unlocked"
	CWalletStatusLockedAndSk string = "locked-anonymization"
	CWalletStatusUnEncrypted string = "unencrypted"

	// Progress constants
	cProgress1  string = ">     "
	cProgress2  string = "=>    "
	cProgress3  string = "==>   "
	cProgress4  string = "===>  "
	cProgress5  string = "====> "
	cProgress6  string = "=====>"
	cProgress7  string = " ====="
	cProgress8  string = "  ===="
	cProgress9  string = "   ==="
	cProgress10 string = "    =="
	cProgress11 string = "     ="

	cUtfTick     string = "\u2713"
	CUtfTickBold string = "\u2714"

	cCircProg1 string = "\u25F7"
	cCircProg2 string = "\u25F6"
	cCircProg3 string = "\u25F5"
	cCircProg4 string = "\u25F4"

	cUtfLock string = "\u1F512"
)

// APPType - either APPTCLI, APPTCLICompiled, APPTInstaller, APPTUpdater, APPTServer
type APPType int

const (
	// APPTCLI - e.g. godivi
	APPTCLI APPType = iota
	// APPTCLICompiled - e.g. cli
	APPTCLICompiled
	// APPTInstaller e.g. godivi-installer
	APPTInstaller
	// APPTUpdater e.g. update-godivi
	APPTUpdater
	// APPTUpdaterCompiled e.g. updater
	APPTUpdaterCompiled
	// APPTServer e.g. godivis
	APPTServer
	// APPTServerCompiled e.g. web
	APPTServerCompiled
)

// OSType - either ostArm, ostLinux or ostWindows
type OSType int

const (
	// OSTArm - Arm
	OSTArm OSType = iota
	// OSTLinux - Linux
	OSTLinux
	// OSTWindows - Windows
	OSTWindows
)

// ProjectType - To allow external to determine what kind of wallet we are working with
type ProjectType int

const (
	// PTDivi - Divi
	PTDivi ProjectType = iota
	// PTPhore - Phore
	PTPhore
	// PTPIVX - PIVX
	PTPIVX
	// PTTrezarcoin - TrezarCoin
	PTTrezarcoin
)

type progressBarType int

const (
	pbType progressBarType = iota
	pbtBCSBar
	pbtMNSBar
)

type application struct {
	infoLog *log.Logger
}

type args struct {
	startFresh *bool
	dbug       *bool
}

// ServerResponse - Determine REST Server response
type ServerResponse int

const (
	NotRequired      ServerResponse = 0
	MalformedRequest ServerResponse = 1
	NoServerError    ServerResponse = 2
)

// Server Request COnstants
const (
	CServRequestGenerateToken  string = "GenerateToken"
	CServRequestShutdownServer string = "ShutdownServer"
)

type listTransactions []struct {
	Account         string        `json:"account"`
	Address         string        `json:"address"`
	Category        string        `json:"category"`
	Amount          float64       `json:"amount"`
	Vout            int           `json:"vout"`
	Confirmations   int           `json:"confirmations"`
	Bcconfirmations int           `json:"bcconfirmations"`
	Blockhash       string        `json:"blockhash"`
	Blockindex      int           `json:"blockindex"`
	Blocktime       int           `json:"blocktime"`
	Txid            string        `json:"txid"`
	Walletconflicts []interface{} `json:"walletconflicts"`
	Time            int           `json:"time"`
	Timereceived    int           `json:"timereceived"`
}

type stakingStatusStruct struct {
	Validtime       bool `json:"validtime"`
	Haveconnections bool `json:"haveconnections"`
	Walletunlocked  bool `json:"walletunlocked"`
	Mintablecoins   bool `json:"mintablecoins"`
	Enoughcoins     bool `json:"enoughcoins"`
	Mnsync          bool `json:"mnsync"`
	StakingStatus   bool `json:"staking status"`
}

type tickerStruct struct {
	DIVI struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Symbol            string      `json:"symbol"`
		Slug              string      `json:"slug"`
		NumMarketPairs    int         `json:"num_market_pairs"`
		DateAdded         time.Time   `json:"date_added"`
		Tags              []string    `json:"tags"`
		MaxSupply         interface{} `json:"max_supply"`
		CirculatingSupply float64     `json:"circulating_supply"`
		TotalSupply       float64     `json:"total_supply"`
		Platform          interface{} `json:"platform"`
		CmcRank           int         `json:"cmc_rank"`
		LastUpdated       time.Time   `json:"last_updated"`
		Quote             struct {
			BTC struct {
				Price            float64   `json:"price"`
				Volume24H        float64   `json:"volume_24h"`
				PercentChange1H  float64   `json:"percent_change_1h"`
				PercentChange24H float64   `json:"percent_change_24h"`
				PercentChange7D  float64   `json:"percent_change_7d"`
				MarketCap        float64   `json:"market_cap"`
				LastUpdated      time.Time `json:"last_updated"`
			} `json:"BTC"`
			USD struct {
				Price            float64   `json:"price"`
				Volume24H        float64   `json:"volume_24h"`
				PercentChange1H  float64   `json:"percent_change_1h"`
				PercentChange24H float64   `json:"percent_change_24h"`
				PercentChange7D  float64   `json:"percent_change_7d"`
				MarketCap        float64   `json:"market_cap"`
				LastUpdated      time.Time `json:"last_updated"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"DIVI"`
}

var lastBCSyncStatus string = ""
var lastMNSyncStatus string = ""

// AddProjectPath - Add the coin project path to the profile
func AddProjectPath() error {
	if runtime.GOOS != "windows" {
		u, err := user.Current()
		if err != nil {
			return err
		}
		sProfile := AddTrailingSlash(u.HomeDir) + ".profile"
		gdf, err := GetAppsBinFolder()
		if err != nil {
			return fmt.Errorf("Unable to GetAppsBinFolder: %v ", err)
		}

		if FileExists(sProfile) {
			// First make sure that the path hasn't already been added
			seif, err := stringExistsInFile(cGoDiviExportPath+gdf, sProfile)
			if err != nil {
				return fmt.Errorf("err checking profile: %v ", err)
			}
			if !seif {
				f, err := os.OpenFile(sProfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
				if err != nil {
					return fmt.Errorf("Unable to open profile: %v ", err)
				}

				defer f.Close()
				fmt.Fprintf(f, "%s\n", cGoDiviExportPath+gdf)

				//run "source $HOME/.profile" so that the godivi path is live
				cmd := exec.Command("bash", "-c", "source "+sProfile)
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("error: %v", err)
				}
			}
		}
	}
	return nil
}

// ConvertBCVerification - Convert Blockchain verification progress
func ConvertBCVerification(verificationPG float64) string {
	var sProg string
	var fProg float64

	fProg = verificationPG * 100
	sProg = fmt.Sprintf("%.2f", fProg)

	return sProg
}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
	}
	return err
}

// GetAppsBinFolder - Returns the directory of where the apps binary files are stored
func GetAppsBinFolder() (string, error) {
	var s string
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	//hd := getUserHomeDir()
	hd := u.HomeDir
	if runtime.GOOS == "windows" {
		// add the "appdata\roaming" part.
		switch gwconf.ProjectType {
		case PTDivi:
			s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(cDiviBinDirWin)
		case PTPIVX:
			s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(cPIVXBinDirWin)
		case PTTrezarcoin:
			s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(cTrezarcoinBinDirWin)
		default:
			err = errors.New("Unable to determine ProjectType")
		}

	} else {
		switch gwconf.ProjectType {
		case PTDivi:
			s = AddTrailingSlash(hd) + AddTrailingSlash(cDiviBinDir)
		case PTPIVX:
			s = AddTrailingSlash(hd) + AddTrailingSlash(cPIVXBinDir)
		case PTTrezarcoin:
			s = AddTrailingSlash(hd) + AddTrailingSlash(cTrezarcoinBinDir)
		default:
			err = errors.New("Unable to determine ProjectType")
		}
	}
	return s, nil
}

// GetAppFileName - Returns the name of the app binary file e.g. godivi, godivis, godivi-installer
func GetAppFileName(an APPType) (string, error) {
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		switch an {
		case APPTCLI:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileGoDivi, nil
			case "linux":
				return CAppCLIFileGoDivi, nil
			case "windows":
				return CAppCLIFileWinGoDivi, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}

		case APPTCLICompiled:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileCompiled, nil
			case "linux":
				return CAppCLIFileCompiled, nil
			case "windows":
				return CAppCLIFileCompiledWin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTInstaller:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileInstallerGoDivi, nil
			case "linux":
				return CAppCLIFileInstallerGoDivi, nil
			case "windows":
				return CAppCLIFileInstallerWinGoDivi, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTServer:
			switch runtime.GOOS {
			case "arm":
				return CAppServerFileGoDivi, nil
			case "linux":
				return CAppServerFileGoDivi, nil
			case "windows":
				return CAppServerFileWinGoDivi, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTServerCompiled:
			switch runtime.GOOS {
			case "arm":
				return CAppServerFileCompiled, nil
			case "linux":
				return CAppServerFileCompiled, nil
			case "windows":
				return CAppServerFileCompiledWin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTUpdater:
			if runtime.GOOS == "windows" {
				return CAppUpdaterFileWinGoDivi, nil
			} else {
				return CAppUpdaterFileGoDivi, nil
			}
		default:
			err = errors.New("Unable to determine ProjectType")
		}
	case PTPIVX:
		switch an {
		case APPTCLI:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileGoPIVX, nil
			case "linux":
				return CAppCLIFileGoPIVX, nil
			case "windows":
				return CAppCLIFileWinGoPIVX, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}

		case APPTCLICompiled:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileCompiled, nil
			case "linux":
				return CAppCLIFileCompiled, nil
			case "windows":
				return CAppCLIFileCompiledWin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTInstaller:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileInstallerGoPIVX, nil
			case "linux":
				return CAppCLIFileInstallerGoPIVX, nil
			case "windows":
				return CAppCLIFileInstallerWinGoPIVX, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTServer:
			switch runtime.GOOS {
			case "arm":
				return CAppServerFileGoPIVX, nil
			case "linux":
				return CAppServerFileGoPIVX, nil
			case "windows":
				return CAppServerFileWinGoPIVX, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTServerCompiled:
			switch runtime.GOOS {
			case "arm":
				return CAppServerFileCompiled, nil
			case "linux":
				return CAppServerFileCompiled, nil
			case "windows":
				return CAppServerFileCompiledWin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTUpdater:
			if runtime.GOOS == "windows" {
				return CAppUpdaterFileWinGoPIVX, nil
			} else {
				return CAppUpdaterFileGoPIVX, nil
			}
		default:
			err = errors.New("Unable to determine ProjectType")
		}
	case PTTrezarcoin:
		switch an {
		case APPTCLI:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileGoTrezarcoin, nil
			case "linux":
				return CAppCLIFileGoTrezarcoin, nil
			case "windows":
				return CAppCLIFileWinGoTrezarcoin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTCLICompiled:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileCompiled, nil
			case "linux":
				return CAppCLIFileCompiled, nil
			case "windows":
				return CAppCLIFileCompiledWin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTInstaller:
			switch runtime.GOOS {
			case "arm":
				return CAppCLIFileInstallerGoTrezarcoin, nil
			case "linux":
				return CAppCLIFileInstallerGoTrezarcoin, nil
			case "windows":
				return CAppCLIFileInstallerWinGoTrezarcoin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTServer:
			switch runtime.GOOS {
			case "arm":
				return CAppServerFileGoTrezarcoin, nil
			case "linux":
				return CAppServerFileGoTrezarcoin, nil
			case "windows":
				return CAppServerFileWinGoTrezarcoin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTServerCompiled:
			switch runtime.GOOS {
			case "arm":
				return CAppServerFileCompiled, nil
			case "linux":
				return CAppServerFileCompiled, nil
			case "windows":
				return CAppServerFileCompiledWin, nil
			default:
				err = errors.New("Unable to determine runtime.GOOS")
			}
		case APPTUpdater:
			if runtime.GOOS == "windows" {
				return CAppUpdaterFileWinGoPIVX, nil
			} else {
				return CAppUpdaterFileGoPIVX, nil
			}
		default:
			err = errors.New("Unable to determine ProjectType")
		}
	default:
		err = errors.New("Unable to determine ProjectType")

	}
	return "", nil
}

// GetAppCLIName - Returns the application CLI name e.g. GoDivi CLI
func GetAppCLIName() (string, error) {
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppNameCLIGoDivi, nil
	case PTPIVX:
		return CAppNameCLIGoPIVX, nil
	case PTTrezarcoin:
		return CAppNameCLIGoTrezarcoin, nil
	default:
		err = errors.New("Unable to determine ProjectType")
	}
	return "", nil
}

// GetAppLogfileName - Returns the application logfile name e.g. godivi.log
func GetAppLogfileName() (string, error) {
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppCLILogfileGoDivi, nil
	case PTPIVX:
		return CAppCLILogfileGoPIVX, nil
	case PTTrezarcoin:
		return CAppCLILogfileGoTrezarcoin, nil
	default:
		err = errors.New("Unable to determine ProjectType")

	}
	return "", nil
}

// GetAppServerName - Returns the application Server name e.g. GoDivi Server
func GetAppServerName() (string, error) {
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppNameServerGoDivi, nil
	case PTPIVX:
		return CAppNameServerGoPIVX, nil
	case PTTrezarcoin:
		return CAppNameServerGoTrezarcoin, nil
	default:
		err = errors.New("Unable to determine ProjectType")

	}
	return "", nil
}

// GetAppName - Returns the application name e.g. GoDivi
func GetAppName() (string, error) {
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppNameGoDivi, nil
	case PTPIVX:
		return CAppNameGoPIVX, nil
	case PTTrezarcoin:
		return CAppNameGoTrezarcoin, nil
	default:
		err = errors.New("Unable to determine ProjectType")

	}
	return "", nil
}

// GetCoinDaemonFilename - Return the coin daemon file name e.g. divid
func GetCoinDaemonFilename() (string, error) {
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CDiviDFile, nil
	case PTPIVX:
		return CDiviDFile, nil
	case PTTrezarcoin:
		return CTrezarcoinDFile, nil
	default:
		err = errors.New("Unable to determine ProjectType")

	}
	return "", nil
}

// GetCoinHomeFolder - Returns the ome folder for the coin e.g. .divi
func GetCoinHomeFolder() (string, error) {
	var s string
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	hd := u.HomeDir
	if runtime.GOOS == "windows" {
		// add the "appdata\roaming" part.
		switch gwconf.ProjectType {
		case PTDivi:
			s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(cDiviHomeDirWin)
		case PTPIVX:
			s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(cPIVXHomeDirWin)
		case PTTrezarcoin:
			s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(cTrezarcoinHomeDirWin)
		default:
			err = errors.New("Unable to determine ProjectType")

		}
	} else {
		switch gwconf.ProjectType {
		case PTDivi:
			s = AddTrailingSlash(hd) + AddTrailingSlash(cDiviHomeDir)
		case PTPIVX:
			s = AddTrailingSlash(hd) + AddTrailingSlash(cPIVXHomeDir)
		case PTTrezarcoin:
			s = AddTrailingSlash(hd) + AddTrailingSlash(cTrezarcoinHomeDir)
		default:
			err = errors.New("Unable to determine ProjectType")

		}
	}
	return s, nil
}

// GetCoinName - Returns the name of the coin e.g. Divi
func GetCoinName() (string, error) {
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return cCoinNameDivi, nil
	case PTPIVX:
		return cCoinNamePIVX, nil
	case PTTrezarcoin:
		return cCoinNameTrezarcoin, nil
	default:
		err = errors.New("Unable to determine ProjectType")
	}
	return "", nil
}

func GetNextProgMNIndicator(LIndicator string) string {
	if LIndicator == cProgress1 {
		lastMNSyncStatus = cProgress2
		return cProgress2
	} else if LIndicator == cProgress2 {
		lastMNSyncStatus = cProgress3
		return cProgress3
	} else if LIndicator == cProgress3 {
		lastMNSyncStatus = cProgress4
		return cProgress4
	} else if LIndicator == cProgress4 {
		lastMNSyncStatus = cProgress5
		return cProgress5
	} else if LIndicator == cProgress5 {
		lastMNSyncStatus = cProgress6
		return cProgress6
	} else if LIndicator == cProgress6 {
		lastMNSyncStatus = cProgress7
		return cProgress7
	} else if LIndicator == cProgress7 {
		lastMNSyncStatus = cProgress8
		return cProgress8
	} else if LIndicator == cProgress8 {
		lastMNSyncStatus = cProgress9
		return cProgress9
	} else if LIndicator == cProgress9 {
		lastMNSyncStatus = cProgress10
		return cProgress10
	} else if LIndicator == cProgress10 {
		lastMNSyncStatus = cProgress11
		return cProgress11
	} else if LIndicator == cProgress11 || LIndicator == "" {
		lastMNSyncStatus = cProgress1
		return cProgress1
	} else {
		lastMNSyncStatus = cProgress1
		return cProgress1
	}
}

func GetEncryptWalletResp() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(`Your wallet is currently UNENCRYPTED!

It is *highly* recommended that you encrypt your wallet before proceeding any further.

Encrypt it now?: (y/n)`)
	resp, _ := reader.ReadString('\n')
	resp = strings.ReplaceAll(resp, "\n", "")
	return resp
}

func GetWalletEncryptionPassword() string {
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i <= 2; i++ {
		fmt.Print("\nPlease enter a password to encrypt your wallet: ")
		epw1, _ := reader.ReadString('\n')
		fmt.Print("\nNow please re-enter your password: ")
		epw2, _ := reader.ReadString('\n')
		if epw1 != epw2 {
			fmt.Print("\nThe passwords don't match, please try again...\n")
		} else {
			s := strings.ReplaceAll(epw1, "\n", "")

			return s
		}
	}
	return ""
}

// GetWalletUnlockPassword - Retrieves the wallet unlock password that the user has entered
func GetWalletUnlockPassword() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPlease enter your wallet encryption password: ")
	pw, _ := reader.ReadString('\n')
	s := strings.ReplaceAll(pw, "\n", "")

	return s
}

func getYesNoResp(msg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg + " (y/n)")
	resp, _ := reader.ReadString('\n')
	resp = strings.ReplaceAll(resp, "\n", "")
	return resp
}

// IsGoWalletInstalled - Returns bool if GoWallet has been installed
func IsGoWalletInstalled() bool {
	// First, let's make sure that we have our divi bin folder
	dbf, _ := GetAppsBinFolder()

	if _, err := os.Stat(dbf); !os.IsNotExist(err) {
		// e.g. /home/user/godivi/ bin folder exists..
		return true
	}
	return false
}

// IsAppCLIRunning - Will then work out what wallet this relates to, and return bool whether the CLI app is running
func IsAppCLIRunning() (bool, int, error) {
	var pid int
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return false, pid, err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			pid, _, err = findProcess(CAppCLIFileWinGoDivi)
		} else {
			pid, _, err = findProcess(CAppCLIFileGoDivi)
		}
	case PTPIVX:
		if runtime.GOOS == "windows" {
			pid, _, err = findProcess(CAppCLIFileWinGoPIVX)
		} else {
			pid, _, err = findProcess(CAppCLIFileGoPIVX)
		}
	case PTTrezarcoin:
		if runtime.GOOS == "windows" {
			pid, _, err = findProcess(CAppCLIFileWinGoTrezarcoin)
		} else {
			pid, _, err = findProcess(CAppCLIFileGoTrezarcoin)
		}
	default:
		err = errors.New("Unable to determine ProjectType")
	}

	if err == nil {
		return true, pid, nil //fmt.Printf ("Pid:%d, Pname:%s\n", pid, s)
	} else if err.Error() == "not found" {
		return false, 0, nil
	} else {
		return false, 0, err
	}
}

// IsAppServerRunning - Will then work out what wallet this relates to, and return accurate info
func IsAppServerRunning() (bool, int, error) {
	var pid int
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return false, pid, err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			pid, _, err = findProcess(CAppServerFileWinGoDivi)
		} else {
			pid, _, err = findProcess(CAppServerFileGoDivi)
		}
	case PTPIVX:
		if runtime.GOOS == "windows" {
			pid, _, err = findProcess(CAppServerFileWinGoPIVX)
		} else {
			pid, _, err = findProcess(CAppServerFileGoPIVX)
		}
	case PTTrezarcoin:
		if runtime.GOOS == "windows" {
			pid, _, err = findProcess(CAppServerFileWinGoTrezarcoin)
		} else {
			pid, _, err = findProcess(CAppServerFileGoTrezarcoin)
		}
	default:
		err = errors.New("Unable to determine ProjectType")
	}

	if err == nil {
		return true, pid, nil //fmt.Printf ("Pid:%d, Pname:%s\n", pid, s)
	} else if err.Error() == "not found" {
		return false, 0, nil
	} else {
		return false, 0, err
	}
}

// RunAppServer - Runs the App Server
func RunAppServer(displayOutput bool) error {
	idr, _, _ := IsAppServerRunning()
	if idr == true {
		// Already running...
		return nil
	}
	gwconf, err := GetCLIConfigStruct(false)
	if err != nil {
		return err
	}
	abf, _ := GetAppsBinFolder()

	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			fp := abf + CAppServerFileWinGoDivi
			cmd := exec.Command("cmd.exe", "/C", "start", "/b", fp)
			if err := cmd.Run(); err != nil {
				return err
			}
		} else {
			if displayOutput {
				fmt.Println("Attempting to run " + CAppNameServerGoDivi + "...")
			}

			cmdRun := exec.Command(abf + CAppServerFileGoDivi)
			if err := cmdRun.Start(); err != nil {
				return fmt.Errorf("Failed to start cmd: %v", err)
			}
		}
	case PTPIVX:
		if runtime.GOOS == "windows" {
			fp := abf + CAppServerFileWinGoPIVX
			cmd := exec.Command("cmd.exe", "/C", "start", "/b", fp)
			if err := cmd.Run(); err != nil {
				return err
			}
		} else {
			if displayOutput {
				fmt.Println("Attempting to run " + CAppNameServerGoPIVX + "...")
			}

			cmdRun := exec.Command(abf + CAppServerFileGoPIVX)
			if err := cmdRun.Start(); err != nil {
				return fmt.Errorf("Failed to start cmd: %v", err)
			}
		}
	case PTTrezarcoin:
		if runtime.GOOS == "windows" {
			fp := abf + CAppServerFileWinGoTrezarcoin
			cmd := exec.Command("cmd.exe", "/C", "start", "/b", fp)
			if err := cmd.Run(); err != nil {
				return err
			}
		} else {
			if displayOutput {
				fmt.Println("Attempting to run " + CAppNameServerGoTrezarcoin + "...")
			}

			cmdRun := exec.Command(abf + CAppServerFileGoTrezarcoin)
			if err := cmdRun.Start(); err != nil {
				return fmt.Errorf("Failed to start cmd: %v", err)
			}
		}
	default:
		err = errors.New("Unable to determine ProjectType")
	}

	return nil
}

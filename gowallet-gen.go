package gwcommon

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/inconshreveable/go-update"
	"github.com/mitchellh/go-ps"
)

const (
	// CAppVersion - The app version of the suite of apps
	CAppVersion    string = "0.21.1" // All of the individual apps will have the same version to make it easier for the user
	cUnknown       string = "Unknown"
	cAddNodeURL    string = "https://api.diviproject.org/v1/addnode"
	cDownloadURLDP string = "https://github.com/DiviProject/Divi/releases/download/v1.0.8/"
	// CDownloadURLGD - The download file lotcation for GoDivi
	CDownloadURLGD string = "https://bitbucket.org/rmace/godivi/downloads/"
	cCoinNameDivi  string = "Divi"

	// Divi Wallet Constants
	CDiviAppVersion string = "1.08"
	cDiviHomeDir    string = ".divi"
	cDiviHomeDirWin string = "DIVI"
	cDiviBinDir     string = "godivi"
	cDiviBinDirWin  string = "GoDivi"

	// Public download files
	CDFileGodiviLatetsARM     string = "godivi-arm-latest.zip"
	CDFileGodiviLatetsLinux   string = "godivi-linux-latest.zip"
	CDFileGodiviLatetsWindows string = "godivi-windows-latest.zip"
	cDFileRPi                 string = "divi-1.0.8-RPi2.tar.gz"
	cDFileUbuntu              string = "divi-1.0.8-x86_64-linux-gnu.tar.gz"
	cDFileWindows             string = "divi-1.0.8-win64.zip"

	// Divi project file constants
	CAppNameGoDivi        string = "GoDivi"
	CAppNameUpdaterGoDivi string = "GoDivi Updater"
	CAppNameCLIGoDivi     string = "GoDivi CLI"
	CAppNameServerGoDivi  string = "GoDivi Server"

	cDiviConfFile   string = "divi.conf"
	cDiviCliFile    string = "divi-cli"
	cDiviCliFileWin string = "divi-cli.exe"
	cDiviDFile      string = "divid"
	cDiviDFileWin   string = "divid.exe"
	cDiviTxFile     string = "divi-tx"
	cDiviTxFileWin  string = "divi-tx.exe"

	// GoWallet file constants
	CAppCLIFileCompiled        string = "cli"
	CAppCLIFileCompiledWin     string = "cli.exe"
	CAppServerFileCompiled     string = "web"
	CAppServerFileCompiledWin  string = "web.exe"
	CAppUpdaterFileCompiled    string = "updater"
	CAppUpdaterFileCompiledWin string = "updater.exe"

	// GoDivi file constants
	CAppCLIFileGoDivi             string = "godivi"
	CAppCLIFileWinGoDivi          string = "godivi.exe"
	CAppCLIFileInstallerGoDivi    string = "godivi-installer"
	CAppCLIFileInstallerWinGoDivi string = "godivi-installer.exe"
	CAppServerFileGoDivi          string = "godivis"
	CAppServerFileWinGoDivi       string = "godivis.exe"
	CAppUpdaterFileGoDivi         string = "update-godivi"
	CAppUpdaterFileWinGoDivi      string = "update-godivi.exe"
	CAppCLILogfileGoDivi          string = "godivi.log"
	CWalletSeedFileGoDivi         string = "unsecure-divi-seed.txt"

	// Divi-cli command constants
	cCommandGetBCInfo     string = "GetBlockchainInfo"
	cCommandGetWInfo      string = "GetWalletInfo"
	cCommandMNSyncStatus1 string = "mnsync"
	cCommandMNSyncStatus2 string = "status"

	// Divii-cli wallet commands
	cCommandDisplayWalletAddress string = "getaddressesbyaccount" // ./divi-cli getaddressesbyaccount ""
	cCommandDumpHDinfo           string = "dumphdinfo"            // ./divi-cli dumphdinfo
	// CCommandEncryptWallet - Needed by dash command
	CCommandEncryptWallet  string = "encryptwallet"    // ./divi-cli encryptwallet “a_strong_password”
	cCommandRestoreWallet  string = "-hdseed="         // ./divid -debug-hdseed=the_seed -rescan (stop divid, rename wallet.dat, then run command)
	cCommandUnlockWallet   string = "walletpassphrase" // ./divi-cli walletpassphrase “password” 0
	cCommandUnlockWalletFS string = "walletpassphrase" // ./divi-cli walletpassphrase “password” 0 true
	cCommandLockWallet     string = "walletlock"       // ./divi-cli walletlock

	cRPCUserStr     string = "rpcuser"
	cRPCPasswordStr string = "rpcpassword"

	// Divid Responses
	cDiviDNotRunningError     string = "error: couldn't connect to server"
	cDiviDDIVIServerStarting  string = "DIVI server starting"
	cDividRespWalletEncrypted string = "wallet encrypted"

	cGoDiviExportPath         string = "export PATH=$PATH:"
	CUninstallConfirmationStr string = "Confirm"
	cSeedStoredSafelyStr      string = "Confirm"

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
	// PTTrezar - TrezarCoin
	PTTrezar
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
	CServRequestShutdownServer string = "ShutdownServer"
)

type BlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Verificationprogress float64 `json:"verificationprogress"`
	Chainwork            string  `json:"chainwork"`
}

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

type MNSyncStatus struct {
	IsBlockchainSynced         bool `json:"IsBlockchainSynced"`
	LastMasternodeList         int  `json:"lastMasternodeList"`
	LastMasternodeWinner       int  `json:"lastMasternodeWinner"`
	LastFailure                int  `json:"lastFailure"`
	NCountFailures             int  `json:"nCountFailures"`
	SumMasternodeList          int  `json:"sumMasternodeList"`
	SumMasternodeWinner        int  `json:"sumMasternodeWinner"`
	CountMasternodeList        int  `json:"countMasternodeList"`
	CountMasternodeWinner      int  `json:"countMasternodeWinner"`
	RequestedMasternodeAssets  int  `json:"RequestedMasternodeAssets"`
	RequestedMasternodeAttempt int  `json:"RequestedMasternodeAttempt"`
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

// WalletInfoStruct - The WalletInfoStruct
type WalletInfoStruct struct {
	Walletversion      int     `json:"walletversion"`
	Balance            float64 `json:"balance"`
	UnconfirmedBalance float64 `json:"unconfirmed_balance"`
	ImmatureBalance    float64 `json:"immature_balance"`
	Txcount            int     `json:"txcount"`
	Keypoololdest      int     `json:"keypoololdest"`
	Keypoolsize        int     `json:"keypoolsize"`
	UnlockedUntil      int     `json:"unlocked_until"`
	EncryptionStatus   string  `json:"encryption_status"`
	Hdchainid          string  `json:"hdchainid"`
	Hdaccountcount     int     `json:"hdaccountcount"`
	Hdaccounts         []struct {
		Hdaccountindex     int `json:"hdaccountindex"`
		Hdexternalkeyindex int `json:"hdexternalkeyindex"`
		Hdinternalkeyindex int `json:"hdinternalkeyindex"`
	} `json:"hdaccounts"`
}

var lastBCSyncStatus string = ""
var lastMNSyncStatus string = ""

func AddGoDiviPath() error {
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

// DoRequiredFiles - Download and install required files
func DoRequiredFiles() error {
	var filePath, fileURL string
	abf, err := GetAppsBinFolder()
	if err != nil {
		return fmt.Errorf("Unable to perform GetAppsBinFolder: %v ", err)
	}

	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return fmt.Errorf("Unable to get ConfigStruct: %v ", err)
	}
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			filePath = abf + cDFileWindows
			fileURL = cDownloadURLDP + cDFileWindows
		} else if runtime.GOARCH == "arm" {
			filePath = abf + cDFileRPi
			fileURL = cDownloadURLDP + cDFileRPi
		} else {
			filePath = abf + cDFileUbuntu
			fileURL = cDownloadURLDP + cDFileUbuntu
		}
	}

	log.Print("Downloading required files...")
	if err := DownloadFile(filePath, fileURL); err != nil {
		return fmt.Errorf("Unable to download file: %v - %v", filePath+fileURL, err)
	}

	r, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Unable to open file: %v - %v", filePath, err)
	}

	// Now, uncompress the files...
	log.Print("Uncompressing files...")
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			_, err = UnZip(filePath, "tmp")
			if err != nil {
				return fmt.Errorf("Unable to unzip file: %v - %v", filePath, err)
			}
		} else if runtime.GOARCH == "arm" {
			err = extractTarGz(r)
			if err != nil {
				return fmt.Errorf("Unable to extractTarGz file: %v - %v", r, err)
			}
		} else {
			err = extractTarGz(r)
			if err != nil {
				return fmt.Errorf("Unable to extractTarGz file: %v - %v", r, err)
			}
		}
	}

	log.Print("Installing files...")

	// Copy files to correct location
	var srcPath, srcRoot, srcFileCLI, srcFileD, srcFileTX, srcFileGD, srcFileUGD, srcFileGDS string
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			srcPath = "./tmp/divi-1.0.8/bin/"
			srcRoot = "./tmp/"
			srcFileCLI = cDiviCliFileWin
			srcFileD = cDiviDFileWin
			srcFileTX = cDiviTxFileWin
			srcFileGD = CAppCLIFileWinGoDivi
			srcFileGDS = CAppServerFileWinGoDivi

		} else if runtime.GOARCH == "arm" {
			srcPath = "./divi-1.0.8/bin/"
			srcRoot = "./divi-1.0.8/"
			srcFileCLI = cDiviCliFile
			srcFileD = cDiviDFile
			srcFileTX = cDiviTxFile
			srcFileGD = CAppCLIFileGoDivi
			srcFileUGD = CAppUpdaterFileGoDivi
			srcFileGDS = CAppServerFileGoDivi

		} else {
			srcPath = "./divi-1.0.8/bin/"
			srcRoot = "./divi-1.0.8/"
			srcFileCLI = cDiviCliFile
			srcFileD = cDiviDFile
			srcFileTX = cDiviTxFile
			srcFileGD = CAppCLIFileGoDivi
			srcFileUGD = CAppUpdaterFileGoDivi
			srcFileGDS = CAppServerFileGoDivi
		}
	}

	// coin-cli
	err = FileCopy(srcPath+srcFileCLI, abf+srcFileCLI, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", srcPath+srcFileCLI, err)
	}
	err = os.Chmod(abf+srcFileCLI, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", abf+srcFileCLI, err)
	}
	// coind
	err = FileCopy(srcPath+srcFileD, abf+srcFileD, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", srcPath+srcFileD, err)
	}
	err = os.Chmod(abf+srcFileD, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", abf+srcFileD, err)
	}

	// cointx
	err = FileCopy(srcPath+srcFileTX, abf+srcFileTX, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", srcPath+srcFileTX, err)
	}
	err = os.Chmod(abf+srcFileTX, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", abf+srcFileTX, err)
	}

	// Copy the gowallet binary itself
	ex, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting exe - %v", err)
	}

	err = FileCopy(ex, abf+srcFileGD, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", abf+srcFileGD, err)
	}
	err = os.Chmod(abf+srcFileGD, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", abf+srcFileGD, err)
	}

	// Copy the updater file
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			// TODO Code the Windows part
		} else if runtime.GOARCH == "arm" {
			err = FileCopy("./"+CAppUpdaterFileGoDivi, abf+srcFileUGD, false)
			if err != nil {
				return fmt.Errorf("Unable to copyFile: %v - %v", abf+srcFileUGD, err)
			}
			err = os.Chmod(abf+srcFileUGD, 0777)
			if err != nil {
				return fmt.Errorf("Unable to chmod file: %v - %v", abf+srcFileUGD, err)
			}

		} else {
			err = FileCopy("./"+CAppUpdaterFileGoDivi, abf+srcFileUGD, false)
			if err != nil {
				return fmt.Errorf("Unable to copyFile: %v - %v", abf+srcFileUGD, err)
			}
			err = os.Chmod(abf+srcFileUGD, 0777)
			if err != nil {
				return fmt.Errorf("Unable to chmod file: %v - %v", abf+srcFileUGD, err)
			}
		}
	}

	// Copy the godivis file
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			// TODO Code the Windows part
		} else if runtime.GOARCH == "arm" {
			err = FileCopy("./"+CAppServerFileGoDivi, abf+srcFileGDS, false)
			if err != nil {
				return fmt.Errorf("Unable to copyFile: %v - %v", abf+srcFileGDS, err)
			}
			err = os.Chmod(abf+srcFileGDS, 0777)
			if err != nil {
				return fmt.Errorf("Unable to chmod file: %v - %v", abf+srcFileGDS, err)
			}

		} else {
			err = FileCopy("./"+CAppServerFileGoDivi, abf+srcFileGDS, false)
			if err != nil {
				return fmt.Errorf("Unable to copyFile: %v - %v", abf+srcFileGDS, err)
			}
			err = os.Chmod(abf+srcFileGDS, 0777)
			if err != nil {
				return fmt.Errorf("Unable to chmod file: %v - %v", abf+srcFileGDS, err)
			}
		}
	}

	err = os.RemoveAll(srcRoot)
	if err != nil {
		return fmt.Errorf("error performing os.RemoveAll - %v", err)
	}
	err = FileDelete(filePath)
	if err != nil {
		return fmt.Errorf("error deleting file: %v - %v", filePath, err)
	}

	return nil
}

// DoPrivKeyDisplay - Display the private key
func DoPrivKeyDisplay() error {
	dbf, err := GetAppsBinFolder()
	if err != nil {
		return fmt.Errorf("Unable to GetAppsBinFolder: %v", err)
	}
	// Display instructions for displaying seed recovery

	sWarning := getWalletSeedDisplayWarning()
	fmt.Printf(sWarning)
	fmt.Println("")
	fmt.Println("\nRequesting private seed...")
	s, err := runDCCommand(dbf+cDiviCliFile, cCommandDumpHDinfo, "Waiting for wallet to respond. This could take several minutes...", 30)
	if err != nil {
		return fmt.Errorf("Unable to run command: %v - %v", dbf+cDiviCliFile+cCommandDumpHDinfo, err)
	}

	fmt.Println("\nPrivate seed received...")
	fmt.Println("")
	println(s)

	return nil
}

func DoPrivKeyFile() error {
	dbf, err := GetAppsBinFolder()
	if err != nil {
		return fmt.Errorf("Unable to GetAppsBinFolder: %v", err)
	}

	// User specified that they wanted to save their private key in a file.
	s := getWalletSeedDisplayWarning() + `

Storing your private key in a file is risky.

Please confirm that you understand the risks: `
	yn := getYesNoResp(s)
	if yn == "y" {
		fmt.Println("\nRequesting private seed...")
		s, err := runDCCommand(dbf+cDiviCliFile, cCommandDumpHDinfo, "Waiting for wallet to respond. This could take several minutes...", 30)
		// cmd := exec.Command(dbf+cDiviCliFile, cCommandDumpHDinfo)
		// out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Unable to run command: %v - %v", dbf+cDiviCliFile+cCommandDumpHDinfo, err)
		}

		// Now store the info in file
		err = WriteTextToFile(dbf+CWalletSeedFileGoDivi, s)
		if err != nil {
			return fmt.Errorf("error writing to file %s", err)
		}
		fmt.Println("Now please store the privte seed file somewhere safe. The file has been saved to: " + dbf + CWalletSeedFileGoDivi)
	}
	return nil
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

func doWalletAdressDisplay() error {
	err := RunCoinDaemon(false)
	if err != nil {
		return fmt.Errorf("Unable to RunCoinDaemon: %v ", err)
	}

	dbf, err := GetAppsBinFolder()
	if err != nil {
		return fmt.Errorf("Unable to GetAppsBinFolder: %v", err)
	}
	// Display wallet public address

	fmt.Println("\nRequesting wallet address...")
	s, err := RunDCCommandWithValue(dbf+cDiviCliFile, cCommandDisplayWalletAddress, `""`, "Waiting for wallet to respond. This could take several minutes...", 30)
	if err != nil {
		return fmt.Errorf("Unable to run command: %v - %v", dbf+cDiviCliFile+cCommandDisplayWalletAddress, err)
	}

	fmt.Println("\nWallet address received...")
	fmt.Println("")
	println(s)

	return nil
}

func findProcess(key string) (int, string, error) {
	pname := ""
	pid := 0
	err := errors.New("not found")
	ps, _ := ps.Processes()

	for i := range ps {
		if ps[i].Executable() == key {
			pid = ps[i].Pid()
			pname = ps[i].Executable()
			err = nil
			break
		}
	}
	return pid, pname, err
}

func GetAddNodes() ([]byte, error) {
	addNodesClient := http.Client{
		Timeout: time.Second * 3, // Maximum of 3 secs
	}

	req, err := http.NewRequest(http.MethodGet, cAddNodeURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "godivi")

	res, getErr := addNodesClient.Do(req)
	if getErr != nil {
		return nil, err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	return body, nil
}

func GetBlockchainInfo() (BlockchainInfo, error) {
	// gdConfig, err := getConfStruct("./")
	// if err != nil {
	// 	log.Print(err)
	// }

	bci := BlockchainInfo{}
	dbf, _ := GetAppsBinFolder()

	cmdBCInfo := exec.Command(dbf+cDiviCliFile, cCommandGetBCInfo)
	out, _ := cmdBCInfo.CombinedOutput()
	err := json.Unmarshal([]byte(out), &bci)
	if err != nil {
		return bci, err
	}
	return bci, nil
}

// GetAppsBinFolder - Returns the directory of where the apps binary files are stored
func GetAppsBinFolder() (string, error) {
	var s string
	gwconf, err := GetConfigStruct(false)
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
		}

	} else {
		switch gwconf.ProjectType {
		case PTDivi:
			s = AddTrailingSlash(hd) + AddTrailingSlash(cDiviBinDir)
		}
	}
	return s, nil
}

// GetAppCLIFileName - Returns the name of the app cli binary file e.g. godivi
func GetAppCLIFileName() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppCLIFileGoDivi, nil
	}
	return "", nil
}

// GetAppInstallerFileName - Returns the name of the app installer binary file e.g. "godivi-installer"
func GetAppInstallerFileName() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppCLIFileInstallerGoDivi, nil
	}
	return "", nil
}

// GetAppServerFileName - Returns the name of the app server binary file e.g. godivis
func GetAppServerFileName() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppServerFileGoDivi, nil
	}
	return "", nil
}

// GetAppUpdaterFileName - Returns the name of the app updater binary file e.g. update-godivi
func GetAppUpdaterFileName() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppUpdaterFileGoDivi, nil
	}
	return "", nil
}

// GetAppCLIName - Returns the application CLI name e.g. GoDivi CLI
func GetAppCLIName() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppNameCLIGoDivi, nil
	}
	return "", nil
}

// GetAppServerName - Returns the application Server name e.g. GoDivi Server
func GetAppServerName() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppNameServerGoDivi, nil
	}
	return "", nil
}

// GetAppName - Returns the application name e.g. GoDivi
func GetAppName() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return CAppNameGoDivi, nil
	}
	return "", nil
}

// GetCoinDaemonFilename - Return the coin daemon file name e.g. divid
func GetCoinDaemonFilename() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return cDiviDFile, nil
	}
	return "", nil
}

// GetCoinHomeFolder - Returns the ome folder for the coin e.g. .divi
func GetCoinHomeFolder() (string, error) {
	var s string
	gwconf, err := GetConfigStruct(false)
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
			s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(cDiviHomeDirWin)
		}
	} else {
		switch gwconf.ProjectType {
		case PTDivi:
			s = AddTrailingSlash(hd) + AddTrailingSlash(cDiviHomeDir)
		}
	}
	return s, nil
}

// GetCoinName - Returns the name of the coin e.g. Divi
func GetCoinName() (string, error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		return cCoinNameDivi, nil
	}
	return "", nil
}

// GetCoinDownloadLink - Returns a link to the required file
func GetCoinDownloadLink(ostype OSType) (url, file string, err error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		switch ostype {
		case OSTArm:
			return cDownloadURLDP, cDFileRPi, nil
		case OSTLinux:
			return cDownloadURLDP, cDFileUbuntu, nil
		case OSTWindows:
			return cDownloadURLDP, cDFileWindows, nil
		}
	}
	return "", "", nil
}

// GetGoWalletDownloadLink - Returns a link of both the url and file
func GetGoWalletDownloadLink(ostype OSType) (url, file string, err error) {
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return "", "", err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		switch ostype {
		case OSTArm:
			return CDownloadURLGD, CDFileGodiviLatetsARM, nil
		case OSTLinux:
			return CDownloadURLGD, CDFileGodiviLatetsLinux, nil
		case OSTWindows:
			return CDownloadURLGD, CDFileGodiviLatetsWindows, nil
		}
	}
	return "", "", nil
}

func GetMNSyncStatus() (MNSyncStatus, error) {
	// gdConfig, err := getConfStruct("./")
	// if err != nil {
	// 	log.Print(err)
	// }

	mnss := MNSyncStatus{}
	dbf, _ := GetAppsBinFolder()

	cmdMNSS := exec.Command(dbf+cDiviCliFile, cCommandMNSyncStatus1, cCommandMNSyncStatus2)
	out, _ := cmdMNSS.CombinedOutput()
	err := json.Unmarshal([]byte(out), &mnss)
	if err != nil {
		return mnss, err
	}
	return mnss, nil
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

func GetWalletAddress(attempts int) (string, error) {
	var err error
	var s string = "waiting for wallet."
	dbf, _ := GetAppsBinFolder()
	app := dbf + cDiviCliFile
	arg1 := cCommandDisplayWalletAddress
	arg2 := ""

	for i := 0; i < attempts; i++ {

		cmd := exec.Command(app, arg1, arg2)
		out, err := cmd.CombinedOutput()

		if err == nil {
			return string(out), err
		}

		fmt.Printf("\r"+s+" %d/"+strconv.Itoa(attempts), i+1)

		time.Sleep(3 * time.Second)

		// t := string(out)
		// if strings.Contains(string(out), "Loading block index....") {
		// 	//s = s + "."
		// 	//fmt.Println(s)
		// 	fmt.Printf("\r"+s+" %d/"+strconv.Itoa(attempts), i+1)
		// 	fmt.Println(t)

		// 	time.Sleep(3 * time.Second)

		// }

	}

	return "", err

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

func GetWalletInfo(dispProgress bool) (WalletInfoStruct, error) {
	wi := WalletInfoStruct{}
	s := "waiting for divid server.."
	attempts := 30

	// Start the DiviD server if required...
	err := RunCoinDaemon(false)
	if err != nil {
		return wi, fmt.Errorf("Unable to RunDiviD: %v ", err)
	}

	dbf, err := GetAppsBinFolder()
	if err != nil {
		return wi, fmt.Errorf("Unable to GetAppsBinFolder: %v ", err)
	}

	for i := 0; i < attempts; i++ {
		cmd := exec.Command(dbf+cDiviCliFile, cCommandGetWInfo)
		out, err := cmd.CombinedOutput()
		if err == nil {
			errUM := json.Unmarshal([]byte(out), &wi)
			if errUM == nil {
				return wi, err
			}
		}

		//s = s + "."
		//fmt.Println(s)
		if dispProgress {
			fmt.Printf("\r"+s+" %d/"+strconv.Itoa(attempts), i+1)
		}
		time.Sleep(3 * time.Second)
	}
	return wi, nil
}

func getWalletRestoreResp() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(`Warning - This will overrite your existing wallet.dat file and re-sync the blockchain!

It will take a while for your restored wallet to sync and display any funds.

Restore wallet now?: (y/n)`)
	resp, _ := reader.ReadString('\n')
	resp = strings.ReplaceAll(resp, "\n", "")
	return resp
}

func getWalletSeedDisplayWarning() string {
	return `
A recovery seed can be used to recover your wallet, should anything happen to this computer.
					
It's a good idea to have more than one and keep each in a safe place, other than your computer.`
}

func GetWalletSeedRecoveryResp() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n\n*** WARNING ***" + "\n\n" +
		"You haven't provided confirmation that you've backed up your recovery seed!\n\n" +
		"This is *extremely* important as it's the only way of recovering your wallet in the future\n\n" +
		"To (d)isplay your reovery seed now press: d, to (c)onfirm that you've backed it up press: c, or to (m)ove on, press: m\n\n" +
		"Please enter: [d/c/m]")
	resp, _ := reader.ReadString('\n')
	resp = strings.ReplaceAll(resp, "\n", "")
	return resp
}

func GetWalletSeedRecoveryConfirmationResp() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter the response: " + cSeedStoredSafelyStr)
	resp, _ := reader.ReadString('\n')
	if resp == cSeedStoredSafelyStr+"\n" {
		return true
	}

	return false
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

// IsCoinDaemonRunning - Works out whether the coin Daemon is running e.g. divid
func IsCoinDaemonRunning() (bool, int, error) {
	var pid int
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return false, pid, err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			pid, _, err = findProcess(cDiviDFileWin)
		} else {
			pid, _, err = findProcess(cDiviDFile)
		}
	}

	if err == nil {
		return true, pid, nil //fmt.Printf ("Pid:%d, Pname:%s\n", pid, s)
	} else {
		return false, 0, err
	}
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

// IsGoWalletCLIRunning - Is the GoWallet CLI Running
func IsGoWalletCLIRunning() (bool, int, error) {
	var pid int
	var err error
	if runtime.GOOS == "windows" {
		pid, _, err = findProcess(CAppCLIFileWinGoDivi)
	} else {
		pid, _, err = findProcess(CAppCLIFileGoDivi)
	}

	//pid, _, err := FindProcess(cDiviDFile)
	if err.Error() == "not found" {
		return false, 0, nil
	}
	if err == nil {
		return true, pid, nil //fmt.Printf ("Pid:%d, Pname:%s\n", pid, s)
	} else {
		return false, 0, err
	}
}

// IsAppCLIRunning - Will then work out what wallet this relates to, and return bool whether the CLI app is running
func IsAppCLIRunning() (bool, int, error) {
	var pid int
	gwconf, err := GetConfigStruct(false)
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
	gwconf, err := GetConfigStruct(false)
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
	}

	if err == nil {
		return true, pid, nil //fmt.Printf ("Pid:%d, Pname:%s\n", pid, s)
	} else if err.Error() == "not found" {
		return false, 0, nil
	} else {
		return false, 0, err
	}
}

func runDCCommand(cmdBaseStr, cmdStr, waitingStr string, attempts int) (string, error) {
	var err error
	//var s string = waitingStr
	for i := 0; i < attempts; i++ {
		cmd := exec.Command(cmdBaseStr, cmdStr)
		out, err := cmd.CombinedOutput()

		// cmd := exec.Command(cmdBaseStr, cmdStr)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// err = cmd.Run()

		if err == nil {
			return string(out), err
		}
		//s = s + "."
		//fmt.Println(s)
		fmt.Printf("\r"+waitingStr+" %d/"+strconv.Itoa(attempts), i)

		time.Sleep(3 * time.Second)
	}

	return "", err
}

func RunDCCommandWithValue(cmdBaseStr, cmdStr, valueStr, waitingStr string, attempts int) (string, error) {
	var err error
	//var s string = waitingStr
	for i := 0; i < attempts; i++ {
		cmd := exec.Command(cmdBaseStr, cmdStr, valueStr)
		out, err := cmd.CombinedOutput()

		if err == nil {
			return string(out), err
		}
		//s = s + "."
		//fmt.Println(s)
		fmt.Printf("\r"+waitingStr+" %d/"+strconv.Itoa(attempts), i)
		time.Sleep(3 * time.Second)
	}

	return "", err
}

// RunCoinDaemon - Run the coins Daemon e.g. Run divid
func RunCoinDaemon(displayOutput bool) error {
	idr, _, _ := IsCoinDaemonRunning()
	if idr == true {
		// Already running...
		return nil
	}

	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return err
	}
	abf, _ := GetAppsBinFolder()

	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			//_ = exec.Command(GetAppsBinFolder() + cDiviDFileWin)
			fp := abf + cDiviDFileWin
			cmd := exec.Command("cmd.exe", "/C", "start", "/b", fp)
			if err := cmd.Run(); err != nil {
				return err
			}

		} else {
			if displayOutput {
				fmt.Println("Attempting to run the divid daemon...")
			}

			cmdRun := exec.Command(abf + cDiviDFile)
			stdout, err := cmdRun.StdoutPipe()
			if err != nil {
				return err
			}
			cmdRun.Start()

			buf := bufio.NewReader(stdout) // Notice that this is not in a loop
			num := 1
			for {
				line, _, _ := buf.ReadLine()
				if num > 3 {
					os.Exit(0)
				}
				num++
				if string(line) == "DIVI server starting" {
					return nil
				} else {
					return errors.New("Unable to start Divi server")
				}
			}
		}
	}

	return nil
}

// RunAppServer - Runs the App Server
func RunAppServer(displayOutput bool) error {
	idr, _, _ := IsAppServerRunning()
	if idr == true {
		// Already running...
		return nil
	}
	gwconf, err := GetConfigStruct(false)
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
	}

	return nil
}

func RunInitialDaemon() error {
	abf, err := GetAppsBinFolder()
	if err != nil {
		return fmt.Errorf("Unable to GetAppsBinFolder - %v", err)
	}
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return fmt.Errorf("Unable to GetConfigStruct - %v", err)
	}
	switch gwconf.ProjectType {
	case PTDivi:
		//Run divid for the first time, so that we can get the outputted info to build the conf file
		fmt.Println("About to run divid for the first time...")
		cmdDividRun := exec.Command(abf + cDiviDFile)
		out, err := cmdDividRun.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Unable to run "+abf+cDiviDFile+" - %v", err)
		}
		fmt.Println("Populating " + cDiviConfFile + " for initial setup...")

		scanner := bufio.NewScanner(strings.NewReader(string(out)))
		var rpcuser, rpcpw string
		for scanner.Scan() {
			s := scanner.Text()
			if strings.Contains(s, cRPCUserStr) {
				rpcuser = strings.ReplaceAll(s, cRPCUserStr+"=", "")
			}
			if strings.Contains(s, cRPCPasswordStr) {
				rpcpw = strings.ReplaceAll(s, cRPCPasswordStr+"=", "")
			}
		}

		chd, _ := GetCoinHomeFolder()

		err = WriteTextToFile(chd+cDiviConfFile, cRPCUserStr+"="+rpcuser)
		if err != nil {
			log.Fatal(err)
		}
		err = WriteTextToFile(chd+cDiviConfFile, cRPCPasswordStr+"="+rpcpw)
		if err != nil {
			log.Fatal(err)
		}
		err = WriteTextToFile(chd+cDiviConfFile, "")
		if err != nil {
			log.Fatal(err)
		}
		err = WriteTextToFile(chd+cDiviConfFile, "daemon=1")
		if err != nil {
			log.Fatal(err)
		}
		err = WriteTextToFile(chd+cDiviConfFile, "")
		if err != nil {
			log.Fatal(err)
		}

		// Now get a list of the latest "addnodes" and add them to the file:
		// I've commented out the below, as I think it might cause future issues with blockchain syncing,
		// because, I think that the ipaddresess in the conf file are used before any others are picked up,
		// so, it's possible that they could all go, and then cause issues.

		// gdc.AddToLog(lfp, "Adding latest master nodes to "+gdc.CDiviConfFile)
		// addnodes, _ := gdc.GetAddNodes()
		// sAddnodes := string(addnodes[:])
		// gdc.WriteTextToFile(dhd+gdc.CDiviConfFile, sAddnodes)

		return nil
	}
	return nil
}

// StopCoinDaemon - Stops the coin daemon (e.g. divid) from running
func StopCoinDaemon() error {
	idr, _, _ := IsCoinDaemonRunning() //DiviDRunning()
	if idr != true {
		// Not running anyway ...
		return nil
	}

	dbf, _ := GetAppsBinFolder()
	gwconf, err := GetConfigStruct(false)
	if err != nil {
		return err
	}
	switch gwconf.ProjectType {
	case PTDivi:
		if runtime.GOOS == "windows" {
			// TODO Complete for Windows
		} else {
			cRun := exec.Command(dbf+cDiviCliFile, "stop")
			if err := cRun.Run(); err != nil {
				return fmt.Errorf("Unable to StopDiviD:%v", err)
			}

			for i := 0; i < 50; i++ {
				sr, _, _ := IsCoinDaemonRunning() //DiviDRunning()
				if !sr {
					return nil
				}
				fmt.Printf("\rWaiting for divid server to stop %d/"+strconv.Itoa(50), i+1)
				time.Sleep(3 * time.Second)

			}
		}
	}

	return nil
}

func UnlockWallet(pword string, attempts int, forStaking bool) (bool, error) {
	var err error
	var s string = "waiting for wallet."
	dbf, _ := GetAppsBinFolder()
	app := dbf + cDiviCliFile
	arg1 := cCommandUnlockWalletFS
	arg2 := pword
	arg3 := "0"
	arg4 := "true"
	for i := 0; i < attempts; i++ {

		var cmd *exec.Cmd
		if forStaking {
			cmd = exec.Command(app, arg1, arg2, arg3, arg4)
		} else {
			cmd = exec.Command(app, arg1, arg2, arg3)
		}
		//fmt.Println("cmd = " + dbf + cDiviCliFile + cCommandUnlockWalletFS + `"` + pword + `"` + "0")
		out, err := cmd.CombinedOutput()

		fmt.Println("string = " + string(out))
		//fmt.Println("error = " + err.Error())

		if err == nil {
			return true, err
		}

		if strings.Contains(string(out), "The wallet passphrase entered was incorrect.") {
			return false, err
		}

		if strings.Contains(string(out), "Loading block index....") {
			//s = s + "."
			//fmt.Println(s)
			fmt.Printf("\r"+s+" %d/"+strconv.Itoa(attempts), i+1)

			time.Sleep(3 * time.Second)

		}

	}

	return false, err
}

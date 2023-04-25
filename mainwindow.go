// main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/dialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitee.com/rocket049/discover-go"
)

var console *widget.Entry

func consoleAppend(s string) {
	old := strings.TrimSpace(console.Text)
	console.SetText(old + "\n" + s)

}

func runServer() {
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(mainPage))
	})

	http.ListenAndServe(":6868", nil)
}

func main() {
	var share1 = flag.String("share", "", "share path")
	var upload1 = flag.String("upload", "", "upload path")
	flag.Parse()

	loadConf()
	if *share1 != "" {
		setShareDir(*share1)
	}
	if *upload1 != "" {
		setUploadDir(*upload1)
	}
	defer autoSaveConfs()
	go runServer()
	createGui()
}
func autoSaveConfs() {
	saveConf(uploadDir, share.Get())
}
func saveConf(uploadPath, sharePath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cfgDir := filepath.Join(home, ".config", "FyneLanShare")
	os.MkdirAll(cfgDir, os.ModePerm)
	fp, err := os.Create(filepath.Join(cfgDir, "paths.json"))
	if err != nil {
		return err
	}
	defer fp.Close()
	var cfg struct {
		UploadPath string
		SharePath  string
	}
	cfg.UploadPath = uploadPath
	cfg.SharePath = sharePath
	data, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}
	_, err = fp.Write(data)
	return err
}

func loadConf() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(home, ".config", "FyneLanShare", "paths.json")
	data, err := ioutil.ReadFile(cfgPath)
	var cfg struct {
		UploadPath string
		SharePath  string
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		cfg.SharePath = home
		cfg.UploadPath = home
	}
	setUploadDir(cfg.UploadPath)
	setShareDir(cfg.SharePath)
	return nil
}

func createGui() {
	app := app.New()
	app.Settings().SetTheme(GetMyTheme())
	w := app.NewWindow("File Server")

	console = widget.NewMultiLineEntry()
	console.SetMinRowsVisible(5)

	labelShare := widget.NewLabel(share.Get())
	btShare := widget.NewButton("Share Path", func() {
		button1 := dialog.NewFolderOpen(func(url1 fyne.ListableURI, err error) {
			if err == nil {
				labelShare.SetText(url1.Path())
				share.Set(url1.Path())
			} else {
				//consoleAppend( err.Error())
				log.Println(err.Error())
			}
		}, w)
		button1.Show()
	})

	box1 := container.NewHBox(btShare, labelShare)

	labelUpload := widget.NewLabel(uploadDir)
	btUpload := widget.NewButton("Upload Path", func() {
		button1 := dialog.NewFolderOpen(func(url1 fyne.ListableURI, err error) {
			if err == nil {
				labelUpload.SetText(url1.Path())
				uploadDir = url1.Path()
			} else {
				//consoleAppend( err.Error())
				log.Println(err.Error())
			}
		}, w)
		button1.Show()
	})

	box2 := container.NewHBox(btUpload, labelUpload)

	addrs := getAddress()
	var addrNum = len(addrs)
	if addrNum == 0 {
		panic("Can not get local host IP.")
	}

	qrImage := widget.NewIcon(addrs[0].Image)
	box4 := container.NewGridWrap(fyne.NewSize(400, 400), qrImage)

	qrTitle := widget.NewLabel(fmt.Sprintf("URL: %v  共%v个", addrs[0].Text, addrNum))

	var imgIdx = 0
	btLeft := widget.NewButton("<", func() {
		//qrcodeLabel.SetText("显示前一个二维码")
		imgIdx -= 1
		if imgIdx < 0 {
			imgIdx = addrNum - 1
		}
		qrTitle.SetText(fmt.Sprintf("URL: %v  共%v个", addrs[imgIdx].Text, addrNum))
		qrImage.Resource = addrs[imgIdx].Image
		qrImage.Refresh()
	})
	btRight := widget.NewButton(">", func() {
		//qrcodeLabel.SetText("显示后一个二维码")
		imgIdx = (imgIdx + 1) % addrNum
		qrTitle.SetText(fmt.Sprintf("URL: %v  共%v个", addrs[imgIdx].Text, addrNum))
		qrImage.Resource = addrs[imgIdx].Image
		qrImage.Refresh()
	})
	box3 := container.NewHBox(btLeft, qrTitle, btRight)

	scanServers := func() {
		client := discover.NewClient()
		services := client.Query()
		consoleAppend("Found Servers:")
		for i := range services {
			consoleAppend(services[i].Href)
		}
		consoleAppend("OK\n")
	}

	buttonScan := widget.NewButton("Scan Servers", func() {
		scanServers()
	})

	boxOut := container.NewVBox(box1, box2, box3, box4, console, buttonScan)

	w.SetContent(boxOut)

	w.SetOnClosed(func() {
		app.Quit()
	})

	w.ShowAndRun()
}

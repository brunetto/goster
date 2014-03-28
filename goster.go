package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/brunetto/gowut/gwu"
)

func main() {
	var ( 
	err error
	)	

	// Create and build a window
	win := gwu.NewWindow("slpp-ui", "UI for slpp")
	win.Style().SetFullWidth()	
	
	// inFileLoad panel
	inFileLoadPanel := gwu.NewHorizontalPanel()
	inFileLoadLabel := gwu.NewLabel("Filename")
	inFileLoadBox := gwu.NewTextBox("")
	inFileBrowseButton := gwu.NewButton("Browse")
	inFileLoadButton := gwu.NewButton("Load")
	
	// Actions
	inFileLoadBox.AddSyncOnETypes(gwu.ETYPE_KEY_UP)
	
	inFileBrowseButton.AddEHandlerFunc(func(e gwu.Event) {
		inFileLoadBox.SetText("testBrowse")
		e.MarkDirty(inFileLoadBox)
	}, gwu.ETYPE_CLICK)
	
	inFileLoadButton.AddEHandlerFunc(func(e gwu.Event) {
		inFileLoadBox.SetText("testLoad")
		e.MarkDirty(inFileLoadBox)
	}, gwu.ETYPE_CLICK)
	
	// Add components
	inFileLoadPanel.Add(inFileLoadLabel)
	inFileLoadPanel.Add(inFileLoadBox)
	inFileLoadPanel.Add(inFileBrowseButton)
	inFileLoadPanel.Add(inFileLoadButton)
	win.Add(inFileLoadPanel)
	
	// TODO: image frame with image label
	// imgPanel
	imgPanel := gwu.NewHorizontalPanel()	
	axPanel := gwu.NewVerticalPanel()
	axXPanel := gwu.NewVerticalPanel()

	axXLabel := gwu.NewLabel("X-Axis")
	axXOriginButton := gwu.NewButton("Mark X axis origin")
	axXOriginLabel := gwu.NewLabel("X origin at: " + "--, --")
	
	// Actions
	axXOriginButton.AddEHandlerFunc(func(e gwu.Event) {
		xOrigin := func(eType gwu.EventType) (*Position) {
			x, y := eType.Mouse()
			fmt.Println(x, y)
			return &Position{x, y}
		} (gwu.ETYPE_CLICK)
		fmt.Println(xOrigin)
// 		axXOriginLabel.SetText("X origin at: " + xOrigin.Str())
		e.MarkDirty(axXOriginLabel)
	}, gwu.ETYPE_CLICK)
	
	// Add components
	axXPanel.Add(axXLabel)
	axXPanel.Add(axXOriginButton)
	axPanel.Add(axXPanel)
	imgPanel.Add(axPanel)
	win.Add(imgPanel)
	
	pointsPanel := gwu.NewHorizontalPanel()
	win.Add(pointsPanel)
	
	
	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("interface", "localhost:8081")//"localhost:8081")
	server.SetText("Hola!!")
	server.AddWin(win)
	openbrowser("http://localhost:8081/interface")
	if err = server.Start(""); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}


// Try to open url in a browser. Instruct to do so if it fails.
// Copyed from 
// https://raw.githubusercontent.com/mumax/3/778cd36cfe2bcf3ae397f6f03251a7e8393a3348/cmd/mumax3/browser.go
// (Mumax3 code at https://github.com/mumax/3)
func openbrowser(url string) {
	for _, cmd := range browsers {
		err := exec.Command(cmd, url).Start()
		if err == nil {
			fmt.Println("Opening web interface in", cmd)
			return
		}
	}
	fmt.Println("Please open ", url, " in a browser")
}

// list of browsers to try.
var browsers = []string{"google-chrome-stable", "chromium-browser", "x-www-browser", "google-chrome", "firefox", "ie", "iexplore"}

type Position struct {
	X int
	Y int
}

func (p *Position) Str () (string){
	return strconv.Itoa(p.X) + ", " + strconv.Itoa(p.X)
}
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/brunetto/gowut/gwu"
)

func main() {
	var ( 
		err error
		getOrigin bool = false
		getPoints bool = false
		origin *Position
		points []*Position
		pos *Position
		pointsBox gwu.TextBox
	)	

	points = make([]*Position, 0)
	
	// Create and build a window
	win := gwu.NewWindow("goster", "Goster plot extractor")
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
	
	// imgPanel
	imgPanel := gwu.NewHorizontalPanel()
	axPanel := gwu.NewVerticalPanel()

	axLabel := gwu.NewLabel("Axis")
	axOriginButton := gwu.NewButton("Mark axes origin")
	axClearOriginButton := gwu.NewButton("Clear")
	axOriginLabel := gwu.NewLabel("Origin at: " + "--, --")
	
	// Actions
	axOriginButton.AddEHandlerFunc(func(e gwu.Event) {
		getOrigin = true
	}, gwu.ETYPE_CLICK)
	
	axClearOriginButton.AddEHandlerFunc(func(e gwu.Event) {
		getOrigin = false
		axOriginLabel.SetText("Origin at: " + "--, --")
		e.MarkDirty(axOriginLabel)
	}, gwu.ETYPE_CLICK)
	
	plotImg := gwu.NewImage("Test", "file:///home/ziosi/Dropbox/DCode/go/src/github.com/brunetto/goster/mapelli_progenitor_remnant.png")
	log.Println(plotImg.Style().Size())
	
	plotImg.AddEHandlerFunc(func(e gwu.Event) {
		x, y := e.Mouse()
		if getOrigin == true {
			origin = &Position{x, y}
			axOriginLabel.SetText("Origin at: " + origin.Str())
			getOrigin = false
			e.MarkDirty(axOriginLabel)
		} else if getPoints == true {
			pos = &Position{x/*-origin.X*/, y/*-origin.Y*/}
			points = append(points, pos)
			pointsBox.SetText(pointsBox.Text() + pos.Str() + "\n")
			pointsBox.SetRows(len(points))
			e.MarkDirty(pointsBox)
		}
	}, gwu.ETYPE_CLICK)
	
	// Add components
	axPanel.Add(axLabel)
	axPanel.Add(axOriginButton)
	axPanel.Add(axClearOriginButton)
	axPanel.Add(axOriginLabel)
	imgPanel.Add(axPanel)
	imgPanel.Add(plotImg)
	win.Add(imgPanel)
	
	pointsPanel := gwu.NewVerticalPanel()
	pointsButtonPanel := gwu.NewHorizontalPanel()
	getPointsButton := gwu.NewButton("Get points")
	getPointsButton.AddEHandlerFunc(func(e gwu.Event) {
		getPoints = true
	}, gwu.ETYPE_CLICK)
	
	donePointsButton := gwu.NewButton("Done")
	donePointsButton.AddEHandlerFunc(func(e gwu.Event) {
		getPoints = false
	}, gwu.ETYPE_CLICK)
	
	clearPointsButton := gwu.NewButton("Clear")
	clearPointsButton.AddEHandlerFunc(func(e gwu.Event) {
		getPoints = false
		points = make([]*Position, 0)
		pointsBox.SetText("")
		pointsBox.SetRows(0)
		e.MarkDirty(pointsBox)
	}, gwu.ETYPE_CLICK)
	
	savePointsButton := gwu.NewButton("Save points")
	savePointsButton.AddEHandlerFunc(func(e gwu.Event) {
		// Write results
		var outFile *os.File
		if outFile, err = os.Create("points.dat"); err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
		for _, point := range points {
			fmt.Fprintf(outFile,"%v\n", point.Str())
		}
	}, gwu.ETYPE_CLICK)
	
	pointsBox = gwu.NewTextBox("")
	pointsBox.SetReadOnly(true)
	
	pointsButtonPanel.Add(getPointsButton)
	pointsButtonPanel.Add(donePointsButton)
	pointsButtonPanel.Add(clearPointsButton)
	pointsButtonPanel.Add(savePointsButton)
	pointsPanel.Add(pointsButtonPanel)
	pointsPanel.Add(pointsBox)
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
	return strconv.Itoa(p.X) + ", " + strconv.Itoa(p.Y)
}
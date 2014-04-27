package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/brunetto/gowut/gwu"
)

var f *os.File
var imgBytes *bytes.Buffer

func main() {
	var ( 
		err error
		getOrigin bool = false
		getPoints bool = false
		origin *Position
		points []*Position
		pos *Position
		pointsBox gwu.TextBox
		imgName string
		maxY int
	)	

	points = make([]*Position, 0)
	
	// Create and build a window
	win := gwu.NewWindow("goster", "Goster plot extractor")
	win.Style().SetFullWidth()
	win.SetCellPadding(5)

	mainPanel := gwu.NewHorizontalPanel()
	mainPanel.SetVAlign(gwu.VA_TOP)
	win.Add(mainPanel)
	controlsPanel := gwu.NewVerticalPanel()
	plotImg := gwu.NewImage("Please, load a png image", "")
	
	mainPanel.Add(controlsPanel)
	mainPanel.Add(plotImg)
	
	
	
	
	
	
	// inFileLoad panel
	inFileLoadPanel := gwu.NewVerticalPanel()
	inFileBtnsPanel := gwu.NewHorizontalPanel()
	inFileLoadLabel := gwu.NewLabel("Filename")
	inFileLoadLabel.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD)
	inFileLoadBox := gwu.NewTextBox("mapelli_progenitor_remnant.png")
	inFileBrowseButton := gwu.NewButton("Browse")
	inFileLoadButton := gwu.NewButton("Load")
	
	inFileLoadPanel.Add(inFileLoadLabel)
	inFileLoadPanel.Add(inFileLoadBox)
	inFileBtnsPanel.Add(inFileBrowseButton)
	inFileBtnsPanel.Add(inFileLoadButton)
	inFileLoadPanel.Add(inFileBtnsPanel)
	controlsPanel.Add(inFileLoadPanel)
	
	
	
	axPanel := gwu.NewVerticalPanel()
	axBtnsPanel := gwu.NewHorizontalPanel()
	axLabel := gwu.NewLabel("Axis")
	axLabel.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD)
	axOriginButton := gwu.NewButton("Mark axes origin")
	axClearOriginButton := gwu.NewButton("Clear")
	axOriginLabel := gwu.NewLabel("Origin at: " + "--, --")
	
	axPanel.Add(axLabel)
	axBtnsPanel.Add(axOriginButton)
	axBtnsPanel.Add(axClearOriginButton)
	axPanel.Add(axBtnsPanel)
	axPanel.Add(axOriginLabel)
	controlsPanel.Add(axPanel)
	
	pointsPanel := gwu.NewVerticalPanel()
	pointsLabel := gwu.NewLabel("Points")
	pointsLabel.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD)
	getPointsButton := gwu.NewButton("Get points")
	donePointsButton := gwu.NewButton("Done")
	pointsButtonPanel := gwu.NewHorizontalPanel()
	clearPointsButton := gwu.NewButton("Clear")
	savePointsButton := gwu.NewButton("Save points")
	pointsBox = gwu.NewTextBox("")
	pointsBox.SetReadOnly(false)
	pointsBox.AddSyncOnETypes(gwu.ETYPE_KEY_UP)
	pointsButtonPanel.Add(getPointsButton)
	pointsButtonPanel.Add(donePointsButton)
	pointsButtonPanel.Add(clearPointsButton)
	pointsButtonPanel.Add(savePointsButton)
	pointsPanel.Add(pointsLabel)
	pointsPanel.Add(pointsButtonPanel)
	pointsPanel.Add(pointsBox)
	controlsPanel.Add(pointsPanel)
	
	clearAllButton := gwu.NewButton("Clear All")
	clearAllButton.SetEnabled(false)
	controlsPanel.Add(clearAllButton)
	
	
	
	// Actions
	inFileLoadBox.AddSyncOnETypes(gwu.ETYPE_KEY_UP)
	
	inFileBrowseButton.AddEHandlerFunc(func(e gwu.Event) {
		inFileLoadBox.SetText("testBrowse")
		e.MarkDirty(inFileLoadBox)
	}, gwu.ETYPE_CLICK)
		
	inFileLoadButton.AddEHandlerFunc(func(e gwu.Event) {
		imgName = inFileLoadBox.Text()
		f, _ = os.Open(imgName)
		defer f.Close()
		img, _ := png.Decode(f)
		// FIXME: inserire i punti in nero sul plot
		imgBytes = new(bytes.Buffer)
		png.Encode(imgBytes, img)
		rect := img.Bounds()
		maxY = rect.Max.Y
		plotImg.SetUrl("http://localhost:8081/"+imgName)
		e.MarkDirty(plotImg)
	}, gwu.ETYPE_CLICK)
	
	axOriginButton.AddEHandlerFunc(func(e gwu.Event) {
		getOrigin = true
		clearAllButton.SetEnabled(true)
		e.MarkDirty(clearAllButton)
	}, gwu.ETYPE_CLICK)
	
	axClearOriginButton.AddEHandlerFunc(func(e gwu.Event) {
		getOrigin = false
		axOriginLabel.SetText("Origin at: " + "--, --")
		e.MarkDirty(axOriginLabel)
	}, gwu.ETYPE_CLICK)
	
	getPointsButton.AddEHandlerFunc(func(e gwu.Event) {
		getPoints = true
	}, gwu.ETYPE_CLICK)
	
	
	donePointsButton.AddEHandlerFunc(func(e gwu.Event) {
		getPoints = false
	}, gwu.ETYPE_CLICK)
	
	clearPointsButton.AddEHandlerFunc(func(e gwu.Event) {
		getPoints = false
		points = make([]*Position, 0)
		pointsBox.SetText("")
		pointsBox.SetRows(0)
		e.MarkDirty(pointsBox)
	}, gwu.ETYPE_CLICK)
	
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
	
	plotImg.AddEHandlerFunc(func(e gwu.Event) {
		x, y := e.Mouse()
		y = maxY - y
		if getOrigin == true {
			origin = &Position{x, y}
			axOriginLabel.SetText("Origin at: " + origin.Str())
			getOrigin = false
			e.MarkDirty(axOriginLabel)
		} else if getPoints == true {
			pos = &Position{x-origin.X, y-origin.Y}
			points = append(points, pos)
			pointsBox.SetText(pointsBox.Text() + pos.Str() + "\n")
			pointsBox.SetRows(len(points))
			e.MarkDirty(pointsBox)
		}
	}, gwu.ETYPE_CLICK)
	
	clearAllButton.AddEHandlerFunc(func(e gwu.Event) {
		inFileLoadBox.SetText("mapelli_progenitor_remnant.png")
		e.MarkDirty(inFileLoadBox)
		getOrigin = false
		axOriginLabel.SetText("Origin at: " + "--, --")
		e.MarkDirty(axOriginLabel)
		getPoints = false
		points = make([]*Position, 0)
		pointsBox.SetText("")
		pointsBox.SetRows(0)
		e.MarkDirty(pointsBox)
		plotImg.SetUrl("")
		e.MarkDirty(plotImg)
		clearAllButton.SetEnabled(false)
		e.MarkDirty(clearAllButton)
	}, gwu.ETYPE_CLICK)
	
	
	
	http.HandleFunc("/", ImgHandler3)
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

func ImgHandler (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func ImgHandler2 (w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, "image", time.Now(), f)
}

func ImgHandler3 (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/png")
	imgBytes.WriteTo(w)
}


package main

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type Route struct {
	Name  string   `json:"name"`
	Stops []string `json:"stops"`
}

func fetchRoutes(apiURL string) ([]Route, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var routes []Route
	if err := json.Unmarshal(body, &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

func getTimeEverySecond() <-chan string {
	// Create a channel to send the current time
	timeChannel := make(chan string)

	go func() {
		for {
			currentTime := time.Now().Format("15:04:05")
			timeChannel <- currentTime
			time.Sleep(1 * time.Second)
		}
	}()
	return timeChannel
}

func main() {
	// App entry point
	myApp := app.New()
	myWindow := myApp.NewWindow("Route Display")
	os.Setenv("FYNE_SCALE", "5")

	// Colors
	colorRed := color.RGBA{R: 250, G: 0, B: 0, A: 250}

	// Top container
	currentTimeText := canvas.NewText("Loading time...", colorRed)
	currentRouteText := canvas.NewText("Loading route...", colorRed)

	topHContainer := container.New(layout.NewHBoxLayout(), currentTimeText, layout.NewSpacer(), currentRouteText)
	topContainer := container.New(layout.NewVBoxLayout(), topHContainer)

	// Middle container
	nextStopText := canvas.NewText("Loading next stop...", colorRed)
	middleContainer := container.New(layout.NewCenterLayout(), layout.NewSpacer(), nextStopText)

	// Main container
	mainContainer := container.New(layout.NewStackLayout(), topContainer, middleContainer)

	timeChannel := getTimeEverySecond()
	apiURL := "http://192.168.64.20:3000/getAllRoutes"

	go func() {
		for {
			// Set time
			currentTime := <-timeChannel
			currentTimeText.Text = currentTime
			currentTimeText.Refresh()

			// Handle fetching
			routes, err := fetchRoutes(apiURL)
			if err != nil {
				// Log error
				log.Println(err)
			} else {
				// Update route text
				currentRouteText.Text = routes[len(routes)-1].Name
				currentRouteText.Refresh()
			}
		}
	}()

	myWindow.SetFullScreen(true)
	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
}

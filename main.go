package main

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"net/http"
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
	myApp := app.New()
	myWindow := myApp.NewWindow("Route Display")

	currentTimeLabel := canvas.NewText("Loading time...", color.White)
	topContainer := container.New(layout.NewHBoxLayout(), currentTimeLabel)

	timeChannel := getTimeEverySecond()

	go func() {
		for {
			// Set time
			currentTime := <-timeChannel
			currentTimeLabel.Text = currentTime
			currentTimeLabel.Refresh()
		}
	}()

	myWindow.SetFullScreen(true)
	myWindow.SetContent(topContainer)
	myWindow.ShowAndRun()
}

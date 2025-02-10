package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Route Display")

	currentRoute := widget.NewLabel("Route: Loading...")
	nextStop := widget.NewLabel("Next: Loading...")
	container := container.NewVBox(currentRoute, nextStop)

	apiURL := "http://192.168.64.19:3000/getAllRoutes" // Replace with actual API URL

	go func() {
		for {
			routes, err := fetchRoutes(apiURL)
			if err != nil {
				fmt.Println("Error fetching routes:", err)
			} else {
				// If routes are fetched successfully, display the first route's info
				if len(routes) > 0 {
					currentRoute.SetText("Route: " + routes[0].Name)
					if len(routes[0].Stops) > 0 {
						nextStop.SetText("Next stop: " + routes[0].Stops[0])
					} else {
						nextStop.SetText("Next stop: No stops available")
					}
				} else {
					currentRoute.SetText("Route: No routes available")
					nextStop.SetText("Next stop: N/A")
				}
			}
			time.Sleep(5 * time.Second) // Refresh every 5 seconds
		}
	}()

	myWindow.SetFullScreen(true)
	myWindow.SetContent(container)
	myWindow.ShowAndRun()
}

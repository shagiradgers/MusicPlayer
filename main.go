package main

import "MusicPlayer/pkg/app"

func main() {
	a, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	a.Run()
}

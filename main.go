package main

import (
	"flag"
	"log"
	"spot/menus"
	"spot/setup"
	"spot/spotify"
)

func setupComplete(action string) bool {
	return spotify.IsDataSet(spotify.ClientId) &&
		spotify.IsDataSet(spotify.ClientSecret) &&
		(spotify.IsDataSet(spotify.Token) || action == "auth")
}

func main() {
	actionPtr := flag.String("action", "", "action to take")
	artistPtr := flag.String("artist", "", "artist to focus")
	albumPtr := flag.String("album", "", "album to focus")
	trackPtr := flag.String("track", "", "track to play")

	flag.Parse()

	args := flag.Args()

	if !setupComplete(*actionPtr) {
		err := setup.Creds(args)
		if err != nil {
			log.Fatal(err)
			return
		}

		return
	}

	if *actionPtr != "" {
		action := *actionPtr
		if action == "playtrack" {
			err := spotify.PlayTrack(*trackPtr)
			if err != nil {
				log.Fatal(err)
				return
			}
		} else if action == "auth" {
			err := setup.LaunchAuth()
			if err != nil {
				log.Fatal(err)
				return
			}
		} else if action == "server" {
			err := setup.Server()
			if err != nil {
				log.Fatal(err)
				return
			}
		} else {
			err := spotify.ApplyControl(action)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		return
	}

	if *albumPtr != "" {
		err := menus.AlbumDetailMenu(*albumPtr, args)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *artistPtr != "" {
		err := menus.ArtistDetailMenu(*artistPtr, args)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	spotifyStatus, err := spotify.GetStatus()
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(args) > 0 && (args[0] == "c") {
		err = menus.ControlsMenu(spotifyStatus)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else if len(args) > 0 && (args[0] == "auth") {
		err = menus.SetupAuthMenu()
		if err != nil {
			log.Fatal(err)
			return
		}
	} else if len(args) > 0 {
		err = menus.SearchMenu(args)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		err = menus.StatusMenu(spotifyStatus)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

// TODO: playlist behavior, playlist cache
// TODO: alternate actions

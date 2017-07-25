package main

import (
	"flag"
	"log"
	"spot/menus"
	"spot/setup"
	"spot/spotify"
)

func main() {
	actionPtr := flag.String("action", "", "action to take")
	artistPtr := flag.String("artist", "", "artist to focus")
	albumPtr := flag.String("album", "", "album to focus")
	trackPtr := flag.String("track", "", "track to play")

	flag.Parse()

	if !spotify.IsDataSet(spotify.ClientId) || !spotify.IsDataSet(spotify.ClientSecret) {
		args := flag.Args()
		err := setup.ClientCreds(args)
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
		err := menus.AlbumDetailMenu(*albumPtr)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *artistPtr != "" {
		err := menus.ArtistDetailMenu(*artistPtr)
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

	if len(flag.Args()) > 0 && (flag.Args()[0] == "c") {
		err = menus.ControlsMenu(spotifyStatus)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else if len(flag.Args()) > 0 && (flag.Args()[0] == "auth") {
		err = menus.SetupMenu()
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

// TODO: better search handling of albums, use URI? always look up URI?
// TODO: playlist behavior, playlist cache
// TODO: alternate actions

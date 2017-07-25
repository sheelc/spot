package menus

import (
	"fmt"
	client "github.com/zmb3/spotify"
	"spot/alfred"
	"spot/spotify"
)

func SetupMenu() error {
	items := alfred.AlfredItems{
		Items: []alfred.AlfredItem{
			alfred.AlfredItem{
				Title: "Authenticate to Spot",
				Icon: alfred.AlfredIcon{
					Path: "icons/configuration.png",
				},
				Valid:        newFalse(),
				Autocomplete: "--action auth",
			},
		},
	}

	return alfred.PrintMenu(items)
}

func ClientIdMenuInstruction() error {
	items := alfred.AlfredItems{
		Items: []alfred.AlfredItem{
			alfred.AlfredItem{
				Title: "Paste the Client Id above and press enter",
			},
		},
	}

	return alfred.PrintMenu(items)
}

func ClientIdMenuStepFinished() error {
	items := alfred.AlfredItems{
		Items: []alfred.AlfredItem{
			alfred.AlfredItem{
				Title: "The client id is set, press enter to close",
			},
		},
	}

	return alfred.PrintMenu(items)
}

func ClientSecretMenuInstruction() error {
	items := alfred.AlfredItems{
		Items: []alfred.AlfredItem{
			alfred.AlfredItem{
				Title: "Paste the Client Secret above and press enter",
			},
		},
	}

	return alfred.PrintMenu(items)
}

func ClientSecretMenuStepFinished() error {
	items := alfred.AlfredItems{
		Items: []alfred.AlfredItem{
			alfred.AlfredItem{
				Title: "Setup is now complete! Press enter to close",
			},
		},
	}

	return alfred.PrintMenu(items)
}
func StatusMenu(spotifyStatus spotify.Status) error {
	var playerStatusIcon string
	if spotifyStatus.PlayerState == "paused" {
		playerStatusIcon = "icons/playing.png"
	} else {
		playerStatusIcon = "icons/paused.png"
	}

	items := alfred.AlfredItems{
		Items: []alfred.AlfredItem{
			alfred.AlfredItem{
				Title:    spotifyStatus.SongTitle,
				Subtitle: fmt.Sprintf("%s by %s", spotifyStatus.Album, spotifyStatus.Artist),
				Arg:      "-action playpause",
				Icon: alfred.AlfredIcon{
					Path: playerStatusIcon,
				},
			},
			alfred.AlfredItem{
				Title:    spotifyStatus.Artist,
				Subtitle: "More from this artist...",
				Icon: alfred.AlfredIcon{
					Path: "icons/artist.png",
				},
				Valid:        newFalse(),
				Autocomplete: fmt.Sprintf("-artist=\"%s\"", spotifyStatus.Artist),
			},
			alfred.AlfredItem{
				Title:    spotifyStatus.Album,
				Subtitle: "More from this album...",
				Icon: alfred.AlfredIcon{
					Path: "icons/album.png",
				},
				Valid:        newFalse(),
				Autocomplete: fmt.Sprintf("-album=\"%s\"", spotifyStatus.Album),
			},
		},
	}

	return alfred.PrintMenu(items)
}

func ControlsMenu(spotifyStatus spotify.Status) error {
	var playerStatusIcon string
	if spotifyStatus.PlayerState == "paused" {
		playerStatusIcon = "icons/playing.png"
	} else {
		playerStatusIcon = "icons/paused.png"
	}

	items := alfred.AlfredItems{
		Items: []alfred.AlfredItem{
			alfred.AlfredItem{
				Title:    spotifyStatus.SongTitle,
				Subtitle: fmt.Sprintf("%s by %s", spotifyStatus.Album, spotifyStatus.Artist),
				Arg:      "--action playpause",
				Icon: alfred.AlfredIcon{
					Path: playerStatusIcon,
				},
			},
			alfred.AlfredItem{
				Title: "Next Track",
				Arg:   "--action next",
				Icon: alfred.AlfredIcon{
					Path: "icons/next.png",
				},
			},
			alfred.AlfredItem{
				Title: "Previous Track",
				Arg:   "--action previous",
				Icon: alfred.AlfredIcon{
					Path: "icons/previous.png",
				},
			},
		},
	}

	return alfred.PrintMenu(items)
}

func AlbumDetailMenu(album string) error {
	spotifyClient, err := spotify.NewClient()
	if err != nil {
		return err
	}

	searchResults, err := spotifyClient.Search(fmt.Sprintf("album:%s", album), client.SearchTypeTrack)
	if err != nil {
		return err
	}

	tracks := searchResults.Tracks.Tracks
	items := make([]alfred.AlfredItem, 0, len(tracks))

	for _, track := range tracks {
		items = append(items, alfred.AlfredItem{
			Uid:   string(track.URI),
			Title: track.Name,
			Icon: alfred.AlfredIcon{
				Path: "icons/track.png",
			},
			Arg: fmt.Sprintf("--action playtrack --track %s", track.URI),
		})

	}

	alfredItems := alfred.AlfredItems{
		Items: items,
	}

	return alfred.PrintMenu(alfredItems)
}

func ArtistDetailMenu(artist string) error {
	spotifyClient, err := spotify.NewClient()
	if err != nil {
		return err
	}

	searchResults, err := spotifyClient.Search(fmt.Sprintf("artist:%s", artist), client.SearchTypeAlbum)
	if err != nil {
		return err
	}

	albums := searchResults.Albums.Albums
	items := make([]alfred.AlfredItem, 0, len(albums))

	for _, album := range albums {
		items = append(items, alfred.AlfredItem{
			Uid:          string(album.URI),
			Title:        album.Name,
			Valid:        newFalse(),
			Autocomplete: fmt.Sprintf("-artist=\"%s\" -album=\"%s\"", artist, album.Name),
			Icon: alfred.AlfredIcon{
				Path: "icons/album.png",
			},
		})

	}

	alfredItems := alfred.AlfredItems{
		Items: items,
	}

	return alfred.PrintMenu(alfredItems)
}

func newFalse() *bool {
	b := false
	return &b
}
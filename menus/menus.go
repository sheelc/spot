package menus

import (
	"fmt"
	client "github.com/zmb3/spotify"
	"spot/alfred"
	"spot/spotify"
	"strings"
)

func SetupAuthMenu() error {
	items := alfred.Items{
		Items: []alfred.Item{
			alfred.Item{
				Title: "Authenticate to Spot",
				Icon: alfred.Icon{
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
	items := alfred.Items{
		Items: []alfred.Item{
			alfred.Item{
				Title: "Paste the Client Id above and press enter",
			},
		},
	}

	return alfred.PrintMenu(items)
}

func ClientIdMenuStepFinished() error {
	items := alfred.Items{
		Items: []alfred.Item{
			alfred.Item{
				Title: "The client id is set, press enter to close",
			},
		},
	}

	return alfred.PrintMenu(items)
}

func ClientSecretMenuInstruction() error {
	items := alfred.Items{
		Items: []alfred.Item{
			alfred.Item{
				Title: "Paste the Client Secret above and press enter",
			},
		},
	}

	return alfred.PrintMenu(items)
}

func ClientSecretMenuStepFinished() error {
	items := alfred.Items{
		Items: []alfred.Item{
			alfred.Item{
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

	items := alfred.Items{
		Items: []alfred.Item{
			alfred.Item{
				Title:    spotifyStatus.SongTitle,
				Subtitle: fmt.Sprintf("%s by %s", spotifyStatus.Album, spotifyStatus.Artist),
				Arg:      "-action playpause",
				Icon: alfred.Icon{
					Path: playerStatusIcon,
				},
			},
			alfred.Item{
				Title:    spotifyStatus.Artist,
				Subtitle: "More from this artist...",
				Icon: alfred.Icon{
					Path: "icons/artist.png",
				},
				Valid:        newFalse(),
				Autocomplete: fmt.Sprintf("-artist=\"%s\" ", spotifyStatus.Artist),
			},
			alfred.Item{
				Title:    spotifyStatus.Album,
				Subtitle: "More from this album...",
				Icon: alfred.Icon{
					Path: "icons/album.png",
				},
				Valid:        newFalse(),
				Autocomplete: fmt.Sprintf("-album=\"%s\" ", spotifyStatus.Album),
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

	items := alfred.Items{
		Items: []alfred.Item{
			alfred.Item{
				Title: "Next Track",
				Arg:   "--action next",
				Icon: alfred.Icon{
					Path: "icons/next.png",
				},
			},
			alfred.Item{
				Title: "Previous Track",
				Arg:   "--action previous",
				Icon: alfred.Icon{
					Path: "icons/previous.png",
				},
			},
			alfred.Item{
				Title:    spotifyStatus.SongTitle,
				Subtitle: fmt.Sprintf("%s by %s", spotifyStatus.Album, spotifyStatus.Artist),
				Arg:      "--action playpause",
				Icon: alfred.Icon{
					Path: playerStatusIcon,
				},
			},
		},
	}

	return alfred.PrintMenu(items)
}

func popularityToMeter(popularity int) string {
	meter := ""
	for i := 0; i <= 100; i = i + 10 {
		if i < popularity {
			meter = meter + "■"
		} else {
			meter = meter + "□"
		}
	}
	return meter
}

func albumItem(album client.SimpleAlbum) alfred.Item {
	return alfred.Item{
		Uid:          string(album.URI),
		Title:        album.Name,
		Valid:        newFalse(),
		Autocomplete: fmt.Sprintf("-album=\"%s\" ", album.Name),
		Icon: alfred.Icon{
			Path: "icons/album.png",
		},
		Mods: alfred.Mods{
			Ctrl: alfred.Mod{
				Arg:      fmt.Sprintf("--action revealinspotify --context %s", album.URI),
				Subtitle: "Reveal in Spotify",
			},
		},
	}
}

func albumItemWithArtist(album client.SimpleAlbum, artist string) alfred.Item {
	item := albumItem(album)
	item.Autocomplete = fmt.Sprintf("-artist=\"%s\" -album=\"%s\" ", artist, album.Name)
	return item
}

func trackItem(track client.FullTrack) alfred.Item {
	return alfred.Item{
		Uid:      string(track.URI),
		Title:    track.Name,
		Subtitle: popularityToMeter(track.Popularity),
		Icon: alfred.Icon{
			Path: "icons/track.png",
		},
		Arg: fmt.Sprintf("--action playtrack --track %s --context %s", track.URI, track.Album.URI),
		Mods: alfred.Mods{
			Ctrl: alfred.Mod{
				Arg:      fmt.Sprintf("--action revealinspotify --context %s", track.URI),
				Subtitle: "Reveal in Spotify",
			},
		},
	}
}

func artistItem(artist client.FullArtist) alfred.Item {
	return alfred.Item{
		Uid:      string(artist.URI),
		Title:    artist.Name,
		Subtitle: popularityToMeter(artist.Popularity),
		Icon: alfred.Icon{
			Path: "icons/artist.png",
		},
		Valid:        newFalse(),
		Autocomplete: fmt.Sprintf("-artist=\"%s\" ", artist.Name),
		Mods: alfred.Mods{
			Ctrl: alfred.Mod{
				Arg:      fmt.Sprintf("--action revealinspotify --context %s", artist.URI),
				Subtitle: "Reveal in Spotify",
			},
		},
	}
}

func AlbumDetailMenu(album string, args []string) error {
	spotifyClient, err := spotify.NewClient()
	if err != nil {
		return err
	}

	searchString := strings.Join(args, " ")
	searchResults, err := spotifyClient.Search(fmt.Sprintf("album:\"%s\" %s", album, searchString), client.SearchTypeTrack, 12)
	if err != nil {
		return err
	}

	tracks := searchResults.Tracks.Tracks
	items := make([]alfred.Item, 0, len(tracks))

	for _, track := range tracks {
		items = append(items, trackItem(track))

	}

	alfredItems := alfred.Items{
		Items: items,
	}

	return alfred.PrintMenu(alfredItems)
}

func ArtistDetailMenu(artist string, args []string) error {
	spotifyClient, err := spotify.NewClient()
	if err != nil {
		return err
	}

	searchString := strings.Join(args, " ")
	searchResults, err := spotifyClient.Search(fmt.Sprintf("artist:\"%s\" %s", artist, searchString), client.SearchTypeAlbum|client.SearchTypeTrack, 3)
	if err != nil {
		return err
	}

	albums := searchResults.Albums.Albums
	tracks := searchResults.Tracks.Tracks
	items := make([]alfred.Item, 0, len(albums)+len(tracks))

	for _, track := range tracks {
		items = append(items, trackItem(track))

	}

	for _, album := range albums {
		items = append(items, albumItemWithArtist(album, artist))

	}

	alfredItems := alfred.Items{
		Items: items,
	}

	return alfred.PrintMenu(alfredItems)
}

func SearchMenu(args []string) error {
	spotifyClient, err := spotify.NewClient()
	if err != nil {
		return err
	}

	searchString := strings.Join(args, " ")
	searchResults, err := spotifyClient.Search(searchString, client.SearchTypeArtist|client.SearchTypeAlbum|client.SearchTypeTrack, 2)
	if err != nil {
		return err
	}

	albums := searchResults.Albums.Albums
	tracks := searchResults.Tracks.Tracks
	artists := searchResults.Artists.Artists
	items := make([]alfred.Item, 0, len(albums)+len(tracks)+len(artists))

	for _, track := range tracks {
		items = append(items, trackItem(track))

	}

	for _, artist := range artists {
		items = append(items, artistItem(artist))

	}

	for _, album := range albums {
		items = append(items, albumItem(album))

	}

	alfredItems := alfred.Items{
		Items: items,
	}

	return alfred.PrintMenu(alfredItems)
}

func newFalse() *bool {
	b := false
	return &b
}

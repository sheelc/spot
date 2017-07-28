package menus

import (
	"fmt"
	client "github.com/zmb3/spotify"
	"spot/alfred"
	"spot/spotify"
	"strings"
)

const (
	albumContext  = iota
	artistContext = iota
)

func SetupAuthMenu() error {
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
				Autocomplete: fmt.Sprintf("-artist=\"%s\" ", spotifyStatus.Artist),
			},
			alfred.AlfredItem{
				Title:    spotifyStatus.Album,
				Subtitle: "More from this album...",
				Icon: alfred.AlfredIcon{
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

	items := alfred.AlfredItems{
		Items: []alfred.AlfredItem{
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
			alfred.AlfredItem{
				Title:    spotifyStatus.SongTitle,
				Subtitle: fmt.Sprintf("%s by %s", spotifyStatus.Album, spotifyStatus.Artist),
				Arg:      "--action playpause",
				Icon: alfred.AlfredIcon{
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

func albumItem(album client.SimpleAlbum) alfred.AlfredItem {
	return alfred.AlfredItem{
		Uid:          string(album.URI),
		Title:        album.Name,
		Valid:        newFalse(),
		Autocomplete: fmt.Sprintf("-album=\"%s\" ", album.Name),
		Icon: alfred.AlfredIcon{
			Path: "icons/album.png",
		},
	}
}

func albumItemWithArtist(album client.SimpleAlbum, artist string) alfred.AlfredItem {
	item := albumItem(album)
	item.Autocomplete = fmt.Sprintf("-artist=\"%s\" -album=\"%s\" ", artist, album.Name)
	return item
}

func trackItem(track client.FullTrack, contextType int) alfred.AlfredItem {
	var arg string
	if contextType == albumContext {
		arg = fmt.Sprintf("--action playtrack --track %s --context %s", track.URI, track.Album.URI)
	} else if contextType == artistContext {
		arg = fmt.Sprintf("--action playtrack --track %s --context %s", track.URI, track.Artists[0].URI)
	} else {
		arg = fmt.Sprintf("--action playtrack --track %s", track.URI)
	}

	return alfred.AlfredItem{
		Uid:      string(track.URI),
		Title:    track.Name,
		Subtitle: popularityToMeter(track.Popularity),
		Icon: alfred.AlfredIcon{
			Path: "icons/track.png",
		},
		Arg: arg,
	}
}

func artistItem(artist client.FullArtist) alfred.AlfredItem {
	return alfred.AlfredItem{
		Uid:      string(artist.URI),
		Title:    artist.Name,
		Subtitle: popularityToMeter(artist.Popularity),
		Icon: alfred.AlfredIcon{
			Path: "icons/artist.png",
		},
		Valid:        newFalse(),
		Autocomplete: fmt.Sprintf("-artist=\"%s\" ", artist.Name),
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
	items := make([]alfred.AlfredItem, 0, len(tracks))

	for _, track := range tracks {
		items = append(items, trackItem(track, albumContext))

	}

	alfredItems := alfred.AlfredItems{
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
	items := make([]alfred.AlfredItem, 0, len(albums)+len(tracks))

	for _, track := range tracks {
		items = append(items, trackItem(track, artistContext))

	}

	for _, album := range albums {
		items = append(items, albumItemWithArtist(album, artist))

	}

	alfredItems := alfred.AlfredItems{
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
	items := make([]alfred.AlfredItem, 0, len(albums)+len(tracks)+len(artists))

	for _, track := range tracks {
		items = append(items, trackItem(track, artistContext))

	}

	for _, artist := range artists {
		items = append(items, artistItem(artist))

	}

	for _, album := range albums {
		items = append(items, albumItem(album))

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

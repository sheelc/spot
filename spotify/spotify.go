package spotify

import (
	"bytes"
	"fmt"
	client "github.com/zmb3/spotify"
	"log"
	"os/exec"
	"strings"
)

type Client struct {
	transport *client.Client
}

type Status struct {
	SongTitle   string
	Album       string
	Artist      string
	TrackUrl    string
	PlayerState string
}

func NewClient() (*Client, error) {
	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	auth, err := NewAuth()
	if err != nil {
		return nil, err
	}

	transport := auth.NewClient(token)
	return &Client{
		transport: &transport,
	}, nil
}

func GetStatus() (Status, error) {
	cmd := exec.Command("osascript", "-e", `tell application "Spotify" to return name of current track & "Ψ" & album of current track & "Ψ" & artist of current track & "Ψ" & spotify url of current track & "Ψ" & player state`)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return Status{}, err
	}

	str := out.String()
	spotifyStr := str[:len(str)-1]
	spotifyData := strings.Split(spotifyStr, "Ψ")
	return Status{
		SongTitle:   spotifyData[0],
		Album:       spotifyData[1],
		Artist:      spotifyData[2],
		TrackUrl:    spotifyData[3],
		PlayerState: spotifyData[4],
	}, nil
}

var controlsDict = map[string]string{
	"playpause": "playpause",
	"next":      "play next track",
	"previous":  "play previous track",
}

func ApplyControl(control string) error {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Spotify\" to %s", controlsDict[control]))
	return cmd.Run()
}

func PlayTrack(trackUri string, contextUri string) error {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Spotify\" to play track \"%s\" in context \"%s\"", trackUri, contextUri))
	return cmd.Run()
}

func Reveal(contextUri string) error {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Spotify\" to open location \"%s\" & (activate)", contextUri))
	return cmd.Run()
}

func (c *Client) Search(searchStr string, st client.SearchType, limit int) (*client.SearchResult, error) {
	opts := client.Options{
		Limit: &limit,
	}
	return c.transport.SearchOpt(searchStr, st, &opts)
}

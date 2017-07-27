package setup

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"spot/menus"
	"spot/spotify"
	"time"
)

func LaunchAuth() error {
	auth, err := spotify.NewAuth()
	if err != nil {
		return err
	}

	url := auth.AuthURL("state")

	daemon := exec.Command("./spot", "--action=server")
	daemon.Env = os.Environ()
	err = daemon.Start()
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	cmd := exec.Command("open", url)
	return cmd.Run()
}

func SaveClientId(clientid string) error {
	return spotify.Set(spotify.ClientId, []byte(clientid))
}

func SaveClientSecret(clientsecret string) error {
	return spotify.Set(spotify.ClientSecret, []byte(clientsecret))
}

func setupClientId(args []string) error {
	if len(args) > 0 {
		err := SaveClientId(args[0])
		if err != nil {
			return err
		}

		err = menus.ClientIdMenuStepFinished()
		if err != nil {
			return err
		}
		return nil
	}

	err := menus.ClientIdMenuInstruction()
	if err != nil {
		return err
	}

	return nil
}

func setupClientSecret(args []string) error {
	if len(args) > 0 {
		err := SaveClientSecret(args[0])
		if err != nil {
			return err
		}

		err = menus.ClientSecretMenuStepFinished()
		if err != nil {
			return err
		}
		return nil
	}

	err := menus.ClientSecretMenuInstruction()
	if err != nil {
		return err
	}

	return nil
}

func Creds(args []string) error {
	if !spotify.IsDataSet(spotify.ClientId) {
		return setupClientId(args)
	}

	if !spotify.IsDataSet(spotify.ClientSecret) {
		return setupClientSecret(args)
	}

	if !spotify.IsDataSet(spotify.Token) {
		return menus.SetupAuthMenu()
	}

	return nil
}

func Server() error {
	go func() {
		time.Sleep(5 * time.Minute)
		os.Exit(0)
	}()

	http.HandleFunc("/callback", callbackHandler)
	return http.ListenAndServe(":11075", nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	auth, err := spotify.NewAuth()
	if err != nil {
		http.Error(w, "Failed reading client config", http.StatusInternalServerError)
		return
	}

	token, err := auth.Token("state", r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	err = spotify.SetToken(token)
	if err != nil {
		http.Error(w, "Couldn't write token file", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "success! this window can be closed now")
}

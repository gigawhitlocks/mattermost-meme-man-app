package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
	"github.com/mattermost/mattermost-server/v5/model"
)

//go:embed icon.png
var iconData []byte

//go:embed manifest.json
var manifestData []byte

//go:embed bindings.json
var bindingsData []byte

const (
	host = "localhost"
	port = 8080
)

func main() {
	// Serve its own manifest as HTTP for convenience in dev. mode.
	http.HandleFunc("/manifest.json", writeJSON(manifestData))

	// Returns the Channel Header and Command bindings for the app.
	http.HandleFunc("/bindings", writeJSON(bindingsData))

	http.HandleFunc("/post-meme/submit", post)

	// Serves the icon for the app.
	http.HandleFunc("/static/icon.png", writeData("image/png", iconData))

	addr := fmt.Sprintf("%v:%v", host, port)
	fmt.Printf(`memes app listening at http://%s`, addr)
	http.ListenAndServe(addr, nil)
}

func post(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	json.NewDecoder(req.Body).Decode(&c)

	// v, ok := c.Values["message"]
	// if ok && v != nil {
	// 	message += fmt.Sprintf(" ...and %s!", v)
	// }
	if c.Context.ActingUserAccessToken == "" {
		json.NewEncoder(w).Encode(apps.CallResponse{
			Type:     apps.CallResponseTypeError,
			Markdown: "empty",
		})

	}

	mmclient.AsActingUser(c.Context).CreatePost(
		&model.Post{
			UserId:    c.Context.ActingUserID,
			ChannelId: c.Context.ChannelID,
			Message:   fmt.Sprintf("![](%s%s/static/icon.png)", c.Context.MattermostSiteURL, c.Context.AppPath),
		},
	)

	json.NewEncoder(w).Encode(apps.CallResponse{
		Type: apps.CallResponseTypeOK,
	})
}

func writeData(ct string, data []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", ct)
		w.Write(data)
	}
}

func writeJSON(data []byte) func(w http.ResponseWriter, r *http.Request) {
	return writeData("application/json", data)
}

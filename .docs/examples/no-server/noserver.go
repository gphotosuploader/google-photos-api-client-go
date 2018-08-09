package main

import (
	"log"

	"github.com/nmrshll/google-photos-api-client-go/gphotoslib"
	gphotosclient "github.com/nmrshll/google-photos-api-client-go/noserver-gphotosclient"
)

var (
	// example credentials. DON'T USE THESE.
	// now they are public they could be overused by someone else, get blocked by google, thus breaking your app
	// create your at https://console.cloud.google.com/apis/credentials . And keep them private
	apiAppCredentials = gphotoslib.APIAppCredentials{
		ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
		ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
	}
)

func main() {
	// ask the user to authenticate on google in the browser
	photosClient, err := gphotosclient.NewClient(
		gphotosclient.AuthenticateUser(
			gphotoslib.NewOAuthConfig(gphotoslib.APIAppCredentials{
				ClientID:     "________________",
				ClientSecret: "____________________"},
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = photosClient.UploadFile("/path/to/file")
	if err != nil {
		log.Fatal(err)
	}
}

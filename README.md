[![Go Report Card](https://goreportcard.com/badge/github.com/nmrshll/gphotos-uploader-api-client-go)](https://goreportcard.com/report/github.com/nmrshll/gphotos-uploader-api-client-go)

# Google photos API client (Go library)
The [official google photos client](google.golang.org/api/photoslibrary/v1) doesn't have upload functionality, this repo aims to complete it with upload functionality and improve ease of use for classic use cases.   

It contains three packages, [github.com/nmrshll/google-photos-api-client-go/gphotoslib](github.com/nmrshll/google-photos-api-client-go/gphotoslib) simply wraps the official library to offer the same functionality plus upload,
the two other packages try to make it easier to use in classic cases:    
- [noserver-gphotosclient](github.com/nmrshll/google-photos-api-client-go/noserver-gphotosclient) for CLI tool or desktop/mobile app (when your app is fully client-side and doesn't have its own server-side part): [examples]()    
- [server-gphotosclient](github.com/nmrshll/google-photos-api-client-go/server-gphotosclient) for classic web apps that have a server side: [examples]()

# Quick start

Download using `go get github.com/nmrshll/google-photos-api-client-go`

Then use this way:

- In a CLI tool or desktop/mobile app (when your app is fully client-side and doesn't have its own server-side part):
[embedmd]:# (./.docs/examples/no-server/noserver.go)
```go
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
			gphotoslib.NewOAuthConfig(apiAppCredentials),
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
```

- on the server side:
[embedmd]:# (./.docs/examples/server-side/serverside.go)
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello world !")
}
```

# More Examples

[server-side]():
- []()

[no server, all client side (desktop,cli,mobile,...)]():
- []()

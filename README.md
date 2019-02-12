[![Go Report Card](https://goreportcard.com/badge/github.com/nmrshll/gphotos-uploader-api-client-go)](https://goreportcard.com/report/github.com/nmrshll/gphotos-uploader-api-client-go)

# Google photos API client (Go library)
The [official google photos client](https://google.golang.org/api/photoslibrary/v1) doesn't have upload functionality, this repo aims to complete it with upload functionality and improve ease of use for classic use cases.   

It contains three packages, [github.com/nmrshll/google-photos-api-client-go/lib-gphotos](https://github.com/nmrshll/google-photos-api-client-go/tree/master/lib-gphotos) simply wraps the official library to offer the same functionality plus upload,
the two other packages try to make it easier to use in classic cases:    
- [noserver-gphotosclient](github.com/nmrshll/google-photos-api-client-go/noserver-gphotosclient) for CLI tool or desktop/mobile app (when your app is fully client-side and doesn't have its own server-side part): [examples]()    
- [server-gphotosclient](github.com/nmrshll/google-photos-api-client-go/server-gphotosclient) for classic web apps that have a server side: [examples]()

# Quick start

Download using `go get github.com/nmrshll/google-photos-api-client-go`

Then use this way:

#### In a CLI tool or desktop/mobile app (when your app is fully client-side and doesn't have its own server-side part):    

[embedmd]:# (./.docs/examples/no-server/noserver.go /func main/ $)
```go
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
```

#### on the server side:    

[embedmd]:# (./.docs/examples/server-side/serverside.go /func main/ $)
```go
func main() {
	fmt.Println("Hello world !")
}
```

# More Examples

[server-side]():
- []()

[no server, all client side (desktop,cli,mobile,...)]():
- []()

package gphotos

import gphotos "github.com/nmrshll/google-photos-api-client-go/lib-gphotos"

type Client struct {
	gphotos.Client
}

var (
	NewClient = gphotos.NewClient
)

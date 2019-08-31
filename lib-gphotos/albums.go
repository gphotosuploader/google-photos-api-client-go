package gphotos

import (
	"errors"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

const (
	defaultPageSize = 50
)

var (
	ErrAlbumNotFound = errors.New("specified album was not found")
)

// albumByNameWithPageToken checks for the specified album recursively. Google Photos returns
// albums list in a paginated way, so we need to go through pages in order to check if the album
// exists.
//
// An error (gphotos.ErrAlbumNotFound) is returned if the album doesn't exists.
func (c *Client) albumByName(name, pageToken string) (album *photoslibrary.Album, err error) {
	albumListCall := c.Albums.List().PageSize(defaultPageSize).PageToken(pageToken)
	response, err := albumListCall.Do()
	if err != nil {
		return nil, err
	}
	for _, album := range response.Albums {
		if album.Title == name {
			// Album was found, so return the object.
			return album, nil
		}
	}

	if response.NextPageToken != "" {
		// There are more pages to check, go for next page.
		return c.albumByName(name, response.NextPageToken)
	}

	// The album doesn't exists.
	return nil, ErrAlbumNotFound
}

// AlbumByName returns the album which match with the specified name.
//
// NOTE: We are maintaining backwards compatibility, but `found` should be DEPRECATED and
// returning an error (gphotos.ErrAlbumNotFound) instead of it. (TODO)
func (c *Client) AlbumByName(name string) (album *photoslibrary.Album, found bool, err error) {
	a, err := c.albumByName(name, "")
	if err != nil {
		if err == ErrAlbumNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return a, true, nil
}

// GetOrCreateAlbumByName returns an Album with the specified album name.
// If the album doesn't exists it will try to create it.
func (c *Client) GetOrCreateAlbumByName(name string) (*photoslibrary.Album, error) {
	album, found, err := c.AlbumByName(name)
	if err != nil {
		return nil, err
	}

	if found {
		return album, nil
	}

	return c.Albums.Create(&photoslibrary.CreateAlbumRequest{
		Album: &photoslibrary.Album{Title: name},
	}).Do()
}

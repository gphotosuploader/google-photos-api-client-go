package gphotos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
)

const (
	maxPageSize = 50

	// Store media (albumGallery) for performance reasons.
	// See https://developers.google.com/photos/library/guides/best-practices#caching
	albumCacheTTL = 60 * time.Minute
)

var (
	ErrAlbumNotFound = errors.New("album was not found")
)

// ListAlbums returns all the albums in Google Photos service.
func (c *Client) ListAlbums(ctx context.Context) ([]*photoslibrary.Album, error) {
	var results []*photoslibrary.Album
	err := c.ListAlbumsWithCallback(ctx, func(albums []*photoslibrary.Album, stop func()) {
		results = append(results, albums...)
	})

	return results, err
}

// ListAlbumsFunc is called for each response of 50 albumGallery.
// If this calls stop, ListAlbums stops the loop.
type ListAlbumsFunc func(albums []*photoslibrary.Album, stop func())

// ListAlbumsWithCallback iterate on each Google Photos album executing the provided callback.
func (c *Client) ListAlbumsWithCallback(ctx context.Context, callback ListAlbumsFunc) error {
	var pageToken string
	for {
		res, err := c.service.ListAlbums(ctx, maxPageSize, pageToken)
		if err != nil {
			return fmt.Errorf("error listing albums: err=%w", err)
		}

		// cache albums.
		for _, album := range res.Albums {
			_ = c.cache.PutAlbum(ctx, *album, albumCacheTTL)
		}

		var stop bool
		callback(res.Albums, func() { stop = true })
		if stop {
			return nil
		}
		if res.NextPageToken == "" {
			return nil
		}
		pageToken = res.NextPageToken
	}
}

// CreateAlbum creates an Album in Google Photos library and returns the created object.
// If the Album was already on the library, it will return the Album.
// It's ensuring that Albums are unique (for this app) on Google Photos.
func (c *Client) CreateAlbum(ctx context.Context, title string) (*photoslibrary.Album, error) {
	c.m.Lock()
	defer c.m.Unlock()
	album, err := c.FindAlbum(ctx, title)
	if !errors.Is(err, ErrAlbumNotFound) {
		// Album was found or there was an error with the cache.
		return album, err
	}

	album, err = c.service.CreateAlbum(ctx, &photoslibrary.CreateAlbumRequest{
		Album: &photoslibrary.Album{Title: title},
	})
	if err != nil {
		return &photoslibrary.Album{}, fmt.Errorf("could not create an album. err: %w", err)
	}

	// Cache the created album.
	_ = c.cache.PutAlbum(ctx, *album, albumCacheTTL)

	return album, nil
}

// FindAlbum search the Album with the specified title in Google Photos library and returns it.
// If the Album is not found, it will return ErrAlbumNotFound.
func (c *Client) FindAlbum(ctx context.Context, title string) (*photoslibrary.Album, error) {
	matched, err := c.cache.GetAlbum(ctx, title)
	if !errors.Is(err, cache.ErrCacheMiss) {
		// Album was found or there was an error with the cache.
		return &matched, err
	}

	if err := c.ListAlbumsWithCallback(ctx, func(albums []*photoslibrary.Album, stop func()) {
		for _, album := range albums {
			if album.Title == title {
				stop()
				matched = *album
				return
			}
		}
	}); err != nil {
		return &photoslibrary.Album{}, err
	}

	if len(matched.Title) == 0 {
		return &photoslibrary.Album{}, ErrAlbumNotFound
	}

	return &matched, nil
}

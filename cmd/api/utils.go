package main

import (
	"context"
	"errors"

	"github.com/jcorry/morellis/pkg/models"
	"googlemaps.github.io/maps"
)

func (app *application) geocodeStore(s *models.Store) error {
	c, err := maps.NewClient(maps.WithAPIKey(app.mapsApiKey))
	if err != nil {
		return err
	}
	r := &maps.GeocodingRequest{
		Address:  s.AddressString(),
		Language: "en",
		Region:   "us",
	}
	resp, err := c.Geocode(context.Background(), r)
	if err != nil {
		return err
	}

	if len(resp) < 1 {
		return errors.New("No geocoding response")
	}

	location := resp[0]
	s.Lat = location.Geometry.Location.Lat
	s.Lng = location.Geometry.Location.Lng

	return nil
}

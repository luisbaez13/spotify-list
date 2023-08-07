package main

type ArtistRespItems struct {
	Items []ArtistObj `json:"items"`
}

type ArtistObj struct {
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Genres []string `json:"genres"`
	Href   string   `json:"href"`
	ID     string   `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
	Type       string `json:"type"`
	URI        string `json:"uri"`
}

type TrackRespItems struct {
	Items []TrackObj `json:"items"`
}

type TrackObj struct {
	Album struct {
		AlbumType        string      `json:"album_type"`
		Artists          []ArtistObj `json:"artists"`
		AvailableMarkets []string    `json:"available_markets"`
		ExternalUrls     struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name                 string `json:"name"`
		ReleaseDate          string `json:"release_date"`
		ReleaseDatePrecision string `json:"release_date_precision"`
		TotalTracks          int    `json:"total_tracks"`
		Type                 string `json:"type"`
		URI                  string `json:"uri"`
	} `json:"album"`
	Artists []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"artists"`
	DurationMs       int      `json:"duration_ms"`
	Explicit         bool     `json:"explicit"`
	ExternalIds      struct {
		Isrc string `json:"isrc"`
	} `json:"external_ids"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href        string `json:"href"`
	ID          string `json:"id"`
	IsLocal     bool   `json:"is_local"`
	Name        string `json:"name"`
	Popularity  int    `json:"popularity"`
	PreviewURL  string `json:"preview_url"`
	TrackNumber int    `json:"track_number"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
}

type UserProfile struct {
	DisplayName  string `json:"display_name"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Type      string `json:"type"`
	URI       string `json:"uri"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Country         string `json:"country"`
	Product         string `json:"product"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	Email string `json:"email"`
}

//only save what is needed from response
type TokenResp struct {
	AccessToken string `json:"access_token"`
}

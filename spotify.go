package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func getTopArtists(spotifyUserToken, time_range, limit string) *ArtistRespItems {
	client := &http.Client{Timeout: time.Second * 10}

	q := "https://api.spotify.com/v1/me/top/artists?"
	q += "time_range=" + time_range + "&"
	q += "limit=" + limit

	fmt.Printf("TopArtistsQ: %s\n", q)

	req, err := http.NewRequest("GET", q, nil)
	if err != nil {
		fmt.Printf("Err creating get request:\n\t%v\n", err)
	}
	bearerHeader := "Bearer " + spotifyUserToken
	req.Header.Set("Authorization", bearerHeader)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Err doing GET request:\n\t%v\n", err)
	}
	fmt.Printf("TOP SONGS API RESPONSE: %v\n", resp.Status)

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("RESP Body: %s\n", b)
		}
		return &ArtistRespItems{
			Items: []ArtistObj{
				{Name: "Artist Not Found"},
			},
		}

	}

	userTop := ArtistRespItems{}
	json.NewDecoder(resp.Body).Decode(&userTop)

	return &userTop
}

func getTopTracks(spotifyUserToken, time_range, limit string) *TrackRespItems {
	client := &http.Client{Timeout: time.Second * 10}

	q := "https://api.spotify.com/v1/me/top/tracks?"
	q += "time_range=" + time_range + "&"
	q += "limit=" + limit

	fmt.Printf("TopTracksQ: %s\n", q)

	req, err := http.NewRequest("GET", q, nil)
	if err != nil {
		fmt.Printf("Err creating get request:\n\t%v\n", err)
	}
	bearerHeader := "Bearer " + spotifyUserToken
	req.Header.Set("Authorization", bearerHeader)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Err doing GET request:\n\t%v\n", err)
	}
	fmt.Printf("TOP TRACKS API RESPONSE: %v\n", resp.Status)

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("RESP Body: %s\n", b)
		}
		//if it fails return an empty item so things can be printed... not sure if this is the best idea
		return &TrackRespItems{
			Items: []TrackObj{
				{Name: "Songs Not Found"},
			},
		}
	}

	userTop := TrackRespItems{}
	json.NewDecoder(resp.Body).Decode(&userTop)

	return &userTop
}

func getUserProfile(spotifyUserToken string) *UserProfile {
	client := &http.Client{Timeout: time.Second * 10}

	q := "https://api.spotify.com/v1/me"
	req, err := http.NewRequest("GET", q, nil)
	if err != nil {
		fmt.Printf("Err creating get request:\n\t%v\n", err)
	}
	bearerHeader := "Bearer " + spotifyUserToken
	req.Header.Set("Authorization", bearerHeader)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Err doing GET request:\n\t%v\n", err)
	}
	userInfo := UserProfile{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	return &userInfo
}

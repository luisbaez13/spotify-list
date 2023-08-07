# Spotify Top Songs Visuals

A website to make a simple server written in Go, it takes your top Spotify songs and create a png of that list (similar idea to [recieptify](https://receiptify.herokuapp.com/)).

# Running it locally
1) Create a Spotify app [here](https://developer.spotify.com/dashboard)
1) Set the following environment variables using your client id and secret from Spotify `PORT, SITE_DOMAIN, SPOT_CLIENT_ID,` and `SPOT_CLIENT_SECRET`
1) In the settings page on the [Spotify developers page](https://developer.spotify.com/dashboard) add the callback link for your domain (for localhost it would be http://localhost:PORT/callback)
1) Run ```go run main.go``` and you should be able to login at `http://localhost:PORT`


## Demo Video

https://github.com/luisbaez13/spotify-stats/assets/95329512/53b39b1c-fc91-4dc1-b493-098492a368ba


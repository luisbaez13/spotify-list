package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var currentDomain = os.Getenv("SITE_DOMAIN")

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("resources/*.html")
	router.Static("/static", "resources/static")

	router.GET("/", landingPage)
	router.GET("/login", spotifyAuth)
	router.GET("/callback", callback)
	router.GET("/settings", settings)
	router.POST("/custom", custom)
	router.GET("/custom", custom)        // using GET on custom returns default settings
	router.GET("/json-data", dataAsJson) // return all data for a timeframe as json, using limit=50
	router.GET("/make-png", makePNG)     // returns image
	router.Run(":" + os.Getenv("PORT"))
}

func landingPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"loginlink": currentDomain + "/login",
	})
}

func spotifyAuth(c *gin.Context) {
	// Get this app approved by the user with their spotify acct and use the code that is returned
	// once the browser gets redirected and approved, they are redirected to {{domain}}/callback
	q := "https://accounts.spotify.com/authorize?"
	q += "client_id=" + os.Getenv("SPOT_CLIENT_ID") + "&"
	q += "client_secret=" + os.Getenv("SPOT_CLIENT_SECRET") + "&"
	q += "redirect_uri=" + url.QueryEscape(currentDomain+"/callback") + "&"
	q += "response_type=code&"
	q += "scope=" + url.QueryEscape("user-top-read user-library-read")

	c.Redirect(http.StatusFound, q)
}

func callback(c *gin.Context) {
	code := c.Query("code")
	client := &http.Client{Timeout: time.Second * 10}

	bodyStr := "grant_type=authorization_code&"
	bodyStr += "code=" + code + "&"
	bodyStr += "redirect_uri=" + url.QueryEscape(currentDomain+"/callback")

	//use code to get token
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(bodyStr)))
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}

	mySpotifyApp := os.Getenv("SPOT_CLIENT_ID") + ":" + os.Getenv("SPOT_CLIENT_SECRET")
	encodedId := b64.StdEncoding.EncodeToString([]byte(mySpotifyApp))
	req.Header.Set("Authorization", "Basic "+encodedId)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		fmt.Printf("Stat:%v\n", resp.Status)
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("%s\n", bodyString)
	}

	//Process json and get response
	respInfo := TokenResp{}
	json.NewDecoder(resp.Body).Decode(&respInfo)
	userToken := respInfo.AccessToken

	//save token as cookie
	c.SetCookie("spotify_token", userToken, 60*60*24, "/", "", true, true)

	//redirect to get user to not refresh on this route
	c.Redirect(http.StatusFound, "/custom")
}

func settings(c *gin.Context) {
	_, err := c.Cookie("spotify_token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
	}
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"api_search_link": currentDomain + "/custom",
	})
}

func custom(c *gin.Context) {
	_, err := c.Cookie("spotify_token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
	}

	kind := c.PostForm("type-select")
	time_range := c.PostForm("time-select")
	limit := c.PostForm("num-select")

	if kind == "" {
		kind = "tracks"
	}
	if time_range == "" {
		time_range = "short_term"
	}
	if limit == "" {
		limit = "15"
	}
	fmt.Printf("custom():: Formdata: %v, %v, %v\n", kind, time_range, limit)

	c.HTML(http.StatusOK, "findpic.html", gin.H{
		"settingslink": currentDomain + "/settings",
		"query":        currentDomain + "/make-png?kind=" + kind + "&time=" + time_range + "&limit=" + limit,
	})
}

func dataAsJson(c *gin.Context) {
	cookie, err := c.Cookie("spotify_token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
	}

	timeFrame := c.Query("time")

	limit := "50"

	trackRespData := getTopTracks(cookie, timeFrame, limit)
	artistRespData := getTopArtists(cookie, timeFrame, limit)
	userInfo := getUserProfile(cookie)

	type combined struct {
		TrackItems  []TrackObj  `json:"trackitems"`
		ArtistItems []ArtistObj `json:"artistitems"`
		UserName    string      `json:"username"`
		TimeFrame   string      `json:"timeframe"`
	}

	UserTracksAndArtist := &combined{
		TrackItems:  trackRespData.Items,
		ArtistItems: artistRespData.Items,
		UserName:    userInfo.DisplayName,
		TimeFrame:   limit,
	}

	c.JSON(200, UserTracksAndArtist)
}

func makePNG(c *gin.Context) {
	cookie, err := c.Cookie("spotify_token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
	}

	kind := c.Query("kind")
	timeFrame := c.Query("time")
	limit := c.Query("limit")

	if kind == "" {
		kind = "tracks"
	}
	if timeFrame == "" {
		timeFrame = "short_term"
	}
	if limit == "" {
		limit = "15"
	}

	trackRespData := &TrackRespItems{}
	ArtistRespData := &ArtistRespItems{}

	if kind == "tracks" {
		trackRespData = getTopTracks(cookie, timeFrame, limit)
	} else if kind == "artists" {
		ArtistRespData = getTopArtists(cookie, timeFrame, limit)
	}
	userInfo := getUserProfile(cookie)

	im, err := gg.LoadImage("resources/static/white.png")
	if err != nil {
		log.Fatal(err)
	}
	pt := im.Bounds()

	// read font file
	fontBytes, err := ioutil.ReadFile("resources/static/PermanentMarker-Regular.ttf")
	if err != nil {
		fmt.Printf("makePNG():: Error Reading font")
	}
	// parse font
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		fmt.Printf("makePNG():: Error Parsing font")
	}

	sz30Marker := truetype.NewFace(f, &truetype.Options{Size: 30})

	// Measures height of a wrapped string may be off by one line if there is a short word at the end of a line that is wrapped
	measureStringHeight := func(s string, f1 font.Face, lineSpacing float64, width float64) (multilineHeight, singlelineHeight float64, numlines int) {
		d := &font.Drawer{
			Face: f1,
		}
		a := d.MeasureString(s)
		numLines := (float64(a) / 64.0) / width
		roundedNumlines := math.Ceil(numLines)
		oneLineHeight := float64(f1.Metrics().Height) / 64.0
		return oneLineHeight * roundedNumlines, oneLineHeight, int(roundedNumlines)
	}

	// set up so the contents height can be measured before creating
	newPage := func(x, y int, title string) *gg.Context {
		dc := gg.NewContext(x, y)
		dc.SetRGB(0, 0, 0)
		dc.DrawImage(im, 0, 0)

		titleMargin := 60.0
		dc.SetFontFace(truetype.NewFace(f, &truetype.Options{Size: 60}))
		dc.DrawString(title, titleMargin, 90)
		titleLength, _ := dc.MeasureString(title)

		dc.SetFontFace(sz30Marker)
		dc.DrawString(userInfo.DisplayName, titleMargin+titleLength+20, 60)
		tagline := generateTagline(timeFrame)
		dc.DrawString(tagline, titleMargin+titleLength+40, 100)
		return dc
	}
	var dc *gg.Context
	pageWidth := 700.0 // based on the image used
	headerSize := 120  // How far down content starts after the title
	textHeight := 0.0
	footerMargin := 30.0

	if kind == "tracks" {
		topList := make([]string, len(trackRespData.Items))

		// take response and form into a list of strings with numbers, and measure the size of each
		for i, track := range trackRespData.Items {
			artistsAsString := ""
			if len(track.Artists) > 1 {
				artistsAsString += "by "
				for j, artist := range track.Artists {
					artistsAsString += artist.Name
					if j < len(track.Artists)-1 {
						artistsAsString += ", "
					}
				}
			} else {
				artistsAsString += " by " + track.Artists[0].Name
			}
			trackString := strconv.Itoa(i+1) + ". " + track.Name + " " + artistsAsString
			topList[i] = trackString
			mh, sh, _ := measureStringHeight(trackString, sz30Marker, 1.0, pageWidth-30)
			textHeight += mh + sh*0.5 //add the height of text + half of line space
		}

		dc = newPage(pt.Max.X-50, int(textHeight+footerMargin)+headerSize, "Top Songs:")
		verticalOffset := 0.0
		for i := 0; i < len(topList); i++ {
			mh, sh, _ := measureStringHeight(topList[i], sz30Marker, 1.0, pageWidth-30)

			dc.DrawStringWrapped(topList[i], 50, float64(headerSize)+(sh*(1.5)*float64(i))+verticalOffset, 0, 0, pageWidth, 1.0, gg.AlignLeft)

			if mh != sh {
				verticalOffset += (mh - sh)
			}
		}

	} else if kind == "artists" {
		topList := make([]string, len(ArtistRespData.Items))

		// take response and form into a list of strings with numbers, and measure the size of each
		for i, artist := range ArtistRespData.Items {
			artistNameString := strconv.Itoa(i+1) + ". " + artist.Name
			topList[i] = artistNameString
			mh, sh, _ := measureStringHeight(artistNameString, sz30Marker, 1.0, pageWidth-30)
			textHeight += mh + sh*0.5 //add the height of text + half of line space
		}

		dc = newPage(pt.Max.X-150, int(textHeight+footerMargin)+headerSize, "Top Artists:")
		verticalOffset := 0.0
		for i := 0; i < len(topList); i++ {
			mh, sh, _ := measureStringHeight(topList[i], sz30Marker, 1.0, pageWidth-30)

			dc.DrawStringWrapped(topList[i], 50, float64(headerSize)+(sh*(1.5)*float64(i))+verticalOffset, 0, 0, pageWidth, 1.0, gg.AlignLeft)

			if mh != sh {
				verticalOffset += (mh - sh)
			}
		}
	}

	var b bytes.Buffer
	dc.EncodePNG(&b)

	//send back png as data
	c.Data(http.StatusOK, "image/x-png", b.Bytes())
}

func generateTagline(timeFrame string) string {
	// Returns either the current month and year as a string, or 'All-Time'
	if timeFrame == "long_term" {
		return "All-Time"
	} else if timeFrame == "short_term" {
		//get month
		current := time.Now().UTC()
		m := time.Month(current.Month())
		y := current.Year()
		return m.String() + " " + strconv.Itoa(y)
	}
	return ""
}

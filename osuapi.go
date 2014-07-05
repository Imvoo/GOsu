package GOsu

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	API_URL          string = "https://osu.ppy.sh/api/"
	API_RECENT_PLAYS string = "get_user_recent"
	USER_ID          string
)

type Database struct {
	API_KEY string
}

type Song struct {
	Beatmap_ID   string
	Score        string
	MaxCombo     string
	Count50      string
	Count100     string
	Count300     string
	CountMiss    string
	CountKatu    string
	CountGeki    string
	Perfect      string
	Enabled_Mods string
	User_ID      string
	Date         string
	Rank         string
}

func (d *Database) SetAPIKey() error {
	tempKey, err := ioutil.ReadFile("./APIKEY.txt")

	// If there is no file, try find the API Key in the Environment Variables.
	if err != nil {
		d.API_KEY = os.Getenv("APIKEY")

		if len(d.API_KEY) <= 1 {
			err = errors.New("API Key: unable to locate API Key in environment variables or in local APIKEY.txt file.")
			return err
		} else {
			err = nil
		}
	} else {
		d.API_KEY = string(tempKey)
	}

	// Trims spaces and trailing newlines from the API key so that the URL
	// to retrieve songs can be built properly.
	d.API_KEY = strings.TrimSpace(d.API_KEY)
	d.API_KEY = strings.Trim(d.API_KEY, "\r\n")

	return err
}

func SetUserID(user string) {
	USER_ID = user
}

func (d Database) BuildRecentURL(IN_USER_ID string, GAME_TYPE int) string {
	return API_URL + API_RECENT_PLAYS + "?k=" + d.API_KEY + "&u=" + IN_USER_ID
}

func GetRecentPlays(url string) ([]Song, error) {
	var songs []Song

	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return nil, errors.New("HTTP: Could not open a connection to the Osu! API server.")
	}

	html, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.New("HTML: Could not read the HTML page grabbed.")
	}

	err = json.Unmarshal(html, &songs)

	if err != nil {
		return nil, errors.New("JSON: Couldn't process the HTML page grabbed, most likely due to not being in the right format. Make sure you aren't being redirected due to a network proxy or invalid API Key.")
	}

	return songs, err
}

// ONLY A TEMPORARY FUNCTION.
// Use this function if you are behind a proxy/corporate network and want to work off a local file.
// It will serve as a local HTML file for you to test the website.
func GetLocalPlays(path string) ([]Song, error) {
	var songs []Song

	html, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.New("HTML: Could not read the local HTML page properly.")
	}

	err = json.Unmarshal(html, &songs)

	if err != nil {
		return nil, errors.New("JSON: Couldn't process the local HTML page, most likely due to not being in the right format.")
	}

	return songs, err
}

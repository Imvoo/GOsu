package GOsu

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Types
const (
	OSU        = "0"
	TAIKO      = "1"
	CTB        = "2"
	MANIA      = "3"
	BEATMAPSET = "s"
	BEATMAPID  = "b"
	USERID     = "u"
)

var (
	API_URL          string = "https://osu.ppy.sh/api/"
	API_RECENT_PLAYS string = "get_user_recent"
	API_GET_BEATMAPS string = "get_beatmaps"
	API_GET_USER     string = "get_user"
)

type Database struct {
	API_KEY string
}

type Beatmap struct {
	Beatmapset_ID     string
	Beatmap_ID        string
	Approved          string
	Approved_Date     string
	Last_Update       string
	Total_Length      string
	Hit_Length        string
	Version           string
	Artist            string
	Title             string
	Creator           string
	Bpm               string
	Source            string
	Difficulty_Rating string
	Diff_Size         string
	Diff_Overall      string
	Diff_Approach     string
	Diff_Drain        string
	Mode              string
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

type User struct {
	User_ID       string
	Username      string
	Count300      string
	Count100      string
	Count50       string
	PlayCount     string
	Ranked_Score  string
	Total_Score   string
	PP_Rank       string
	Level         string
	PP_Raw        string
	Accuracy      string
	Count_Rank_SS string
	Count_Rank_S  string
	Count_Rank_A  string
	Country       string
	Events        []Event
}

type Event struct {
	Display_HTML  string
	Beatmap_ID    string
	Beatmapset_ID string
	Date          string
	EpicFactor    string
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

func (d Database) BuildRecentURL(USER_ID string, GAME_TYPE string) string {
	return API_URL + API_RECENT_PLAYS + "?k=" + d.API_KEY + "&u=" + USER_ID + "&m=" + GAME_TYPE
}

func (d Database) BuildBeatmapURL(ID string, TYPE string) string {
	return API_URL + API_GET_BEATMAPS + "?k=" + d.API_KEY + "&" + TYPE + "=" + ID
}

func (d Database) BuildUserURL(USER_ID string, GAME_TYPE string, DAYS string) string {
	return API_URL + API_GET_USER + "?k=" + d.API_KEY + "&u=" + USER_ID + "&m=" + GAME_TYPE + "&event_days=" + DAYS
}

func RetrieveHTML(URL string) ([]byte, error) {
	res, err := http.Get(URL)
	defer res.Body.Close()

	if err != nil {
		return nil, errors.New("HTTP: Could not open a connection to the Osu! API server.")
	}

	html, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.New("HTML: Could not read the HTML page grabbed.")
	}

	return html, err
}

func (d Database) GetUser(USER_ID string, GAME_TYPE string, DAYS string) ([]User, error) {
	var user []User
	url := d.BuildUserURL(USER_ID, GAME_TYPE, DAYS)
	html, err := RetrieveHTML(url)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(html, &user)

	if err != nil {
		return nil, errors.New("JSON: Couldn't process HTML into JSON data. You might have the wrong page or a wrong API key. The HTML grabbed at " + url + " will be displayed below:\n" + string(html))
	}

	return user, err
}

func (d Database) GetBeatmaps(ID string, TYPE string) ([]Beatmap, error) {
	var beatmaps []Beatmap
	url := d.BuildBeatmapURL(ID, TYPE)
	html, err := RetrieveHTML(url)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(html, &beatmaps)

	if err != nil {
		return nil, errors.New("JSON: Couldn't process HTML into JSON data. You might have the wrong page or a wrong API key. The HTML grabbed at " + url + " will be displayed below:\n" + string(html))
	}

	return beatmaps, err
}

func (d Database) GetRecentPlays(USER_ID string, GAME_TYPE string) ([]Song, error) {
	var songs []Song
	url := d.BuildRecentURL(USER_ID, GAME_TYPE)
	html, err := RetrieveHTML(url)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(html, &songs)

	if err != nil {
		return nil, errors.New("JSON: Couldn't process HTML into JSON data. You might have the wrong page or a wrong API key. The HTML grabbed at " + url + " will be displayed below:\n" + string(html))
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

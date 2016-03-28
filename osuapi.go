package GOsu

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Game Types which can be accessed via GOsu.<GAMETYPE>.
const (
	OSU   = "0"
	TAIKO = "1"
	CTB   = "2"
	MANIA = "3"
)

// Beatmap Identifier Types for use with GetBeatmaps.
const (
	BEATMAPSET = "s"
	BEATMAPID  = "b"
	USERID     = "u"
)

// Osu Mods which are accessed by the GetMod function.
// iota moves the 'binary 1' to the left.
// e.g. NoFail = 1, Easy = 10 in binary which is 2.
// NoVideo = 100 in binary which is 4.
// This is done in order for any combination of mods to be enabled at once.
const (
	None   = 0
	NoFail = 1 << (iota - 1)
	Easy
	NoVideo
	Hidden
	HardRock
	SuddenDeath
	DoubleTime
	Relax
	HalfTime
	Nightcore
	Flashlight
	Autoplay
	SpunOut
	Relax2
	Perfect
	Key4
	Key5
	Key6
	Key7
	Key8
	FadeIn
	Random
	LastMod
)

var (
	API_URL           string = "https://osu.ppy.sh/api/"
	API_RECENT_PLAYS  string = "get_user_recent"
	API_GET_BEATMAPS  string = "get_beatmaps"
	API_GET_USER      string = "get_user"
	API_GET_SCORES    string = "get_scores"
	API_GET_USER_BEST string = "get_user_best"
	API_GET_MATCH     string = "get_match"
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
	PP_Country_Rank  string
	Events        []Event
}

type Event struct {
	Display_HTML  string
	Beatmap_ID    string
	Beatmapset_ID string
	Date          string
	EpicFactor    string
}

type Score struct {
	Score        string
	Username     string
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
	PP           string
}

type PPSong struct {
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
	PP           string
}

type Game struct {
	Match MPMatch
	Games []MPGame
}

type MPMatch struct {
	Match_ID   string
	Name       string
	Start_Time string
	End_Time   string
}

type MPGame struct {
	Game_ID      string
	Start_Time   string
	End_Time     string
	Beatmap_ID   string
	Play_Mode    string
	Match_Type   string
	Scoring_Type string
	Team_Type    string
	Mods         string
	Scores       []MPScore
}

type MPScore struct {
	Slot      string
	Team      string
	User_ID   string
	Score     string
	MaxCombo  string
	Rank      string
	Count50   string
	Count100  string
	Count300  string
	CountMiss string
	CountGeki string
	CountKatu string
	Perfect   string
	Pass      string
}

func (d *Database) SetAPIKey(API_KEY string) error {
	var err error
	if len(API_KEY) <= 1 {
		err = errors.New("API Key: invalid/non-existant API key.")
		return err
	} else {
		d.API_KEY = API_KEY
	}

	return err
}

func (d Database) BuildRecentURL(USER_ID string, GAME_TYPE string) string {
	return API_URL + API_RECENT_PLAYS + "?k=" + d.API_KEY + "&u=" + USER_ID + "&m=" + GAME_TYPE + "&limit=50"
}

func (d Database) BuildBeatmapURL(ID string, TYPE string) string {
	return API_URL + API_GET_BEATMAPS + "?k=" + d.API_KEY + "&" + TYPE + "=" + ID
}

func (d Database) BuildUserURL(USER_ID string, GAME_TYPE string, DAYS string) string {
	return API_URL + API_GET_USER + "?k=" + d.API_KEY + "&u=" + USER_ID + "&m=" + GAME_TYPE + "&event_days=" + DAYS
}

func (d Database) BuildUserBestURL(USER_ID string, GAME_TYPE string) string {
	return API_URL + API_GET_USER_BEST + "?k=" + d.API_KEY + "&u=" + USER_ID + "&m=" + GAME_TYPE
}

func (d Database) BuildScoreURL(BEATMAP_ID string, USER_ID string, GAME_TYPE string) string {
	return API_URL + API_GET_SCORES + "?k=" + d.API_KEY + "&b=" + BEATMAP_ID + "&m=" + GAME_TYPE + "&u=" + USER_ID
}

func (d Database) BuildMatchURL(MATCH_ID string) string {
	return API_URL + API_GET_MATCH + "?k=" + d.API_KEY + "&mp=" + MATCH_ID
}

func RetrieveHTML(URL string) ([]byte, error) {
	res, err := http.Get(URL)

	if err != nil {
		return nil, errors.New("HTTP: Could not open a connection to the Osu! API server.")
	}

	defer res.Body.Close()

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
		return user, err
	}

	err = json.Unmarshal(html, &user)

	if err != nil {
		return user, errors.New("JSON: Couldn't process the HTML into JSON. Possible wrong page/wrong API key/offline server.")
	}

	return user, err
}

func (d Database) GetBeatmaps(ID string, BM_TYPE string) ([]Beatmap, error) {
	var beatmaps []Beatmap
	url := d.BuildBeatmapURL(ID, BM_TYPE)
	html, err := RetrieveHTML(url)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(html, &beatmaps)

	if err != nil {
		return nil, errors.New("JSON: Couldn't process the HTML into JSON. Possible wrong page/wrong API key/offline server.")
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
		return nil, errors.New("JSON: Couldn't process the HTML into JSON. Possible wrong page/wrong API key/offline server.")
	}

	return songs, err
}

func (d Database) GetScores(BEATMAP_ID string, USER_ID string, GAME_TYPE string) ([]Score, error) {
	var scores []Score
	url := d.BuildScoreURL(BEATMAP_ID, USER_ID, GAME_TYPE)
	html, err := RetrieveHTML(url)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(html, &scores)

	if err != nil {
		return nil, errors.New("JSON: Couldn't process the HTML into JSON. Possible wrong page/wrong API key/offline server.")
	}

	return scores, err
}

func (d Database) GetUserBest(USER_ID string, GAME_TYPE string) ([]PPSong, error) {
	var songs []PPSong
	url := d.BuildUserBestURL(USER_ID, GAME_TYPE)
	html, err := RetrieveHTML(url)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(html, &songs)

	if err != nil {
		return nil, errors.New("JSON: Couldn't process the HTML into JSON. Possible wrong page/wrong API key/offline server.")
	}

	return songs, err
}

func (d Database) GetMatch(MATCH_ID string) (Game, error) {
	var game Game
	url := d.BuildMatchURL(MATCH_ID)
	html, err := RetrieveHTML(url)

	if err != nil {
		return game, err
	}

	err = json.Unmarshal(html, &game)

	if err != nil {
		return game, errors.New("JSON: Couldn't process the HTML into JSON. Possible wrong page/wrong API key/offline server.")
	}

	return game, err
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

func GetMods(mods int) []string {
	var resultantMods []string

	if mods == 0 {
		resultantMods = append(resultantMods, "None")
		return resultantMods
	}

	AddMod(&resultantMods, &mods, LastMod, "Last Mod")
	AddMod(&resultantMods, &mods, Random, "Random")
	AddMod(&resultantMods, &mods, FadeIn, "Fade In")
	AddMod(&resultantMods, &mods, Key8, "Key 8")
	AddMod(&resultantMods, &mods, Key7, "Key 7")
	AddMod(&resultantMods, &mods, Key6, "Key 6")
	AddMod(&resultantMods, &mods, Key5, "Key 5")
	AddMod(&resultantMods, &mods, Key4, "Key 4")
	AddMod(&resultantMods, &mods, Perfect, "Perfect")
	AddMod(&resultantMods, &mods, Relax2, "Relax 2")
	AddMod(&resultantMods, &mods, SpunOut, "Spun Out")
	AddMod(&resultantMods, &mods, Autoplay, "Autoplay")
	AddMod(&resultantMods, &mods, Flashlight, "Flashlight")
	AddMod(&resultantMods, &mods, Nightcore, "Nightcore")
	AddMod(&resultantMods, &mods, HalfTime, "Half Time")
	AddMod(&resultantMods, &mods, Relax, "Relax")
	AddMod(&resultantMods, &mods, DoubleTime, "Double Time")
	AddMod(&resultantMods, &mods, SuddenDeath, "Sudden Death")
	AddMod(&resultantMods, &mods, HardRock, "Hard Rock")
	AddMod(&resultantMods, &mods, Hidden, "Hidden")
	AddMod(&resultantMods, &mods, NoVideo, "No Video")
	AddMod(&resultantMods, &mods, Easy, "Easy")
	AddMod(&resultantMods, &mods, NoFail, "No Fail")

	return resultantMods
}

func AddMod(array *[]string, remainingMods *int, mod int, modName string) {
	if mod <= *remainingMods {
		*remainingMods -= mod
		*array = append(*array, modName)
	}
}

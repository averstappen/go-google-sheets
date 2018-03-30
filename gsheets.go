package gsheets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var (
	Ctx    context.Context
	client *http.Client
	Srv    *sheets.Service
)

func init() {
	Ctx = context.Background()
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file %v\n", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v\n", err)
	}

	client = getClient(Ctx, config)

	Srv, err = sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v\n", err)
	}
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(Ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(Ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	pwd, err := os.Getwd()
	tokenCacheDir := filepath.Join(pwd, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gsheets.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func ValueRange() *sheets.ValueRange {
	return &sheets.ValueRange{}
}

// GetValues gets the values from a google sheet range
// it takes a sheetId and a range as arguments
func GetValues(sheetId string, readRange string) (*sheets.ValueRange, error) {
	return Srv.Spreadsheets.Values.Get(sheetId, readRange).Do()
}

// UpdateValues updates a google sheet range
func UpdateValues(sheetId string, writeRange string, values *sheets.ValueRange, inputoption string) (*sheets.UpdateValuesResponse, error) {
	return Srv.Spreadsheets.Values.Update(sheetId, writeRange, values).ValueInputOption(inputoption).Context(Ctx).Do()
}

// NumberToLetter returns the column string for column number
func NumberToLetter(number int) string {
	if number <= 0 {
		return ""
	}

	return NumberToLetter(int((number-1)/26)) + string(byte(65+(number-1)%26))
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

type APICred struct {
	APIKEY       string
	APISECRET    string
	ACCESSTOKEN  string
	ACCESSSECRET string
}

type Quotes struct {
	QuoteList []QuoteObject `json:"quotes"`
}

type QuoteObject struct {
	Tweetid int    `json:"tweetid"`
	Quote   string `json:"quote"`
}

func LoadEnv() (env APICred) {
	err := godotenv.Load()

	// Establish credentials for accessing API
	env = APICred{
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
		os.Getenv("ACCESS_TOKEN"),
		os.Getenv("ACCESS_SECRET")}
	// Credential error checking
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return env
}

// Random number generator used to grab index of quote
// Max should be the number of quotes in the JSON file
// TODO: Create a dynamic list to make sure the last few indexes were not chosen
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// Get quote to serve to API
// TODO: Convert to grabbing from DB
func GrabQuote() string {
	quote_file, err := os.Open("quotes.json")
	if err != nil {
		log.Fatal(err)
	}

	defer quote_file.Close()

	quote_bytes, _ := ioutil.ReadAll(quote_file)

	var quotes Quotes

	error := json.Unmarshal(quote_bytes, &quotes)
	if error != nil {
		log.Fatal(error)
	}

	// Count number of quotes available
	max_quotes := len(quotes.QuoteList) - 1

	// Get random index
	random_tweetid := random(0, max_quotes)
	// Get a random quote from the list using the random index above
	quote_to_serv := quotes.QuoteList[random_tweetid].Quote

	return quote_to_serv
}

func SendTweet() {
	env := LoadEnv()

	anaconda.SetConsumerKey(env.APIKEY)
	anaconda.SetConsumerSecret(env.APISECRET)

	api := anaconda.NewTwitterApi(env.ACCESSTOKEN, env.ACCESSSECRET)

	tweet := GrabQuote()

	_, err := api.PostTweet(tweet, url.Values{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(SendTweet)
	//SendTweet()
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dghubble/oauth1"
)

func CheckError(e error) {
	if e != nil {
		log.Printf(e.Error())
	}
}

type OAuth1Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessTokenKey    string
	AccessTokenSecret string
}

func CreateClient(twitterAuthConfig OAuth1Config) *http.Client {
	config := oauth1.NewConfig(twitterAuthConfig.ConsumerKey, twitterAuthConfig.ConsumerSecret)
	token := oauth1.NewToken(twitterAuthConfig.AccessTokenKey, twitterAuthConfig.AccessTokenSecret)

	return config.Client(oauth1.NoContext, token)
}

func GetTweetsByHashtag(client *http.Client, hashtag string) map[string]interface{} {
	req, _ := http.NewRequest(http.MethodGet, TWITTER_SEARCH_URL, nil)

	//Query Params
	params := req.URL.Query()
	params.Add("q", hashtag)
	params.Add("result_type", "recent")
	params.Add("include_entities", "true")
	params.Add("until", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	req.URL.RawQuery = params.Encode()

	response, err := client.Do(req)
	fmt.Println(req.URL.String())
	CheckError(err)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	var parsedTweets map[string]interface{}
	json.Unmarshal(body, &parsedTweets)
	return parsedTweets
}

// GetTweetID takes a single tweet and returns its ID
func GetTweetID(tweet map[string]interface{}) string {
	// TODO: auto parsing of tweet id as float64 is messing the numbers
	// hence using id_str for now
	return tweet["id_str"].(string)
}

func Retweet(client *http.Client, tweetID string) error {
	fmt.Println("Retweeting:", tweetID)
	response, err := client.Post(fmt.Sprintf(RETWEET_URL, tweetID), "application/json", nil)
	CheckError(err)
	fmt.Println(response)
	return nil
}

func Like(client *http.Client, tweetID string) error {
	fmt.Println("Liking:", tweetID)
	response, err := client.Post(fmt.Sprintf(LIKE_URL, tweetID), "application/json", nil)
	CheckError(err)
	fmt.Println(response)
	return nil
}

func main() {

	var twitterAuthConfig OAuth1Config
	_, err := toml.DecodeFile("AuthConfig.toml", &twitterAuthConfig)
	fmt.Println(twitterAuthConfig)
	CheckError(err)

	twitterClient := CreateClient(twitterAuthConfig)

	tweets := GetTweetsByHashtag(twitterClient, "100daysofcode")

	for _, v := range tweets["statuses"].([]interface{}) {
		tweet := v.(map[string]interface{})
		tweetID := GetTweetID(tweet)
		Retweet(twitterClient, tweetID)
		Like(twitterClient, tweetID)
	}
}

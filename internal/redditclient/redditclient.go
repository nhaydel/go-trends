package redditclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"math/rand"
	"strconv"
)

func GetNewPosts(subreddit string) (*http.Response, error) {
	baseurl := "https://reddit.com/r/" + subreddit + "/new.json"
	client := http.Client{}
	req, err := http.NewRequest("GET", baseurl, nil)
	if err != nil {
		return nil, err
	}
	user_agent := strconv.Itoa(rand.Int())
	req.Header.Set("User-Agent",  user_agent)
	return client.Do(req)
}

func ParsePostData(response *http.Response) []map[string]string {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var resp map[string]interface{}
	posts := []map[string]string{}
	json.Unmarshal(body, &resp)
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		log.Printf("got data of type %T but wanted map[string]interface{}", resp["data"])
	} else {
		children := data["children"].([]interface{})
		for _, child := range children {
			child_data := child.(map[string]interface{})["data"].(map[string]interface{})
			post_data := make(map[string]string)
			post_data["title"] = ToString(child_data["title"])
			post_data["id"] = ToString(child_data["id"])
			post_data["author"] = ToString(child_data["author"])
			posts = append(posts, post_data)
		}
	}

	return posts
}

func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

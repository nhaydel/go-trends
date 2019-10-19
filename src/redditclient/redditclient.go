package redditclient

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

func GetNewPosts(subreddit string) (*http.Response, error) {
	baseurl := "https://reddit.com/r/" + subreddit + "/new.json"
	client := http.Client{}
	req, err := http.NewRequest("GET", baseurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Go-Trends")
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
	data := resp["data"].(map[string]interface{})
	children := data["children"].([]interface{})
	for _, child := range children {
		child_data := child.(map[string]interface{})["data"].(map[string]interface{})
		post_data := make(map[string]string)
		post_data["title"] = ToString(child_data["title"])
		post_data["id"] = ToString(child_data["id"])
		post_data["author"] = ToString(child_data["author"])
		posts = append(posts, post_data)
	}

	return posts
}

func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// func main() {
// 	subreddit := "news"
// 	resp, err := GetNewPosts(subreddit)
// 	if err != nil {
// 		fmt.Println("Error pulling posts from subreddit: " + subreddit)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	fmt.Println(string(body))
// }
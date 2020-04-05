package main

import (
	"flag"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
	"gopkg.in/jdkato/prose.v2"
	client "github.com/nhaydel/go-trends/internal/redditclient"
	structures "github.com/nhaydel/go-trends/internal/structures"
	trendsmap "github.com/nhaydel/go-trends/internal/trendsmap"
)

func PostExtractor(id int, queue *structures.CircularQueue, channel chan map[string]string,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		subreddit := queue.PopString()
		response, err := client.GetNewPosts(subreddit)
		if err != nil {
			panic(err)
		}
		data := client.ParsePostData(response)
		for _, post := range data {
			channel <- post
		}
	}
}

func PostTransformer(id int, set *structures.SyncSet,
	postChannel chan map[string]string,
	titleChannel chan string,
	wg *sync.WaitGroup) {
	defer wg.Done()
	reg, err := regexp.Compile("[^a-zA-Z0-9_\t\n\f\r ]+")
	if err != nil {
		log.Fatal(err)
	}
	for {
		post := <-postChannel
		if set.Insert(post["id"]) {
			titleChannel <- reg.ReplaceAllString(post["title"], "")
		}
	}
}

func PostLoader(id int, trends *trendsmap.TrendsMap, titleChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		title := <-titleChannel
		doc, _ := prose.NewDocument(title)
		var important_words []string
		for _, ent := range doc.Tokens() {
			if strings.HasPrefix(ent.Tag, "NN") || strings.HasPrefix(ent.Tag, "VB") {
				important_words = append(important_words, ent.Text)
			}
		}
		trends.Add(title, important_words)
	}
}

func CheckTrends(trends *trendsmap.TrendsMap, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		trends.DisplayTrends()
		time.Sleep(10000 * time.Millisecond)
	}

}

func main() {
	queue := structures.NewQueue()
	queue.InitFromFile("./resources/subreddits.txt")
	set := structures.NewSyncSet()
	trends := trendsmap.NewTrendsMap()
	// set up channels
	postsChannel := make(chan map[string]string)
	titlesChannel := make(chan string)
	var wg sync.WaitGroup
	defer close(postsChannel)
	defer close(titlesChannel)

	// command line arguments
	extractors := flag.Int("e", 5, "The number of extractors you want pulling data from reddit.")
	transformers := flag.Int("t", 5, "The number of transformers you want parsing reddit posts.")
	loaders := flag.Int("l", 5, "The number of loaders you want loading titles from reddit posts.")
	flag.Parse()

	// create goroutines
	for i := 0; i < *extractors; i++ {
		wg.Add(1)
		go PostExtractor(i, queue, postsChannel, &wg)
	}

	for j := 0; j < *transformers; j++ {
		wg.Add(1)
		go PostTransformer(j, set, postsChannel, titlesChannel, &wg)
	}

	for k := 0; k < *loaders; k++ {
		wg.Add(1)
		go PostLoader(k, trends, titlesChannel, &wg)
	}

	wg.Add(1)
	go CheckTrends(trends, &wg)
	wg.Wait()
}

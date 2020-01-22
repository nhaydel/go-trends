package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	graph "../../internal/graph"
	client "../../internal/redditclient"
	structures "../../internal/structures"

	"gopkg.in/jdkato/prose.v2"
)

type byPathLength []*graph.Path

func (s byPathLength) Len() int {
	return len(s)
}
func (s byPathLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPathLength) Less(i, j int) bool {
	return s[i].Length() < s[j].Length()
}

var pos []string

//Initialize the parts of speech we want in our graph
var acceptablePOS *structures.Set

func PostExtractor(id int, queue *structures.CircularQueue, channel chan map[string]string,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		subreddit := queue.PopString()
		// fmt.Println("Extractor ", id, " extracting posts from subreddit: ", subreddit)
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
	for {
		post := <-postChannel
		if set.Insert(post["id"]) {
			// fmt.Println("Transformer ", id, " parsing post ", post["id"])
			doc, _ := prose.NewDocument(post["title"])
			formattedTitle := ""
			for _, tok := range doc.Tokens() {
				if acceptablePOS.Contains(tok.Tag) {
					formattedTitle = formattedTitle + " " + tok.Text
				}
			}
			titleChannel <- FormatPostTitle(formattedTitle)
		}
	}
}

func FormatPostTitle(title string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9_\t\n\f\r ]+")
	if err != nil {
		log.Fatal(err)
	}
	lowerTitle := strings.ToLower(title)
	return reg.ReplaceAllString(lowerTitle, "")
}

func processTitle(titleGraph *graph.WeightedGraph, title string) {
	words := strings.Split(title, " ")
	for i := 1; i < len(words)-1; i++ {
		titleGraph.AddEdge(words[i], words[i+1])
	}
}

func PostLoader(id int, titleGraph *graph.WeightedGraph, titleChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		title := <-titleChannel
		processTitle(titleGraph, title)
	}
}

func FindTrends(titleGraph *graph.WeightedGraph, depth int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Println("Current Trends: ")
		paths := titleGraph.GetAllPathsToDepth(depth)
		sort.Reverse(byPathLength(paths))
		i := 0
		for _, path := range paths {
			fmt.Println(strings.Join(path.Steps(), " "))
			i++
			if i == 5 {
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	queue := structures.NewQueue()
	queue.InitFromFile("../../resources/subreddits.txt")
	set := structures.NewSyncSet()
	words := graph.NewWeightedGraph()
	pos = []string{"JJ", "MD", "VB", "NN", "SYM", "JJS", "JJP", "NNP", "NNS", "NNPS", "VBD", "VBG", "VBN", "VBP", "VBZ"}
	acceptablePOS = structures.NewSet()
	for _, tag := range pos {
		acceptablePOS.Insert(tag)
	}
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
	depth := flag.Int("d", 2, "The maximum size of a generated trend phrase")
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
		go PostLoader(k, words, titlesChannel, &wg)
	}

	wg.Add(1)
	go FindTrends(words, *depth, &wg)
	wg.Wait()
}

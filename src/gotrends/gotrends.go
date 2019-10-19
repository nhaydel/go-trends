package main

import (
	"fmt"
	"sync"
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"regexp"
	"gopkg.in/jdkato/prose.v2"
	client "redditclient"
)
type Item interface{}

type CircularQueue struct {
	items []Item
	lock sync.RWMutex
}

type SyncSet struct {
	items map[Item] bool
	lock sync.RWMutex
}

type SyncMap struct {
	data map[Item]Item
	lock sync.RWMutex
}

func NewQueue() *CircularQueue {
	return &CircularQueue{items: make([]Item, 0)}
}

func (queue *CircularQueue) Push (item interface{}) {
	queue.lock.Lock()
	queue.items = append(queue.items, item)
	queue.lock.Unlock()
}

func (queue *CircularQueue) InitFromFile (filepath string) {
	file, err := os.Open(filepath)
    if err != nil {
        log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
        queue.Push(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func (queue *CircularQueue) PopString () string {
	queue.lock.Lock()
	item := queue.items[0]
	if (len(queue.items) > 1) {
		queue.items = queue.items[1:]
		queue.items = append(queue.items, item)
	}
	queue.lock.Unlock()
	return fmt.Sprintf("%v", item)

}

func (queue *CircularQueue) PrintStatus () {
	for _, item := range queue.items[0:] {
		fmt.Println(item)
	}
}

func NewSyncSet () *SyncSet {
	return &SyncSet{items: make(map[Item]bool)}
}

func (set *SyncSet) Insert(item Item) bool {
	set.lock.Lock()
	if set.Contains(item) {
		set.lock.Unlock()
		return false
	}

	set.items[item] = true
	set.lock.Unlock()
	return true
}

func (set *SyncSet) Contains (item Item) bool {
	_,ok := set.items[item]
	return ok
}

func NewSyncMap () *SyncMap {

}

func PostExtractor(id int, queue *CircularQueue, channel chan map[string]string, 
						wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		subreddit := queue.PopString()
		fmt.Println("Extractor ", id, " extracting posts from subreddit: ", subreddit)
		response, err := client.GetNewPosts("news")
		if err != nil {
			panic(err)
		}
		data := client.ParsePostData(response)
		for _, post := range data {
			channel <- post
		}
	}
}

func PostTransformer(id int, set *SyncSet, postChannel chan map[string] string, 
						wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		post := <- postChannel
		if set.Insert(post["id"]) {
			fmt.Println("Transformer ", id, " parsing post ", post["id"])
			doc, _ := prose.NewDocument(post["title"])
			fmt.Println(doc.Text)

			for _, ent := range doc.Entities() {
				fmt.Println(ent.Text, ent.Label)
			}
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

// func PostLoader(dataMap *SyncMap) {

// }

func main() {
	queue := NewQueue()
	queue.InitFromFile("../../resources/subreddits.txt")
	set := NewSyncSet()

	// set up channels
	postsChannel := make(chan map[string]string)
	entitiesChannel := make(chan string)
	var wg sync.WaitGroup
	defer close(postsChannel)

	// command line arguments
	extractors := flag.Int("e", 5, "The number of extractors you want pulling data from reddit.")
	transformers := flag.Int("t", 5, "The number of transformers you want parsing reddit posts.")
	flag.Parse()

	// create goroutines
	for i := 0; i < *extractors; i++ {
		wg.Add(1)
		go PostExtractor(i, queue, postsChannel, &wg)
	}

	for j := 0; j < *transformers; j++ {
		wg.Add(1)
		go PostTransformer(j, set, postsChannel, &wg)
	}

	wg.Wait()
}

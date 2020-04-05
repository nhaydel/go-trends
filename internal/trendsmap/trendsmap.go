package structures

import (
	"fmt"
	"sync"
	"strings"
	"gopkg.in/jdkato/prose.v2"
	"time"
)

type TrendsMap struct {
	data map[string]int
	top []string
	lock sync.RWMutex
}


func NewTrendsMap() *TrendsMap {
	return &TrendsMap{data: make(map[string]int)}
}

func sliceToSet(tokens []string) map[string]bool {
	set := make(map[string]bool)
	for i := range tokens {
		set[tokens[i]] = true
	}
	return set
}

func tokenSimilarity(tokens []string, insert_tokens []string) float64 {
	set := sliceToSet(tokens)
	common := 0
	for i := range insert_tokens {
		_, ok := set[insert_tokens[i]]
		if ok {
			common = common + 1
		}
	}
	return float64(common/len(tokens))
}

func (trends *TrendsMap) Add(sentence string, tokens []string) {
	trends.lock.Lock()
	for topic, v := range trends.data {
		doc, _ := prose.NewDocument(topic)
		var important_words []string
		for _, ent := range doc.Tokens() {
			if strings.HasPrefix(ent.Tag, "NN") || strings.HasPrefix(ent.Tag, "VB") {
				important_words = append(important_words, ent.Text)
			}
		}
		if tokenSimilarity(tokens, important_words) > 0.5 {
			trends.data[topic] = v + 1
			trends.updateTopTrends(topic)
			trends.lock.Unlock()
			return
		}
	}

	trends.data[sentence] = 1
	trends.updateTopTrends(sentence)
	trends.lock.Unlock()
}

func (trends *TrendsMap) updateTopTrends(title string) {
	if len(trends.top) < 10 {
		trends.top = append(trends.top, title)
	} else {
		for i, topic := range(trends.top) {
			if trends.data[topic] < trends.data[title] {
				trends.top[i] = title
				return
			}
		}
	}
}

func (trends *TrendsMap) Print() {
	for k, v := range trends.data {
		fmt.Println(fmt.Sprintf("%s: %d", k, v))
	}
}

func (trends *TrendsMap) DisplayTrends() {
	trends.lock.Lock()
	fmt.Println(fmt.Sprintf("Top trends as of %s \n\n", time.Now().String()))
	for _, topic := range trends.top {
		fmt.Println(fmt.Sprintf("%s: %d", topic, trends.data[topic]))
	}
	fmt.Println("\n\n")
	trends.lock.Unlock()
}

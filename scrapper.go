package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/akhiltn/rssagg-go/internal/database"
)

func startScrapper(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scapping on %v goroutine every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, feed, db)
		}
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, feed database.Feed, db *database.Queries) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
	}
	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

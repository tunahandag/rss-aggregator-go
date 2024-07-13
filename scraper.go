package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tunahandag/rss-aggregator/internal/database"
)




func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Starting scrapping with %d concurrency and %s between requests", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <- ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(),int32(concurrency))
		if err != nil {
			log.Printf("Error getting feeds to fetch: %v", err)
			continue
		}
		
		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			
			go scrapeFeed(db,wg,feed)
		}
		wg.Wait();
} }

func scrapeFeed (db *database.Queries,wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()	

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post: ", item.Title, " on feed ", feed.Name)
		description := sql.NullString{}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			continue
		}

		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID: feed.ID,
			Title: item.Title,
			Description: description,
			Url: item.Link,
			PublishedAt: pubAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error creating post: %v", err)
		}


	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
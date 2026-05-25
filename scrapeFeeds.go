package main 

import (
	"fmt"
	"context"
)

func scrapeFeeds(s *state) error{ 
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err!= nil{ 
		return fmt.Errorf("%v", err)
	}

	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	RSSFeed, err := fetchFeed(context.Background(), feed.Url)
	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	for _, item := range RSSFeed.Channel.Item{ 
		fmt.Printf("%s: %s\n", feed.Name, item.Title)
	}

	return nil
}


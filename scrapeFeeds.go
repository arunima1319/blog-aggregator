package main 

import (
	"fmt"
	"context"
	"time"
	"github.com/araddon/dateparse"
	"database/sql"
	"github.com/arunima1319/blog-aggregator/internal/database"
	"github.com/google/uuid"

)

func scrapeFeeds(s *state) error{ 
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err!= nil{ 
		return fmt.Errorf("%v", err)
	}
	feedNull := uuid.NullUUID{
		UUID: feed.ID, 
		Valid: true,
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

		layout, err := dateparse.ParseFormat(item.PubDate)
		if err!= nil{
			return fmt.Errorf("%v", err)
		}
		timePublishedAt, err := time.Parse(layout, item.PubDate)
		if err!= nil{
			return fmt.Errorf("%v", err)
		}
		timeNull := sql.NullTime{
			Time: timePublishedAt,
			Valid: true,
		}

		descriptionNull := sql.NullString{
			String: item.Description, 
			Valid: true,
		}

		_, err = s.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID: uuid.New(), 
				CreatedAt: time.Now(), 
				UpdatedAt: time.Now(), 
				Title: item.Title, 
				Url: item.Link, 
				Description: descriptionNull, 
				PublishedAt: timeNull,
				FeedID: feedNull, 
			})
		if err!= nil{
			return fmt.Errorf("%v", err)
		}
		//fmt.Printf("%s: %s\n", feed.Name, item.Title)
	}

	return nil
}


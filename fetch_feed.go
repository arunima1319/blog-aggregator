package main 

import (
	"net/http"
	"fmt"
	"html"
	"context"
	"io"
	"encoding/xml"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){ 

	
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err!= nil{ 
		return nil, fmt.Errorf("Error Creating Request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err!= nil{
		return nil, fmt.Errorf("Error getting response: %v", err)
	}

	defer res.Body.Close()


	data, err := io.ReadAll(res.Body)
	if err!= nil{ 
		return nil, fmt.Errorf("Error reading the data: %v", err)
	}
	//fmt.Printf("%s\n", string(data))

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err!= nil{ 
		return nil, fmt.Errorf("Error unmarshaling into struct: %v", err)
	}
	//fmt.Printf("%v\n", feed)

	pointerFeed := &feed
	pointerFeed.Channel.Title = html.UnescapeString(pointerFeed.Channel.Title)
	pointerFeed.Channel.Description = html.UnescapeString(pointerFeed.Channel.Description)
	for _, item := range pointerFeed.Channel.Item{
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return pointerFeed, nil
}


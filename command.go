package main 

import (
	"fmt"
	"os"
	"github.com/google/uuid"
	"context"
	"time"
	"github.com/arunima1319/blog-aggregator/internal/database"

)

type command struct{ 
	name string 
	arguments []string
}

type commands struct{ 
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error{
	commandFunc, ok := c.commandMap[cmd.name]
	if ok{
		err := commandFunc(s, cmd)
		if err!= nil{
			return err
		}
	} else{
		return fmt.Errorf("command not found")
	}

	return nil
}

func (c *commands) register (name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func handlerLogin(s *state, cmd command) error{ 
	if len(cmd.arguments) == 0{
		return fmt.Errorf("no username argument passed")
		
	}
	username := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err!= nil{
		return fmt.Errorf("The name '%s' was not found in the database", username)
	}
	if err:=s.pointerConfig.SetUser(username); err!= nil{
		return err
	}
	fmt.Println("The user has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error{
	if len(cmd.arguments) == 0{
		return fmt.Errorf("no username argument passed")
		os.Exit(1)
	}
	name := cmd.arguments[0]

	user, err := s.db.CreateUser(
		context.Background(), 
		database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(), 
			Name: name,
		})
	if err!= nil{
		return fmt.Errorf("The name '%s' already exists", name) //
		os.Exit(1)
	}

	err = s.pointerConfig.SetUser(name)

	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	fmt.Println("The user was created: %v", user)

	return nil

} 

func handlerReset(s *state, cmd command) error { 
	err := s.db.DeleteUsers(context.Background())
	if err!=nil{ 
		return fmt.Errorf("%v", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error { 
	currentUser := s.pointerConfig.CurrentUserName

	names, err := s.db.GetAllUsers(context.Background())
	if err!=nil{
		return fmt.Errorf("%v", err)
	}
	for _, name := range names{ 
		if name == currentUser{
			fmt.Printf("* %s (current)\n", name)
		}else{
			fmt.Printf("* %s\n", name)
		}
	}
	return nil
}

func handlerAgg(_ *state, cmd command) error{ 
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("%v", *feed)
	return nil
	
}

func handlerAddfeed(s *state, cmd command, user database.User) error{
	if len(cmd.arguments) < 2{ 
		return fmt.Errorf("not enough arguments passed")
		os.Exit(1)
	} 
	nameFeed := cmd.arguments[0]
	url := cmd.arguments[1]

	//instead of getting the user like this, we are just using the user that has been passed into the function
	/*
	user, err := s.db.GetUser(context.Background(), s.pointerConfig.CurrentUserName)
	if err!= nil{ 
		return fmt.Errorf("error in getting current user: %v", err)
	}	
	*/

	userUUID := uuid.NullUUID{ 
		UUID : user.ID, 
		Valid: true, 
	}

	feed, err := s.db.CreateFeed(
		context.Background(), 
		database.CreateFeedParams{
			ID : uuid.New(),
			CreatedAt : time.Now(),
			UpdatedAt : time.Now(),
			Name : nameFeed,
			Url: url, 
			UserID : userUUID,

	})
	if err!= nil{
		return fmt.Errorf("error in creating new feed: %v", err)
	}

	followCommand := command{name: "follow", arguments: []string{url}}
	err = handlerFollow(s, followCommand, user)
	if err!=nil{
		return fmt.Errorf("error in following feed: %v", err)
	}

	fmt.Println("%v", feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error{ 

	feeds, err := s.db.ListFeeds(context.Background())
	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	for _, feed := range feeds{

		if feed.UserID.Valid{
			userId := feed.UserID.UUID
			user, err := s.db.GetUserFromID(context.Background(), userId)
			if err!= nil{
				return fmt.Errorf("%v", err)
			}else{
				fmt.Printf("name: %s, url: %s, user: %s\n", feed.Name, feed.Url, user.Name)
			}
		}else{
			return fmt.Errorf("the user ID was Null")
		}
	
	}

	return nil 
	
}


func handlerFollow(s *state, cmd command, user database.User) error{

	if len(cmd.arguments) == 0{
		return fmt.Errorf("no url passed")
	}

	urlFeed := cmd.arguments[0]
	feed, err := s.db.GetFeed(context.Background(), urlFeed)
	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	feedFollowRow, err := s.db.CreateFeedFollow(
		context.Background(), 
		database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID: user.ID,
			FeedID: feed.ID,
	})
	if err!=nil{
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("feed: %s, user: %s", feedFollowRow.FeedName, feedFollowRow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error{ 


	feedFollowsForUser, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	for _, feedFollowForUser := range feedFollowsForUser{
		fmt.Printf("%s\n", feedFollowForUser.FeedName)
	}

	return nil

}

func handlerUnfollow(s *state, cmd command, user database.User) error{ 

	if len(cmd.arguments) == 0{
		return fmt.Errorf("no url passed")
	}
	urlFeed := cmd.arguments[0]

	feed, err := s.db.GetFeed(context.Background(), urlFeed)
	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	err = s.db.DeleteFeedFollows(
		context.Background(), 
		database.DeleteFeedFollowsParams{
			FeedID: feed.ID, 
			UserID: user.ID,
	})

	if err!= nil{
		return fmt.Errorf("%v", err)
	}

	return nil

}
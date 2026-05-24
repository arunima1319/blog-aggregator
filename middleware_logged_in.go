package main 

import (
	"fmt"
	"context"
	"github.com/arunima1319/blog-aggregator/internal/database"

)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error{ 

	// instead of getting user from s.db.GetUser(context.Backgorund(), s.pointerConfig.CurrentUserName), we want 
	//the user to be there already.

	return func(s *state, cmd command) error { 

		user, err := s.db.GetUser(context.Background(), s.pointerConfig.CurrentUserName)
		if err!= nil{ 
			return fmt.Errorf("error in getting current user: %v", err)
		}	

		return handler(s, cmd, user)

	}
	
	
}	
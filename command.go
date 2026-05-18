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
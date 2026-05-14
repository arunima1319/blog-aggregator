package main 

import (
	"fmt"
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
	if err:=s.pointerConfig.SetUser(username); err!= nil{
		return err
	}
	fmt.Println("The user has been set")

	return nil
}
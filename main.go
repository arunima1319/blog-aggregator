package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"os"
	"database/sql"
	"github.com/arunima1319/blog-aggregator/internal/config"
	"github.com/arunima1319/blog-aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err!= nil{
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	dbQueries := database.New(db)


	currentState := state{
		db : dbQueries,
		pointerConfig: &cfg,
	}

	newCommandsSet := commands{
		commandMap : map[string]func(*state, command)error{},
	}

    newCommandsSet.register("login", handlerLogin)
	newCommandsSet.register("register", handlerRegister)
	newCommandsSet.register("reset", handlerReset)
	newCommandsSet.register("users", handlerUsers)
	newCommandsSet.register("agg", handlerAgg)
	newCommandsSet.register("addfeed", handlerAddfeed)

	argsEntered := os.Args 

	if len(argsEntered) < 2{
		fmt.Println("No argument entered")
		os.Exit(1) 
	}

	commandGiven := &command{}

	commandName := argsEntered[1]
	commandGiven.name = commandName
	
	var commandArguments []string
	if len(argsEntered) > 2{
		commandArguments = argsEntered[2:]
		
	}else{
		commandArguments = []string{}
	}
	commandGiven.arguments = commandArguments

	if err:=newCommandsSet.run(&currentState, *commandGiven); err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}

	cfg, err = config.Read()
	if err!= nil{
		fmt.Println(err)
	}else{
		fmt.Printf("db_url: %s, current_user_name: %s\n", cfg.DbURL, cfg.CurrentUserName)
	}
}
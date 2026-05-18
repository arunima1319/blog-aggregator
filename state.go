package main 

import (
	"github.com/arunima1319/blog-aggregator/internal/config"
	"github.com/arunima1319/blog-aggregator/internal/database"
	
)

type state struct{ 
	db *database.Queries
	pointerConfig *config.Config
}


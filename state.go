package main 

import "github.com/arunima1319/blog-aggregator/internal/config"

type state struct{ 
	pointerConfig *(config.Config)
}


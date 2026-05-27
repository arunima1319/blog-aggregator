# Blog Aggregator 

A command line RSS Feed aggregator built in Go

## Prerequisites 
- PostgreSQL
- Go
- Goose

## Installation 

`go install github.com/arunima1319/blog-aggregator@latest` 

## Set up

### 1. Create a database

- Start the postgres server and create a new database `gator`. 
- Run migrations using `goose -dir sql/schema postgres "postgres://username:@localhost:5432/gator" up`  

### 2. Configuration
Create a config file `.gatorconfig.json` in your home directory with the following contents: 

```
{"db_url":"postgres://username:@localhost:5432/gator?sslmode=disable","current_user_name":"username"}
```

## Usage 

| Command | Output |
| -------- | -------- |
| `blog-aggregator register <username>` | Registers a new user |
| `blog-aggregator addfeed <feed_name> <feed_url> `| Adds a feed to the database |
| `blog-aggregator agg <time_interval>` | Scrapes feeds at regular intervals and stores posts in the database
| `blog-aggregator browse` | Displays the most recent posts |


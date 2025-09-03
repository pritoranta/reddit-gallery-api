package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
)

func QuerySubredditsEndpoint(c *gin.Context) {
	var query subredditQuery
	err := c.BindQuery(&query)
	if err != nil {
		log.Printf("Error binding query %s", err)
		c.JSON(400, err.Error())
		return
	}
	log.Printf("Querying subreddits with query %v", query)
	subreddits, err := querySubreddits(query)
	if err == nil {
		c.JSON(200, subreddits)
	} else {
		c.JSON(500,  err.Error())
	}
}

type subredditQuery struct {
	SearchPhrase string			`form:"searchPhrase" binding:"required"`
	ShouldIncludeOver18 bool	`form:"shouldIncludeOver18" default:"false"`
}

type subreddit struct {
	Id string			`json:"id" validate:"required"`
	IsOver18 bool		`json:"isOver18"`
}

type subredditQueryResponse struct {
	Data struct {
		Children []struct {
			Kind string
			Data struct {
				Display_name string
				Display_name_prefixed string
				Id string
				Name string
				Public_description string
				Subreddit_type string
				Over18 bool
			}
		}
	}
	Message string
	Error int
}

func querySubreddits(query subredditQuery) (s []subreddit, e error) {
	res, err := HttpClient.Get(getSubredditQueryUrl(query))
	if err != nil {
		log.Printf("Error sending HTTP request: %s", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading HTTP response: %s", err)
		return nil, err
	}
	return getSubreddits(body)
}

func getSubredditQueryUrl(query subredditQuery) string {
	paddedSearchPhrase := url.QueryEscape(fmt.Sprintf("%3s", query.SearchPhrase)) // Reddit wants this to be at least 3 long
	over18QueryParam := "&include_over_18=off"
	if (query.ShouldIncludeOver18) {
		over18QueryParam = "&include_over_18=on"
	}
	url := "https://api.reddit.com/search.json?type=sr&q=" + paddedSearchPhrase + over18QueryParam
	return url
}

func getSubreddits(body []byte) (s []subreddit, e error) {
	var response subredditQueryResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error parsing subreddits: %s", err)
		return nil, err
	}
	if response.Message != "" {
		return nil, errors.New(response.Message)
	}
	var subreddits []subreddit
	for _, sub := range response.Data.Children {
		subreddits = append(subreddits, subreddit{
			Id: sub.Data.Display_name,
			IsOver18: sub.Data.Over18,
		})
	}
	return subreddits, nil
}

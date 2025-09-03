package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

func QueryMediaEndpoint(c *gin.Context) {
	var query mediaQuery
	err := c.BindQuery(&query)
	if err != nil {
		log.Printf("Error binding query %s", err)
		c.JSON(400, err.Error())
		return
	}
	log.Printf("Querying media with query %v", query)
	media, err := queryMedia(query)
	if err == nil {
		c.JSON(200, media)
	} else {
		c.JSON(500, err.Error())
	}
}

func queryMedia(query mediaQuery) (m mediaList, e error) {
	res, err := HttpClient.Get(getMediaQueryUrl(query))
	if err != nil {
		log.Printf("Error sending HTTP request: %s", err)
		return mediaList{nil, ""}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	fmt.Println(string(body))
	if err != nil {
		log.Printf("Error reading HTTP response: %s", err)
		return mediaList{nil, ""}, err
	}
	return getMediaList(body)
}

func getMediaQueryUrl(query mediaQuery) string {
	pageIdQueryParam := ""
	if query.PageId != "" {
		pageIdQueryParam = "&after=" + query.PageId
	}
	url := "https://api.reddit.com/r/" + query.SubredditId + "/top.json?t=all" + pageIdQueryParam
	return url
}

func getMediaList(body []byte) (s mediaList, e error) {
	var response mediaQueryResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error parsing media list: %s", err)
		return mediaList{nil, ""}, err
	}
	if response.Message != "" {
		return mediaList{nil, ""}, errors.New(response.Message)
	}
	var media []media
	for _, post := range response.Data.Children {
		if hasMedia(post.Data) {
			media = append(media, getMedia(post.Data))
		}
	}
	return mediaList{media, response.Data.After}, nil
}

func hasMedia(post postResponse) bool {
	return post.Post_hint == "image" ||
		strings.Contains(post.Post_hint, "video") ||
		(post.Post_hint == "link" && strings.Contains(post.Url, "imgur.com") && !reflect.ValueOf(post.Preview.Reddit_video_preview).IsZero())
}

func getMedia(post postResponse) media {
	if post.Secure_media.Oembed.Type == "video" {
		return media{
			PostPermalink: post.Permalink,
			SubredditId: post.Subreddit,
			PostTitle: post.Title,
			Url: post.Secure_media_embed.Media_domain_url,
			Height: post.Secure_media_embed.Height,
			Width: post.Secure_media_embed.Width,
			ThumbnailUrl: post.Thumbnail,
			IsOver18: post.Over_18,
			IsVideo: true,
		}
	} else if !reflect.ValueOf(post.Secure_media.Reddit_video).IsZero() {
		return media{
			PostPermalink: post.Permalink,
			SubredditId: post.Subreddit,
			PostTitle: post.Title,
			Url: post.Secure_media.Reddit_video.Fallback_url,
			Height: post.Secure_media.Reddit_video.Height,
			Width: post.Secure_media.Reddit_video.Width,
			ThumbnailUrl: post.Thumbnail,
			IsOver18: post.Over_18,
			IsVideo: true,
		}
	} else if !reflect.ValueOf(post.Preview.Reddit_video_preview).IsZero() {
		return media{
			PostPermalink: post.Permalink,
			SubredditId: post.Subreddit,
			PostTitle: post.Title,
			Url: post.Preview.Reddit_video_preview.Fallback_url,
			Height: post.Preview.Reddit_video_preview.Height,
			Width: post.Preview.Reddit_video_preview.Width,
			IsOver18: post.Over_18,
			IsVideo: false,
			ThumbnailUrl: "",
		}
	} else {
		return media{
			PostPermalink: post.Permalink,
			SubredditId: post.Subreddit,
			PostTitle: post.Title,
			Url: post.Url,
			Height: post.Preview.Images[0].Source.Height,
			Width: post.Preview.Images[0].Source.Width,
			IsOver18: post.Over_18,
			IsVideo: false,
			ThumbnailUrl: "",
		}
	}
}

type mediaQuery struct {
	SubredditId string	`form:"subredditId" binding:"required"`
	PageId string		`form:"pageId"`
}

type media struct {
	PostPermalink string	`json:"postPermalink"`
	SubredditId string		`json:"subredditId"`
	PostTitle string		`json:"postTitle"`
	Url string				`json:"url"`
	Height int				`json:"height"`
	Width int				`json:"width"`
	IsOver18 bool			`json:"isOver18"`
	IsVideo bool			`json:"isVideo"`
	ThumbnailUrl string		`json:"thumbnailUrl"`
}

type mediaList struct {
	Media []media		`json:"media"`
	NextPageId string	`json:"nextPageId"`
}

type mediaQueryResponse struct {
	Data struct {
		After string
		Children []struct {
			Data postResponse
		}
	}
	Message string
	Error int
}

type postResponse struct {
	Preview struct {
		Images []struct {
			Source struct {
				Url string
				Height int
				Width int
			}
		}
		Reddit_video_preview struct {
			Bitrate_kbps int
			Duration int
			Fallback_url string
			Height int
			Width int
		}
	}
	Secure_media struct {
		Type string
		Oembed struct {
			Thumbnail_url string
			Type string
			Height int
			Width int
		}
		Reddit_video struct {
			Fallback_url string
			Height int
			Width int
		}
	}
	Secure_media_embed struct {
		Media_domain_url string
		Height int
		Width int
	}
	Post_hint string
	Permalink string
	Url string
	Subreddit string
	Title string
	Thumbnail string
	Over_18 bool		
}

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type jsonResponse struct {
	Data data `json:"data"`
}
type data struct {
	Dist     int        `json:"dist"`
	Children []children `json:"children"`
}
type children struct {
	Data child_data `json:"data"`
}
type child_data struct {
	Url     string  `json:"url"`
	Id      string  `json:"id"`
	Over18  bool    `json:"over_18"`
	Preview preview `json:"preview"`
}

// Just for getting width and height of image
type preview struct {
	Images []json_image `json:"images"`
}
type json_image struct {
	Source json_image_source `json:"source"`
}
type json_image_source struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func main() {
	var subreddit string
	var minH int
	var minW int
	var mode string
	var filter_nsfw bool

	flag.StringVar(&subreddit, "subreddit", "wallpapers", "Set subreddit name")
	flag.StringVar(&mode, "mode", "random", "Set download mode (random, hot, new, rising)")
	flag.BoolVar(&filter_nsfw, "filter_nsfw", true, "Filter NSFW posts")
	flag.IntVar(&minH, "minh", 0, "Set minimum height")
	flag.IntVar(&minW, "minw", 0, "Set minimum width")
	flag.Parse()

	switch mode {
	case "random", "hot", "new", "rising":
	default:
		fmt.Fprintf(os.Stderr, "%s mode not exists\nUsage:\n", mode)
		flag.PrintDefaults()
		os.Exit(1)
	}

	imageUrl, err := get_url(subreddit, "random", filter_nsfw)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Printf(imageUrl)
}

func get_url(subreddit string, mode string, filter_nsfw bool) (string, error) {
	real_mode := mode
	if mode == "random" {
		real_mode = "hot"
	}
	url := fmt.Sprintf("https://reddit.com/r/%s/%s.json", subreddit, real_mode)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-agent", "wallpaper downloader")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get from r/%s. Error:%s\n", subreddit, err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	var jo jsonResponse
	json.NewDecoder(resp.Body).Decode(&jo)

	childs := jo.Data.Children
	if mode == "random" {
		childs = Shuffle(jo.Data.Children)
	}

	for _, child := range childs {
		if strings.HasSuffix(child.Data.Url, ".png") || strings.HasSuffix(child.Data.Url, ".jpg") {
			width := child.Data.Preview.Images[0].Source.Width
			height := child.Data.Preview.Images[0].Source.Height

			// Ignore filtered and portrait images
			if (filter_nsfw && child.Data.Over18) || width < height {
				continue
			} else {
				return child.Data.Url, nil
			}
		}
	}

	return "", errors.New("No image found")
}

func Shuffle(vals []children) []children {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]children, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

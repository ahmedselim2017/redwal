package main

import (
	"flag"
	"fmt"
	"os"
)


func main() {
	var subreddit string
	var minH int
	var minW int
	var mode string

	flag.StringVar(&subreddit, "subreddit", "wallpapers", "Set subreddit name")
	flag.StringVar(&mode, "mode", "random", "Set download mode (random, hot, new, rising)")
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

}

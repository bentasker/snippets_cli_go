package main

import "fmt"
import "github.com/mmcdole/gofeed"

func main() {
    fp := gofeed.NewParser()
    fp.UserAgent = "snippets_cli_go"
    feed, _ := fp.ParseURL("https://www.bentasker.co.uk/rss.xml")
    fmt.Println(feed.Title)
}

package main

import "fmt"
import "github.com/mmcdole/gofeed"

func main() {
    fp := gofeed.NewParser()
    feed, _ := fp.ParseURL("https://www.bentasker.co.uk/rss.xml")
    fmt.Println(feed.Title)
}

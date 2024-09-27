package main

import "fmt"
import "github.com/mmcdole/gofeed"


func fetchFeed() (feed *gofeed.Feed, err error){
    // TODO: caching
    fp := gofeed.NewParser()
    fp.UserAgent = "snippets_cli_go"
    return fp.ParseURL("https://www.bentasker.co.uk/rss.xml")
}


func main() {
    
    feed, _ := fetchFeed()
    fmt.Println(feed.Title)
    
    for _, item := range feed.Items{
        fmt.Println(item.Title)
        fmt.Println(item.Link)
        //fmt.Println(item.Description)
        for _, cat := range item.Categories{
            fmt.Println(cat)
        }
    }
}

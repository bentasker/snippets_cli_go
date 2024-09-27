package main

import (
    "fmt"
    "os"
    
    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/mmcdole/gofeed"
)

type searchResult struct{
    id          int
    title       string
    language    string // TODO - don't know if we'll be able to reliably populate this
}

type search struct{
    term    string
}

func fetchFeed() (feed *gofeed.Feed, err error){
    // TODO: caching
    fp := gofeed.NewParser()
    fp.UserAgent = "snippets_cli_go"
    return fp.ParseURL("https://www.bentasker.co.uk/rss.xml")
}


func printTable(res []searchResult, s search) {
    
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetTitle(fmt.Sprintf("Search results: %s", s.term))
    t.AppendHeader(table.Row{"#", "Title", "Language"})
    
    for _, r := range res{
        t.AppendRow([]interface{}{r.id, r.title, r.language})
    }
    t.Render()
}


/** Iterate through the feed and apply the desired search term
 * 
 */
func searchFeed(feed *gofeed.Feed, search search) []searchResult{
    results := []searchResult{}
    
    for _, item := range feed.Items{
        // fmt.Println(item.Title)
        // fmt.Println(item.Link)
        //fmt.Println(item.Description)
        //for _, cat := range item.Categories{
        //    fmt.Println(cat)
        //}
        
        var res searchResult
        res.id = 1 // TODO
        res.title = item.Title
        res.language = "N/A"
        results = append(results, res)
    }
    return results
}

func main() {
    
    feed, _ := fetchFeed()
    fmt.Println(feed.Title)
    
    var search search
    search.term = "foo"

    results := searchFeed(feed, search)
    
    // Render the results
    printTable(results, search)
}

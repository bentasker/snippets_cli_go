package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    
    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/mmcdole/gofeed"
)

type searchResult struct{
    id          int
    title       string
    language    string // TODO - don't know if we'll be able to reliably populate this
    matchtype   string
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
    var idMatchMode bool
    var searchID int
    var err error
    results := []searchResult{}
    searchTerm := strings.ToLower(search.term)
    // IDs decrement as we iterate through, so get the number
    // of items
    id := feed.Len()
    
    // Have we been passed something that's simply a number?
    if searchID, err = strconv.Atoi(searchTerm); err == nil {
        // We have, switch to ID match mode
        idMatchMode = true
    }
    
    for _, item := range feed.Items{
        // fmt.Println(item.Title)
        // fmt.Println(item.Link)
        //fmt.Println(item.Description)
        //for _, cat := range item.Categories{
        //    fmt.Println(cat)
        //}
        
        var matched bool
        var res searchResult
        res.id = id
        res.title = item.Title
        res.language = "N/A" // TODO

        if !idMatchMode {
            // Does the title match?
            if strings.Contains(strings.ToLower(item.Title), searchTerm) {
                res.matchtype = "title"
                matched = true
            }
            
            // What about keywords?
            for _, cat := range item.Categories{
                if strings.Contains(strings.ToLower(cat), searchTerm){
                    res.matchtype = "keyword"
                    matched = true
                    break;
                }
            }
            
            if matched {
                results = append(results, res)
            }
        }else{
            if id == searchID {
                // TODO - print the snippet
                fmt.Println("%i matches", id)
            }
        }
    
        id -= 1
    }
    return results
}

func main() {
    
    feed, _ := fetchFeed()
    fmt.Println(feed.Title)
    
    var search search
    search.term = "10"

    results := searchFeed(feed, search)
    
    // Render the results
    printTable(results, search)
}

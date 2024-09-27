package main

import (
    "bytes"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    
    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/mmcdole/gofeed"
    "github.com/PuerkitoBio/goquery"
    "golang.org/x/net/html"
    md "github.com/JohannesKaufmann/html-to-markdown"
    "github.com/JohannesKaufmann/html-to-markdown/plugin"
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


/** Fetch a snippet and process it into something which can be
 * output to console
 * 
 */
func printSnippet(id int, title string, link string){
    // Fetch the page
    fmt.Println(link)
    //resp, err := http.Get(link)
    resp, err := http.Get("https://snippets.bentasker.co.uk/page-2409261238-List-Resource-Requests-and-Limits-for-Kubernetes-pods-Misc.html")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        log.Fatalf("Failed with Status code: %d %s", resp.StatusCode)
    }
    
    // Parse the page
    doc, err := html.Parse(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    
    
    // Define the map that we'll write content into
    entry := make(map[string]*html.Node)
    
    // Iterate through to select the items we need
    var f func(*html.Node)
    f = func(n *html.Node){
        // Check whether it's a div
        if n.Type == html.ElementNode && n.Data == "div" {
            attrs := n.Attr
            for _, attr := range attrs {
                    // TODO: match on ID rather than class
                    if attr.Key == "id" {
                        if strings.Contains(attr.Val, "pageContent"){
                            entry["body"] = n
                        }
                    }
                }
        }
    
        // Recurse through children
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }    
        
    f(doc)
    
    // We don't want link targets to be included - we're looking to try and 
    // make something for a user to read in a CLI
    links := md.Rule {
        Filter: []string{"a"},
        Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
            return md.String(content)
        },
    }
    
    images := md.Rule {
        Filter: []string{"img"},
        Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
            return md.String("")
        },
    }
    
    
    
    _, ok := entry["body"]; if ok {
        var h bytes.Buffer
        html.Render(&h, entry["body"])
        converter := md.NewConverter("", true, nil)
        converter.Use(plugin.GitHubFlavored())
        converter.AddRules(links)
        converter.AddRules(images)
        markdown, _ := converter.ConvertString(h.String())
        fmt.Println(markdown)
        
    }

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
                printSnippet(id, item.Title, item.Link)
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

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

/** Fetch and parse the RSS feed
 * 
 */
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
    // we use a map because I may later want to extract
    // some additional info
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
    
    // Trigger iteration
    f(doc)
    
    // Render
    _, ok := entry["body"]; if ok {
        renderSnippet(entry["body"], id, title, link)
    }else{
        log.Fatal("Unable to retrieve snippet")
    }
}    
 
/** Generate console output for a snippet
 * 
 */
func renderSnippet(snippet *html.Node, id int, title string, link string){
    
    // Define some rules to generate more readable output
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
    
    // Call html-to-markdown
    var h bytes.Buffer
    html.Render(&h, snippet)
    converter := md.NewConverter("", true, nil)
    converter.Use(plugin.GitHubFlavored())
    
    // Add the rules
    converter.AddRules(links)
    converter.AddRules(images)
    
    // Render
    markdown, _ := converter.ConvertString(h.String())
    
    // Print
    fmt.Println(fmt.Sprintf("%d: %s\n", id, title))
    fmt.Println(markdown)
    fmt.Println("\n### HTML Link\n")
    fmt.Println(link)
    fmt.Println("")
}

/** Output a tabulated set of results
 * 
 */
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
                // print the snippet
                printSnippet(id, item.Title, item.Link)
                return results
            }
        }
    
        id -= 1
    }
    return results
}

func main() {
    var search search
    
    if len(os.Args[1:]) < 1 {
        log.Fatal("No search term")
    }

    // Take search terms from the command line
    search.term = strings.Join(os.Args[1:], " ")
    
    // Fetch the feed
    feed, err := fetchFeed(); if err != nil {
        log.Fatal(err)
    }
    
    // Run the search
    // note: if an ID was provided this function
    // will instead trigger printing of the 
    // snippet
    results := searchFeed(feed, search)
    
    // Render the results if any were returned
    if len(results) > 0 {
        printTable(results, search)
    }
}

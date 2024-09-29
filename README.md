# snippets_cli

A utility to search snippets.bentasker.co.uk (and other sites) from the command line.

It's a little poorly named because I named it based on my original intent and then built something more genericised.

Basically, it provides CLI based search for sites with RSS feeds:

* Fetch the feed
* Iterate through trying to match the search term against title or keywords
* Print either a results table or the associated page

This is the successor to my [original python script](https://github.com/bentasker/snippets_cli), rewritten in anticipation of a change in static site generator (which went live in September 2024).

Basically a tool for my own convenience, allows me to grab and search snippets without leaving the comfort of my terminal.


---

### Usage

```sh
snippets_cli [search term | id]
```

The CLI can be passed a search term or an integer taken from the results table, in which case the snippet will be printed to console.

If only a single snippet matches a search term, no results table will be displayed and the relevant snippet will be printed instead.

---

### Multi-Site

The utility supports searching a number of my sites, based on the name of the binary when it's called.

It's therefore possible to do

```sh
ln -s snippets_cli /usr/bin/btcli
ln -s snippets_cli /usr/bin/rbt_cli
```

The first will search `www.bentasker.co.uk` whilst the second will search `recipebook.bentasker.co.uk`. 

See [`utilities/snippets_cli_go#8`](https://projects.bentasker.co.uk/gils_projects/issue/utilities/snippets_cli_go/8.html) for more information.

---

### Example Output

```text
$ ./snippets_cli elemen
+---------------------------------------------------------------------------------------------------+
| Search results: elemen                                                                            |
+-----+--------------------------------------------------------------------------------+------------+
|   # | TITLE                                                                          | LANGUAGE   |
+-----+--------------------------------------------------------------------------------+------------+
| 114 | Detect when Enter is pressed within an Input Element (Javascript)              | Javascript |
| 113 | Select contents of element and its children and copy to clipboard (Javascript) | Javascript |
|  12 | Check if table has element (LUA)                                               | LUA        |
+-----+--------------------------------------------------------------------------------+------------+
```


---

### Copyright

Copyright (C) 2024 [B Tasker](https://www.bentasker.co.uk/). All Rights Reserved. Released under [MIT LICENSE](https://www.bentasker.co.uk/pages/licenses/mit-license.html)

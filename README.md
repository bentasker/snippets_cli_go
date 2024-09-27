# snippets_cli

A utility to search snippets.bentasker.co.uk from the command line.

This is the successor to my [original python script](https://github.com/bentasker/snippets_cli), rewritten in anticipation of a change in static site generator.

Basically a tool for my own convenience, allows me to grab and search snippets without leaving the comfort of my terminal.


---

### Usage

```sh
snippets_cli [search term]
```

The search term may be an integer taken from the results table, in which case the snippet will be printed to console.

If only a single snippet matches a search term, no results table will be displayed: the snippet will be printed instead.

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

Copyright (C) 2024 B Tasker. All Rights Reserved. Released under the GNU GPL V2 License, see LICENSE.

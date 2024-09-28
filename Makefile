build:
	go build -ldflags "-w -s" .

install: build
	cp snippets_cli /usr/local/bin
	ln -s /usr/local/bin/snippets_cli /usr/local/bin/rbt_cli
	ln -s /usr/local/bin/snippets_cli /usr/local/bin/btcli


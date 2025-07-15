GOROOT=$(shell go env GOROOT)

default: SpaceDebris web

SpaceDebris:
	go build -o $@ .

webserver: web
	    python3 -m http.server --directory ./web

web:
	rm -rf web/
	mkdir -p web/
	env GOOS=js GOARCH=wasm go build -o web/game.wasm .
	cp $(GOROOT)/lib/wasm/wasm_exec.js web/
	cp index.html web/

publish: web
	rsync -a web/ kaka:/var/www/fortyfootgames.duckdns.org/SpaceDebris/

clean:
	rm -rf web
	rm -f SpaceDebris

.PHONY: web clean webserver default publish

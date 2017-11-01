export WATCHFOLDER=.
export WORKFLOWURL=http://djb.local/api/inotify/%%PATH%%
export FILESUFFIX=.xml

this=$(shell basename `pwd`)

phoney: run

run:
	env|egrep "WATCH|WORK|FILE"
	go fmt main.go
	http_proxy='' https_proxy='' go run main.go

build:
	go fmt *.go
	go build .
	strip $(this)
	ls -lh $(this)

exec:
	http_proxy='' https_proxy='' ./$(this)

test:
	touch fred.xml

CE=podman run -it --security-opt label=disable
GO=${CE} --rm -e 'CGO_ENABLED=1' -v './:/usr/src/myapp' -w '/usr/src/myapp' docker.io/library/golang:1.19-alpine go

build:
	podman build -t dlv:1.19-alpine .

debug:
	${CE} --rm \
		-v './:/usr/src/myapp' \
		-w '/usr/src/myapp' \
		localhost/dlv:1.19-alpine \
		dlv debug vimm.go

init:
	${GO} mod init oxylabs.io/web-scraping-with-go
	${GO} get github.com/gocolly/colly

run:
	${GO} run vimm.go

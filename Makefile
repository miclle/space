all: gofmt govet golangci-lint golint test install

check: gofmt govet golangci-lint golint test

gofmt:
	gofmt -s -l . | tee .gofmt.log
	test `cat .gofmt.log | wc -l` -eq 0
	rm .gofmt.log

govet:
	go vet ./...

golangci-lint:
	golangci-lint run ./...

golint:
	golint ./... | tee .golint.log
	test `cat .golint.log | wc -l` -eq 0
	rm .golint.log

test:
	go clean -i -testcache
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

coverage: test
	go tool cover -html=coverage.txt -o coverage.html
	open coverage.html

install:
	go install golang.org/x/lint/golint@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cespare/reflex@latest
	go install github.com/hhatto/gocloc/cmd/gocloc@latest
	go install -ldflags '-s -w' ./...
	cd ui; yarn install

# Run the website
site:
	cd ui; yarn start

# Run back-end service and proxy website requests
serve:
	PROXY_WEBSITE=http://localhost:3000 reflex\
		-d none -s\
		-R 'tmp/' \
		-R '\.github' \
		-R 'ui/src/' \
		-R 'ui/build/' \
		-R '/node_modules' \
		-R 'cmd/space/public/' \
		-R '^coverage' \
		-R 'Makefile' \
		-R '.log$$' \
		-R '_test.go$$'\
		-- go run -trimpath cmd/space/*.go -f cmd/space/config.local.yaml | tee -a development.log

# Build website and backend to a single app
# 1. build website
# 2. move build website result to backend embed path: public
# 3. build backend
# 4. clean
build:
	cd ui; yarn install && yarn build
	cd cmd/space; CGO_ENABLED=0 go build -tags=jsoniter -trimpath -ldflags '-s -w' ./...

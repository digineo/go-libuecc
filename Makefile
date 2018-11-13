gotest = $(GOPATH)/bin/gotest

test: $(gotest)
	$(MAKE) -C testdata
	$(gotest)

$(gotest):
	go get -u github.com/rakyll/gotest

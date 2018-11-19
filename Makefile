gotest = $(GOPATH)/bin/gotest

test: $(gotest)
	$(MAKE) -C testdata
	$(gotest) -v

$(gotest):
	go get -u github.com/rakyll/gotest

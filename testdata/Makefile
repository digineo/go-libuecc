CFLAGS = -Wall -Werror -pedantic -std=c99
LDFLAGS = -I./libuecc/include  -luecc

.PHONY: data
data: gen
	mkdir -p ./cases
	./gen

gen: gen.c
	gcc $(LDFLAGS) $(CFLAGS) -o $@ $<

clean:
	rm -f ./gen
	rm -rf ./cases

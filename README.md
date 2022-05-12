Optionals?
==========

Another optionals package in Go.

### Installation

```shell
go get -u github.com/SerialImprovement/go-optionals
```

### What is different in this implementation?

I use `channels` & `select` to attempt to replicate pattern matching in Go.

So to get a value from an optional:

```Go
package main

import (
	"github.com/SerialImprovement/go-optionals/pkg/option"
	"log"
)

func main() {
	var opt = option.Some("World")

	some, none := opt.Get()

	select {
	case v := <-some:
		log.Printf("Hello, %s!", v)
	case <-none:
		log.Print("oh dear i should do something else")
	}
}
```

### Why?

It's cool to play with languages! Learning features through implementation
is a great way to learn, at least in my opinion.

Why is it any better than just `nil` checks? That was exactly my dilemma:

```Go
package main

import "log"

func main() {
	var opt *string

	if opt != nil {
		log.Printf("Hello, %s!", *opt)
	} else {
		log.Print("oh dear i should do something else")
	}
}
```



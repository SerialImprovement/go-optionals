Optionals?
==========

Another optionals package in Go.

### Installation

```shell
go get -u github.com/SerialImprovement/go-optionals
```

Next check out the [examples](/examples).

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

Though its never that simple, is it? In Go we're always dealing with at least 2 things, `values` and `errors`
And more commonly we're dealing with 3 things `values`, `nils` & `errors`.

Take this function signature:

```Go
func FetchUser(id int) (*User, error) {}
```

How many possibilities are there that we need to handle?

```Go
maybeUser, err := FetchUser(5)

// typical err check
if err != nil {
	// hopefully do something useful
	return
}

// now we assume user isnt nil
fmt.Printf("%+v", maybeUser)
```

But can we really assume `maybeUser` isn't nil? this is only by convention, you cannot 
assume that `maybeUser` will be a real value.

So you might think; "Duh, just use another if", lets look at that:

```Go
maybeUser, err := FetchUser(5)

// typical err check
if err != nil {
	// hopefully do something useful
	return
}

// now we assume user isnt nil
if maybeUser != nil {
	fmt.Printf("%+v", maybeUser)
} else {
	// hopefully do something useful
}
```



How does it look with a typical optional:

```Go
func FetchUser(id int) (optional.Option[*User], error)
```

Ok still looks the same... what about the usage:

```Go
maybeUser, err := FetchUser(5)

// typical err check
if err != nil {
	// hopefully do something useful
	return
}

// now we use optionals to guarantee we either get a valid value or not

some, none := maybeUser.Get()

select {
	case v := <-some:
		log.Printf("Hello, %s!", v)
	case <-none:
		log.Print("oh dear i should do something else")
	}
```

Ok so that removes the convention problem. But we've still got a lot of different control
flow happening.

We have a `Some`, a `None` and an `error`. Lets modify a typical optional and create a 
new thing called a SoNoEr (**So**me**No**ne**Er**r).

```Go
// we can now call .Get() directly because out function returns only a SoNoEr
// Note: we also can drop the name 'maybeUser' because we can know for sure it is a user.
some, none, err := FetchUser(5).Get()

select {
	case user := <-some:
		log.Printf("Hello, %s!", user.Name)
	case <-none:
		log.Print("oh dear i should do something else")
	case e := <-err:
		log.Print("an error happened! %s", e)
	}
```

As you can see we've consolidated the control structures into one.

### Is this good?

I don't know! I'll be using this method in some real code. In the event that it is 
found to be good i'll probably remove all this nonsense and treat this like it was 
always the correct way.

### Motivation

Sometimes it's good to try fit a square peg into a round hole. Maybe we don't know what 
peg we're holding or what types of holes might exist!

Also...

It's cool to play with languages! Learning features through implementation
is a great way to learn, at least in my opinion.

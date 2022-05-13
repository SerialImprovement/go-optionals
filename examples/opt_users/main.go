package main

import (
	"errors"
	"github.com/SerialImprovement/go-optionals/pkg/option"
	"log"
)

var (
	users = []*Asgardian{
		{
			Id:   1,
			Name: "Odin",
		},
		{
			Id:   2,
			Name: "Thor",
		},
		{
			Id:   3,
			Name: "Hathor",
		},
	}
)

type Asgardian struct {
	Id   int
	Name string
}

func main() {
	HandlePossibleAsgardian(FetchAsgardian(1))
	HandlePossibleAsgardian(FetchAsgardian(4))
	HandlePossibleAsgardian(FetchAsgardian(7))
}

func HandlePossibleAsgardian(possible option.Option[*Asgardian], err error) {
	some, none := possible.Get()

	if err != nil {
		log.Printf("error: %s", err)
		return
	}

	select {
	case asgardian := <-some:
		log.Printf("we got: %s", asgardian.Name)
	case <-none:
		log.Printf("couldnt find that asgardian")
	}
}

func FetchAsgardian(id int) (option.Option[*Asgardian], error) {
	// for the sake of demo, assume that this id causes an error.
	if id == 4 {
		return option.None[*Asgardian](), errors.New("loki is a secret")
	}

	for _, u := range users {
		if u.Id == id {
			return option.Some(u), nil
		}
	}

	return option.None[*Asgardian](), nil
}

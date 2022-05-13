package main

import (
	"errors"
	"github.com/SerialImprovement/go-optionals/pkg/sne"
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

func HandlePossibleAsgardian(possible sne.SoNoEr[*Asgardian]) {
	some, none, err := possible.Get()

	select {
	case asgardian := <-some:
		log.Printf("we got: %s", asgardian.Name)
	case <-none:
		log.Printf("couldnt find that asgardian")
	case e := <-err:
		log.Printf("error: %s", e)
	}
}

func FetchAsgardian(id int) sne.SoNoEr[*Asgardian] {
	// for the sake of demo, assume that this id causes an error.
	if id == 4 {
		return sne.Err[*Asgardian](errors.New("loki is a secret"))
	}

	for _, u := range users {
		if u.Id == id {
			return sne.Some(u)
		}
	}

	return sne.None[*Asgardian]()
}

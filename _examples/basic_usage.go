package main

import (
	"fmt"
	"strings"

	"github.com/b1tkeeper/goset"
)

type Person struct {
	Name    string
	Hobbies []string
}

// Hash implemented from goset.Hashable, that makes your
// custom type can be stored in set
// Note that you always should return some string correspond
// the field represent your object consistency, it will affect
// set's comparison operations on objects
func (p Person) Hash() string {
	return p.Name + strings.Join(p.Hobbies, "")
}

func main() {
	set := goset.NewSet()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Remove(2)
	fmt.Println(set)
	fmt.Println(set.Size())

	set1 := goset.NewSet(1, 2, 3)
	set2 := goset.NewSet(4, 5, 6)
	fmt.Println(set1.Intersect(set2))
	fmt.Println(set1.Union(set2))
	fmt.Println(set1.SymmetricDifference(set2))
	fmt.Println(set1.Difference(set2))

	set3 := goset.NewSet()
	set3.Add(Person{Name: "James", Hobbies: []string{"basketball", "swiming"}})
	set3.Add(Person{Name: "Briant", Hobbies: []string{"basketball"}})
	set3.Add(Person{Name: "James", Hobbies: []string{"basketball", "swiming"}})
	fmt.Println(set3) // goset.ThreadUnsafeSet{ {James [basketball swiming]}, {Briant [basketball]} }

	set4 := goset.NewThreadUnsafeSet()
	fmt.Println(set4.Contains(1))
}

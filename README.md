# goset
A Fast、Thread-Safe、Portable set that can store any object in Golang

## Principles
1. Store Any Object: Due to the incomplete function of Golang's generics, the commonly used set library only supports storing the basic types in Golang (reference to comparable type constrait in Golang), but goset will support any objects of any custom type.
2. Portable: goset ditches weak generic function of golang and any keyword, to downward compatible with more versions (Goals: go 1.10 Windows XP)
3. Fast: No particular performance loss

## Usage

### Basic Usage

```go
set := goset.NewSet()
set.Add(1)
set.Add(2)
set.Add(3)
set.Remove(2)
fmt.Println(set)
fmt.Println(set.Size())
```
### Set Operations
```go
set1 := goset.NewSet(1, 2, 3)
set2 := goset.NewSet(4, 5, 6)
fmt.Println(set1.Intersect(set2))
fmt.Println(set1.Union(set2))
fmt.Println(set1.SymmetricDifference(set2))
fmt.Println(set1.Difference(set2))
```

### Store Custom Type
```go
// Store Custom Type
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

set3 := goset.NewSet()
	set3.Add(Person{Name: "James", Hobbies: []string{"basketball", "swiming"}})
	set3.Add(Person{Name: "Briant", Hobbies: []string{"basketball"}})
	set3.Add(Person{Name: "James", Hobbies: []string{"basketball", "swiming"}})
	fmt.Println(set3) // goset.ThreadUnsafeSet{ {James [basketball swiming]}, {Briant [basketball]} }
```

### Unsafe Set

```go
set4 := goset.NewThreadUnsafeSet()
fmt.Println(set4.Contains(1))
```

## Methods List
- `Add(val interface{}) bool`
- `Cardinality() int`
- `Size() int`
- `Clear()`
- `Clone() Set`
- `Contains(val ...interface{}) bool`
- `Difference(other Set) Set`
- `Equal(other Set) bool`
- `Intersect(other Set) Set`
- `IsProperSubset(other Set) bool`
- `IsProperSuperset(other Set) bool`
- `IsSubset(other Set) bool`
- `IsSuperset(other Set) bool`
- `Each(func(elem interface{}) bool)`
- `Iter() <-chan interface{}`
- `Iterator() *Iterator`
- `Remove(i interface{})`
- `String() string`
- `SymmetricDifference(other Set) Set`
- `Union(other Set) Set`
- `Pop() (interface{}, bool)`
- `ToSlice() []interface{}`
- `MarshalJSON() ([]byte, error)`
- `UnmarshalJSON(b []byte) error`
# go-page

go-page is a collection of functions that makes it easy to store/retrieve data to/from files.

## Supported formats
+ [X] JSON
+ [X] YAML
+ [X] XML

## Example
```golang
import "github.com/henrikac/go-page"

type Group struct {
	Name    string
	Members []Person
}

type Person struct {
	Name string
	Age  int32
}

group := Group{
	Name: "TheGophers",
	Members: []Person{
		{Name: "Alice", Age: 19},
		{Name: "Bob", Age: 19},
	},
}

err := page.WriteJSON("data.json", group)
if err != nil {
	// handle error
}

var readGroup Group

err = page.ReadJSON("data.json", &readGroup)
if err != nil {
	// handle error
}
```

Colprint
========

Easy column formatted printing of structs and slices/arrays in golang.

Colprint is a small go package to help build CLI appliactions where you want to list items in 
human readable form in formatted columns. Colprint builds on [Columnize](https://github.com/ryanuber/columnize) and adds functionality to easy print structs and 
slices/arrays of structs.

Example:
```go
package main

import "github.com/peteabre/colprint"

type Person struct {
        FirstName string `column:"First name,1"`
        LastName string  `column:"Last name,2"`
        Age int          `column:"Age,3"`
} 

func main()  {
        persons := []Person{
                {
                        FirstName: "Ola",
                        LastName:  "Nordmann",
                        Age:        35,
                },
                {
                        FirstName: "Kari",
                        LastName:  "Nordmann",
                        Age:        37,
                },
         }
         colprint.Print(persons)
}
```

As you can see you can simply pass a slice/array and the result will be:

```
First name  Last name  Age
Ola         Nordmann   35
Kari        Nordmann   37

```

Limitations
===========
Supports only fields of primitive types.
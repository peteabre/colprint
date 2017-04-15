Colprint
========

Easy column formatted printing of structs and slices/arrays of structs in golang.

[![Build Status](https://travis-ci.org/peteabre/colprint.svg)](https://travis-ci.org/peteabre/colprint)
[![GoDoc](https://godoc.org/github.com/peteabre/colprint?status.svg)](https://godoc.org/github.com/peteabre/colprint)

Colprint is a small Go package to help build CLI appliactions where you want to list items in 
human readable form in formatted columns. Colprint builds on [Columnize](https://github.com/ryanuber/columnize), and adds functionality to easy print structs and 
slices/arrays of structs. You just have have to add the colprint tag to the fields you want to print.

Small example:
```go
package main

import "github.com/peteabre/colprint"

type Person struct {
        FirstName string `colprint:"First name,1"`
        LastName string  `colprint:"Last name,2"`
        Age int          `colprint:"Age,3"`
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

As you can see, if you have a tagged struct, you can simply pass a slice/array and the result will be:

```
First name  Last name  Age
Ola         Nordmann   35
Kari        Nordmann   37
```

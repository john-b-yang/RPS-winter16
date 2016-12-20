package main

import "flag"
import "fmt"

func main() {
  //Foo and 42 are default values
  wordPtr := flag.String("word", "foo", "String parameter")
  numbPtr := flag.Int("numb", 42, "Int Parameter")

  flag.Parse()

  parameter := *wordPtr
  if parameter == "hi" {
    fmt.Println("If - Conditional 1")
  } else {
    fmt.Println("Else - Conditional 2")
  }

  switch parameter {
  case "hi":
    fmt.Println("Case 1")
  case "bye":
    fmt.Println("Case 2")
  default:
    fmt.Println("Didn't work :(")
  }
}

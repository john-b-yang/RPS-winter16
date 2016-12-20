func main() {
  //Declaration
  var x int = 10 //type is explicitly specified
  x := 10 //type is implied
  const y = 10 //Cannot use := for constant

  //Iteration and Conditionals
  for i := 0; i < 10; i++ {
       fmt.Println(“This is a for loop”)
  }

  for i < 10 {
       fmt.Println(“While loops do not require the init and post conditions of a for loop”)
  }

  for {
      fmt.Println(“The beautiful but deadly infinite loop. As we can see, all loops in Go are based on manipulating the conditions of a for loop”)
  }

  if x * 5 == 50 {
      fmt.Println(“Look mom! No parentheses!”)
  }

  value := 10
  switch (optional value, true by default) {
       case “example1”:
            fmt.Println(“This is case one”)
            fallthrough //Continues on to the next statement
       default:
            defer fmt.Println(“Arguments are evaluated immediately, but call is not executed until return.”)
            defer fmt.Println(“The calls are executed as a stack (LIFO), so this statement will print before the previous defer”)
            fmt.Println(“This is the default”)
  }
}

# Flow

![Generated image of a vibrant jellyfish floating in the ocean.](docs/jellyfish.webp)

> TODO think about logic concerning joining sources<br>
> TODO think about multi threading vs single threaded async of RxJS <br>
> TODO think about implementing nullable types like Kotlin <br>
> TODO think about design by contract like Eiffel, could be implemented in the typing perhaps? <br>
> TODO think about native utility methods like map, filter, reduce, should be chainable en optimized(no multiple loops in execution with two maps or something) <br>
> TODO think about a flexible generics system <br>
> TODO think about arrays, they should be dynamic and implement the iterator interface so you can abstract them<br>
> TODO should really have function overloading, miss this a lot in Go

Reactive programming language experimental outlines.

Hello World in Flow; prints "Hello World" to the terminal.

```
package main


const main = () => {
    "hello world"
        => capitalize
        ~> print
}
```

Example Programs
=

> For testing the  compiler and drawing thoughts

```
package main            // Main package for compiler

const main = () => {    // Main entrypoint for compiler
    let a = 0           // Can reassign new values to a
    const b = a + 2     // Can't reassign b to new value, value is not fixed
    
    b
        => multiply(2)  // curried fn multiply
        ~> print        // prints twice immediately and after reassignment of a (2, 4)
    
    a = 2               // Reasigns new value to a, triggers value change emit in b
}
```

Symbols & reserved keywords
=

## Symbols

| Symbol |                Name                 |                              Meaning                              |                              Context |
| ------ | :---------------------------------: | :---------------------------------------------------------------: | -----------------------------------: |
| `=>`   |            Pipe Operator            |          Pipes values from source to following function           |                    Following a value |
| `=>`   |           Function Sybmol           |          Declares a function, e.g. `(s string) => print`          | Following a argument or return list* |
| `~>`   |         Subscribe Operator          | Subscribes on source, makes values flow into subscribing function |                                      |
| `~`    |         Source Type Symbol          |                Declares a source type, e.g. `~int`                |                     Type declaration |
| `=`    |        Assignment Opterator         |               Assigns value to variable or constant               |                                      |
| `==`   |         Equality Opterator          |                 Compares two values for equality                  |                                      |
| `!=`   |       Non-equality Opterator        |               Compares two values for non-equality                |                                      |
| `*`    |          Pointer Opterator          |                  Points to the value of pointer                   |                     Type declaration |
| `&`    |          Address Opterator          |               Takes address of variable or constant               |                     Type declaration |
| `?`    |         Optional Opterator          |              Declares argument or field as optional               |                                      |
| `_`    |          Blank identifier           |                    Declares given value unused                    |                                      |
| `(`    |    Argument List Open Opterator     |                  Opens argument list of function                  |                                      |
| `(`    |     Return List Open Opterator      |                   Opens return list of function                   |                                      |
| `)`    |    Argument List Close Opterator    |                 Closes argument list of function                  |   Following a function argument list |
| `)`    |     Return List Close Opterator     |                  Closes return list of function                   |     Following a function return list |
| `+`    |          Math Add Operator          |                          Adds two values                          |                   In numeric context |
| `+`    |    String Concatenation Operator    |                        Concats two strings                        |                    In string context |
| `-`    |       Math Substract Operator       |                       Substracts two values                       |                   In numeric context |
| `*`    |       Math Multiply Operator        |                       Multiplies two values                       |                   In numeric context |
| `/`    |       Math Division Operator        |                        Divides two values                         |                   In numeric context |
| `+=`   |    Math Add Assignment Operator     |      Assigns result of added values to variable or constant       |                   In numeric context |
| `-=`   | Math Substract Assignment Operator  |   Assigns result of substracted values to variable or constant    |                   In numeric context |
| `*=`   |  Math Multiply Assignment Operator  |    Assigns result of multiplied values to variable or constant    |                   In numeric context |
| `/=`   |  Math Division Assignment Operator  |     Assigns result of divided values to variable or constant      |                   In numeric context |
| `%=`   |   Math Modulo Assignment Operator   |           Assign modulo result to variable or constant            |                   In numeric context |
| `{`    |     Function Body Open Operator     |                        Opens function body                        |          Following a function symbol |
| `}`    |    Function List Close Opterator    |                       Closes function body                        |            Following a function body |
| `}`    |    Function List Close Opterator    |                       Closes function body                        |            Following a function body |
| `<`    | String Interpolation Open Operator  |     Starts block for string interpolation e.g. `"Value <x>"`      |                Inside string literal |
| `>`    | String Interpolation Close Operator |                Ends block for string interpolation                |                Inside string literal |
| `\`    |       String Escape Character       |              Escapes characters in string e.g. "\\<"              |                Inside string literal |
| `;`    |        Termination Operator         |                 Terminates statement declaration                  |                                      |

> TODO: interface and struct body, emphasises in logical clauses, logical and or xor, array construct and indexing...

## Reserved Keywords

### Types

### General Keywords

| Keyword     |           Name           |                              Meaning                              |                          Context |
| ----------- | :----------------------: | :---------------------------------------------------------------: | -------------------------------: |
| `const`     |   Constant Declaration   |                      Declares a new constant                      |                                  |
| `let`       |   Variable Declaration   |                      Declares a new variable                      |                                  |
| `package`   |   Package Declaration    |                      Declares a new package                       |                                  |
| `import`    |    Import Declaration    |              Imports specified content of a package               |                                  |
| `export`    |    Export Declaration    | Exports following named value; e.g. `export const alfa = "Hello"` |                                  |
| `return`    |     Return Operation     |                  Returns following values; e.g.                   |                                  |
| `for`       |     Loop Declaration     |                      Declares loop construct                      |                                  |
| `if`        |      If Declaration      |                       Declares if construct                       |                                  |
| `else`      |     Else Declaration     |                      Declares else construct                      | Following a if or elif construct |
| `elif`      |   Else If Declaration    |                    Declares else if construct                     |         Following a if construct |
| `switch`    |    Switch Declaration    |                     Declares switch construct                     |                                  |
| `case`      |     Case Declaration     |                           Declares case                           |        Inside a switch construct |
| `default`   |   Default Declaration    |                      Declares default cases                       |        Inside a switch construct |
| `break`     |    Break Declaration     |                        Breaks out of loop                         |                    Inside a loop |
| `continue`  |   Continue Declaration   |                   Skips current loop iteration                    |                    Inside a loop |
| `async`     | Asynchronous Declaration |            Starts asynchronous function in new thread             |                                  |
| `yield`     |    Yield Declaration     |                    Yields value from function                     |             Inside function body |
| `struct`    |    Struct Declaration    |                       Declares a new struct                       |                                  |
| `interface` |  Interface Declaration   |                     Declares a new interface                      |                                  |
| `type`      |     Type Declaration     |                        Declares a new type                        |                                  |
| `count`     |       Count Value        |               Holds the current count in a pipeline               |                    In a pipeline |
| `throw`     |    Throw Declaration     |                          Throws a error                           |                                  |

### Reactive Keywords

| Keyword        |         Name          |                            Meaning                             |                 Context |
| -------------- | :-------------------: | :------------------------------------------------------------: | ----------------------: |
| `error`        |    Error Operator     |                 Adds a error handling function                 |                   Error |
| `closed`       |    Closed Operator    |    Adds a function to be executed when the source is closed    |                  Closed |
| `split`        |    Split Operator     |             Splits source values in seperate emits             | Array and string source |
| `reduce`       |    Reduce Operator    |       Reduces emits into array or user defined construct       | Array and string source |
| `reduceStream` | ReduceStream Operator | Reduces emits into array or user defined construct as a stream | Array and string source |
| `filter`       |    Filter Operator    |                     Filters passing values                     |                     any |
| `flat`         |     Flat Operator     |                Flattens a array down one level                 |                   array |
| `share`        |    Share Operator     |         Shares source value with multiple subscribers          |                     any |
| `for`          |     For Operator      |              Adds iteration logic into a pipeline              |                     any |
| `only`         |     Only Operator     |   Adds a function to be only called upon certain conditions    |                     any |
| `fragment`     |   Fragment Operator   |            Fragments a source into multiple sources            |                     any |

Reactivity
=

## Base Principles

Everything can become a source by using the flow `=>` symbol. Postfixing a source with the pipe symbol maps the
following operators over the source.

```
"Hello World"
    => print
```

```
source()
    => square
```

As many opertors as desired chan be chained following a source.

```
[1,2,3]
    => square
    => substractOne
    => print
```

The newline following the source is optional.

```
42 => print
```

Sources can be stored as `variables` and `constants` prefixing the source.

```
const helloWorld => "Hello World"

let numbers => [1,2,3]

// And with operators
const capitalized = "Hello World"
    => capitalize()
```

## Source type

The type signature is declared using the source `~` symbol.

```
// Variable declaration of int source, not instantiated
let source ~int

// Function taking and returning a string source
(source ~string) ~string => source
```

Notice the difference between the soure value and type.

```
// numbers is of type ~int and has value => [1,2,3]
const numbers ~int => [1,2,3]
```

## Flow

The values of a source don't start "flowing" until it has at least one subscriber. To consume values from a source the
subscribe `~>` is used.

```
[1,2,3]
    => split
    => multiply(2)
    => add(2)
    ~> print
```

## Operators

### Error

> TODO

### Closed

> TODO

### Split

Split splits incoming source value into seperate emits.

```
// Prints 1 2 3 4 5 on seperate lines
[1,2,3,4,5]
    => split
    ~> print
```

Split takes an optional argument in the form of a function. The signature of the function
is `<T>(value T, index int) bool`

```
// Prints 12\n 34\n 5
[1,2,3,4,5]
    => split (_, i) => i % 2 == 0
    ~> print
```

### Reduce

Reduces seperate emits into a collection as a single emit, the default behaviour is joining all emits into a single
array until the source is completed.

```
// Prints [1,2,3,4,5]
[1,2,3,4,5]
    => split
    => reduce
    ~> print
```

Reduce takes an optional reducer function.

```
// Prints [3,7,5]
[1,2,3,4,5]
    => split
    => reduce reduceEverySecondIndex 
    ~> print

// reduceEverySecondIndex adds every second index of an array
const reduceEverySecondIndex = (acc int[], val int, index int) int[] => {
    if index % 2 {
        acc += val
    } else {
        acc[len(acc)] += val
    }

    return acc
}
```

### ReduceStream

ReduceStream works like reduce but takes a second argument for when to emit.

```
//Prints [3]\n [7]\n [5]\n
[1,2,3,4,5]
    => split
    => reduceStream reduceEverySecondIndex 2  // Emits every second received value
    ~> print

// reduceEverySecondIndex adds every second index of an array
const reduceEverySecondIndex = (acc int[], val int, index int) int[] => {
    if index % 2 {
        acc += val
    } else {
        acc[len(acc)] += val
    }

    return acc
}
```

### For

For adds iteration logic to a pipeline, it takes a iterator function as a optional argument. It emits the last received
value.

```
// Prints "Iteration: 1!", "Iteration: 2!"... forever
"Iteration:"
    => for
    => (msg string) => "<msg> <count>!"
    => sleep 1000
    ~> print
```

For takes a optional argument for how may iterations should occur, in the form of a iterator
function `<T>(val T, i int) bool`.

```
// Prints 1 to 10
void
    => for (_, i) => i < 10
    ~> print
```

### Filter

> TODO

### Flat

> TODO

### Once

> TODO

### Share

Share shares the source over multiple subscribers.

```
// Prints two identical random values every second
const random => void
    => for
    => sleep 1000
    => () => random()
    => share

random
    ~> print

random
    ~> print
```

To show what share actually does, take a look at the previous example without the share operator. This version prints
two different values every second, the reason being each subscribe action creates a new source to consume of; using the
share operator this is not the case.

```
// Prints two different random values every second
const random => void
    => for
    => sleep 1000
    => () => random()

random
    ~> print

random
    ~> print
```

Functions
=

## Base structure and arguments

A function is declared by a set of desired arguments within a pair of parentheses followed by the flow `=>` symbol.
Multiple arguments can be seperated by comma's.

```
// A function without arguments, prints Hello World
() => print("Hello World")
```

```
// Prints the input + 1
(input int) => print(input + 1)
```

For multiple arguments of the same type only one type declaration is required; the following two functions are exactly
equal in operation.

```
// Two add functions
(a, b int) => print(a + b)
(a int, b int) => print(a + b)
```

## Using functions

## Multiline functions

Functions requiring multiple lines of statements can define a function body using braces.

```
(a, b int) => {
    a += 2
    b += 3
    print(a + b)
}
```

## Returning values

Functions returing a non-void value declare a return type following the parentheses containing the optional argument
list.

```
// Identity function
(a string) string => a
```

## Functions taking a source as argument

```
(source ~int) => ...
```


## Error handling
```
[0,1,2,3]
    => multiply     // multiply by 0 results in a error
    => add(5)       // add can't act upon error type
    ~> printResult  // print might format expecting integers

// posbile solution
[0,1,2,3]
    => multiply     // multiply by 0 results in a error; pass along in error monad? Does all values passing along get passed in monads then? Conveniance in this?
    => add(5)       // error is passed passed add until the first function that accepts error type
    => catch        // catch accepts the error type function, non error will be passed along untouched, to prevent it causing errors in printResult again dev can choose not to pass error further like filter.
    ~> printResult  // all integer values end up here being printed, all error values never reach here

// Something has to be constructed to quit the whole pipeline after 1 error, or after n errors or whatever case someone wants; can be implemented with a filter?
```

# Inspiration
1. RxJs
2. Svelte
3. TypeScript
4. Eiffel
5. Go
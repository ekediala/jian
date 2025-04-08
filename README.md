# Jian Programming Language Interpreter

Jian is a simple interpreter for a dynamic, C-like programming language, written entirely in Go. This project is heavily inspired by and follows the structure outlined in the book ["Writing An Interpreter In Go"](https://interpreterbook.com/) by Thorsten Ball.

It includes a lexer, parser, abstract syntax tree (AST), object system, and evaluator, providing a functional environment to execute Jian scripts or interact with the language via a REPL (Read-Eval-Print Loop).

## Features

*   **C-like Syntax:** Familiar syntax for variable bindings, function calls, and control flow.
*   **Variable Bindings:** Using the `let` keyword.
*   **Data Types:**
    *   Integers (`int64`)
    *   Booleans (`true`, `false`)
    *   Strings (`"Hello, World!"`)
    *   Arrays (`[1, "two", true]`)
    *   Hashes (Dictionaries/Maps) (`{"key": "value", 1: true}`)
*   **Operators:**
    *   Arithmetic: `+`, `-`, `*`, `/`
    *   Comparison: `==`, `!=`, `<`, `>`
    *   Logical Prefix: `!` (negation)
    *   Integer Prefix: `-` (negation)
*   **Control Flow:** `if`/`else` expressions.
*   **Functions:**
    *   First-class and higher-order functions.
    *   Closures (functions retain access to their definition environment).
*   **Return Statements:** Explicit `return` from functions.
*   **Indexing:** Access elements in Arrays and Hashes (`myArray[0]`, `myHash["key"]`).
*   **Built-in Functions:** Common utilities like `len`, `puts`, `first`, `last`, `rest`, `push`.
*   **REPL:** Interactive command-line interface.
*   **Error Handling:** Reports syntax and runtime errors gracefully.

## Requirements

*   Go version 1.21 or later (as specified in `go.mod`, although slightly older versions might work).

## Installation

You can install the `jian` binary using `go install`:

```bash
go install github.com/ekediala/jian@latest
```

Ensure your Go bin directory (`$GOPATH/bin` or `$HOME/go/bin`) is in your system's `PATH`.

Alternatively, clone the repository and build from source:

```bash
git clone https://github.com/ekediala/jian.git
cd jian
go build .
# Now you can run ./jian
```

## Usage

### 1. Interactive REPL

Simply run the `jian` command to start the REPL:

```bash
jian
```

You'll be greeted with a prompt (`>> `) where you can type Jian code:

```
Hello <your-username>! This is the Jian programming language!
Feel free to type in commands
>> let x = 10;
>> let y = 20;
>> let add = fn(a, b) { return a + b; };
>> add(x, y);
30
>> let myArray = [1, "hello", true];
>> len(myArray);
3
>> myArray[1];
hello
>> let myHash = {"name": "Jian", "version": 1};
>> myHash["name"];
Jian
>> puts("Hello from Jian REPL!");
Hello from Jian REPL!
null
>>
```

Type `Ctrl+D` or `Ctrl+C` to exit the REPL.

### 2. Executing Files

You can execute a file containing Jian code by passing the filename as an argument:

```bash
jian your_script.jian
```

For example, if `your_script.jian` contains:

```jian
let message = "Hello from a Jian script!";
puts(message);

let multiply = fn(a, b) { a * b };
let result = multiply(7, 6);
puts("7 * 6 =", result);
```

Running `jian your_script.jian` would output:

```
Hello from a Jian script!
7 * 6 = 42
```

## Language Overview & Examples

```jian
// Variable binding
let five = 5;
let ten = 10;

// Function definition
let add = fn(x, y) {
  x + y; // Implicit return of the last expression
};

// Function call
let result = add(five, ten);
puts("Result of add:", result); // Output: Result of add: 15

// Strings
let greeting = "Hello";
let name = "Jian";
puts(greeting + " " + name + "!"); // Output: Hello Jian!

// Arrays
let numbers = [1, 2, 3, 4];
puts("First number:", first(numbers)); // Output: First number: 1
puts("Last number:", last(numbers));  // Output: Last number: 4
let rest_nums = rest(numbers);
puts("Rest of numbers:", rest_nums); // Output: Rest of numbers: [2, 3, 4]
let pushed_nums = push(numbers, 5);
puts("Pushed numbers:", pushed_nums); // Output: Pushed numbers: [1, 2, 3, 4, 5]
puts("Original numbers:", numbers);   // Output: Original numbers: [1, 2, 3, 4] (push is non-mutating)
puts("Element at index 2:", numbers[2]); // Output: Element at index 2: 3

// Hashes
let person = {"name": "Alice", "age": 30};
puts(person["name"] + " is " + person["age"] + " years old."); // Output: Alice is 30 years old.
puts("Keys not present return null:", person["city"]); // Output: Keys not present return null: null

// Conditionals
let checkAge = fn(age) {
  if (age < 18) {
    return "Minor";
  } else {
    return "Adult";
  }
};
puts("Age 25 is:", checkAge(25)); // Output: Age 25 is: Adult
puts("Age 15 is:", checkAge(15)); // Output: Age 15 is: Minor

// Closures
let newAdder = fn(x) {
  fn(y) { x + y }; // Inner function closes over x
};
let addTwo = newAdder(2);
puts("2 + 5 is:", addTwo(5)); // Output: 2 + 5 is: 7
```

## Built-in Functions

*   `len(arg)`: Returns the length of a string or array.
    *   `len("hello")` -> `5`
    *   `len([1, 2])` -> `2`
*   `first(array)`: Returns the first element of an array, or `null` if empty.
    *   `first([1, 2])` -> `1`
    *   `first([])` -> `null`
*   `last(array)`: Returns the last element of an array, or `null` if empty.
    *   `last([1, 2])` -> `2`
    *   `last([])` -> `null`
*   `rest(array)`: Returns a *new* array containing all elements *except* the first, or `null` if empty.
    *   `rest([1, 2, 3])` -> `[2, 3]`
    *   `rest([1])` -> `[]`
    *   `rest([])` -> `null`
*   `push(array, element)`: Returns a *new* array with the `element` added to the end.
    *   `push([1, 2], 3)` -> `[1, 2, 3]`
    *   `push([], 1)` -> `[1]`
*   `puts(...)`: Prints arguments to the standard output, separated by newlines, and returns `null`.
    *   `puts("Hello", "World")` -> prints "Hello\nWorld\n"

## Development

### Building

To build the interpreter from source:

```bash
go build .
```

### Testing

The project includes unit tests for the lexer, parser, evaluator, and AST. Run tests using:

```bash
go test ./...
```

Or with coverage:

```bash
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

## Acknowledgements

*   Thorsten Ball for the invaluable ["Writing An Interpreter In Go"](https://interpreterbook.com/) book, which served as the primary guide for this project.

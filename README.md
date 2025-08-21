# CSVParser — Go CSV Parsing Library

This library provides a simple, efficient CSV parser in Go with support for standard CSV quoting and field handling rules. It is designed to read lines from a CSV file one at a time using an `io.Reader`, making it memory-efficient even for very large files.

## 📦 Interface

```go
type CSVParser interface {
    ReadLine(r io.Reader) (string, error)
    GetField(n int) (string, error)
    GetNumberOfFields() int
}
```

## 🛠 Errors

```go
var (
    ErrQuote      = errors.New("excess or missing \" in quoted-field")
    ErrFieldCount = errors.New("wrong number of fields")
)
```

---

## 🧩 Method Details

### `ReadLine(r io.Reader) (string, error)`

- Reads **one line** from the input stream (`io.Reader`).
- Handles newline variations: `\n`, `\r\n`, `\r`.
- Handles fields wrapped in double quotes (`"..."`), including escaped quotes (`""`).
- Returns an `ErrQuote` if a line contains an unmatched or excessive quote.
- Returns `io.EOF` when the end of the file is reached.
- The returned line does not include the newline character.
- Designed to handle large files efficiently by reading only up to the next newline.

### `GetField(n int) (string, error)`

- Returns the **n-th field** (0-based index) from the most recent line read by `ReadLine`.
- Returns `ErrFieldCount` if `n < 0` or `n >= GetNumberOfFields()`.
- Handles quoted fields and automatically removes the quotes.

### `GetNumberOfFields() int`

- Returns the number of fields in the **last line read** by `ReadLine`.
- Behavior is undefined if called before any `ReadLine` call.
- Safe to call after `ReadLine` returns `io.EOF`.

---

## 🔬 Example Usage

```go
package main

import (
    "fmt"
    "io"
    "os"

    "your-module-name/csvparser"
)

func main() {
    file, err := os.Open("example.csv")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    parser := csvparser.NewParser()

    for {
        line, err := parser.ReadLine(file)
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Println("ReadLine error:", err)
            return
        }

        fmt.Println("Line:", line)
        fmt.Println("Fields:", parser.GetNumberOfFields())
        for i := 0; i < parser.GetNumberOfFields(); i++ {
            field, err := parser.GetField(i)
            if err != nil {
                fmt.Printf("  Field %d error: %v\n", i, err)
                continue
            }
            fmt.Printf("  Field %d: %s\n", i, field)
        }
    }
}
```

---

## ✅ Notes

- **Efficiency**: The parser reads one line at a time and does not load the whole file into memory.
- **Robustness**: Handles edge cases like missing quotes, binary files, and arbitrary field lengths.
- **Quoting**: Correctly supports CSV quoting rules including escaped quotes (`""`).
- **Testing**: You should test the parser on different CSV formats and even malformed files to ensure stability.

---

## 🧪 Testing Tips

Try testing with:
- A file with `1,000,000+` lines
- Fields with embedded commas and quotes
- Empty fields and lines
- Binary files (e.g., `.png`, `.exe`) to see how parser behaves

---

## 📁 Project Structure Example

```
csvparser/
├── parser.go       # implementation

main.go             # demo usage
README.md
go.mod
```


package main

import (
	"fmt"
	"io"
	"os"

	"a-library-for-others/csvparser"
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
			val, _ := parser.GetField(i)
			fmt.Printf("Field %d: %s\n", i, val)
		}
	}
}

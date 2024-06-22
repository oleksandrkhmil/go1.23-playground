package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
)

const limit = 10

type RandomValuesGenerator struct{}

// All returns iteration index and value pairs
func (g RandomValuesGenerator) All() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for i := 0; i < limit; i++ {
			if !yield(i, rand.IntN(100)) {
				fmt.Println("Received stop")
				return
			}
		}
		fmt.Println("Limit reached")
	}
}

type FileReader struct {
	file string
}

func NewFileReader(file string) FileReader {
	return FileReader{file: file}
}

func (r FileReader) All() iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		file, err := os.Open(r.file)
		if err != nil {
			yield("", fmt.Errorf("open: %w", err))
			return
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, err := reader.ReadLine()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				yield("", fmt.Errorf("read line: %w", err))
				return
			}
			if !yield(string(line), nil) {
				return
			}
		}
	}
}

func main() {
	{
		fmt.Println("Exercise 1: Base iterator usage with slice")
		sl := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
		for i, s := range slices.All(sl) {
			fmt.Printf("%d: %s; ", i, s)
		}
		// Output: 0: a; 1: b; 2: c; 3: d; 4: e; 5: f; 6: g; 7: h;
	}

	fmt.Print("\n\n")

	{
		fmt.Println("Exercise 2: Base iterator usage with map")
		m := map[string]string{"Apple": "Unites States", "Samsung": "South Korea", "Xiaomi": "China"}
		for k, v := range maps.All(m) {
			fmt.Printf("%s: %s; ", k, v)
		}
		// Output: Apple: Unites States; Samsung: South Korea; Xiaomi: China;
	}

	fmt.Print("\n\n")

	{
		fmt.Println("Exercise 3: Custom iterator usage with range")
		generator := RandomValuesGenerator{}
		for i, v := range generator.All() {
			fmt.Printf("%d: %d; ", i, v)
		}
		// Output: 0: 72; 1: 6; 2: 57; 3: 21; 4: 57; 5: 54; 6: 45; 7: 1; 8: 31; 9: 67; Limit reached
	}

	fmt.Print("\n")

	{
		fmt.Println("Exercise 4: Custom iterator usage with iter.Pull2")
		generator := RandomValuesGenerator{}
		next, stop := iter.Pull2(generator.All())
		defer stop()
		for i, v, ok := next(); ok; i, v, ok = next() {
			fmt.Printf("%d: %d; ", i, v)
		}
		// Output: 0: 27; 1: 88; 2: 57; 3: 52; 4: 66; 5: 7; 6: 24; 7: 44; 8: 41; 9: 34; Limit reached

		fmt.Println("\nExercise 4.1: Call iterator one more time")

		i, v, ok := next()
		fmt.Printf("%d: %d: %v; \n", i, v, ok)
		// Output: 0: 0: false;
	}

	fmt.Print("\n")

	{
		fmt.Println("Exercise 5: Custom iterator usage with iter.Pull2 and custom stop")
		generator := RandomValuesGenerator{}
		next, stop := iter.Pull2(generator.All())

		for i := 0; i < 5; i++ {
			j, v, ok := next()
			if !ok {
				break
			}
			fmt.Printf("%d: %d; ", j, v)
		}
		stop()
		// Output: 0: 22; 1: 0; 2: 18; 3: 93; 4: 11; Received stop

		fmt.Println("\nExercise 5.1: Call iterator one more time")

		i, v, ok := next()
		fmt.Printf("%d: %d: %v; \n", i, v, ok)
		// Output: 0: 0: false;

		fmt.Println("\nExercise 5.2: Call stop one more time")
		stop() // No panic
		fmt.Println("OK")
	}

	fmt.Print("\n")

	{
		fmt.Println("Exercise 6: Read file with iterator")
		reader := NewFileReader("./dump.txt")
		next, stop := iter.Pull2(reader.All())
		defer stop()

		for line, err, ok := next(); ok; line, err, ok = next() {
			switch {
			case err != nil:
				fmt.Println("Error: " + err.Error())
			case strings.Contains(line, "STOP"):
				fmt.Println("Stop: " + line)
				stop()
			default:
				fmt.Println("Read line: " + line)
			}
		}
		/*
			Output:
				Read line: Lorem ipsum dolor sit amet
				Stop: Donec malesuada suscipit nulla, STOP HERE
		*/
	}
}

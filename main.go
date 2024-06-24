package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {

	// Define map to store values
	flags := map[string]interface{}{
		"c": flag.String("c", "", "Count bytes in a file"),
		"l": flag.String("l", "", "Count lines in a file"),
		"w": flag.String("w", "", "Count words in a file"),
		"m": flag.String("m", "", "Count characters in a file"),
	}

	flag.Parse()

	if flag.NFlag() == 0 {
		defaultFile := flag.Args()[0]

		// Perform the default action: count bytes
		byteCount, err := byteCounter(defaultFile)
		if err != nil {
			fmt.Println("Error checking file:", err)
			return
		}
		wordCount, err := wordCounter(defaultFile)
		if err != nil {
			fmt.Println("Error checking file:", err)
			return
		}
		lineCount, err := lineCounter(defaultFile)
		if err != nil {
			fmt.Println("Error checking file:", err)
			return
		}

		fmt.Printf("%d %d %d %s", byteCount, lineCount, wordCount, defaultFile)
		return
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "c":
			countBytesFileName := *flags["c"].(*string)

			if countBytesFileName == "" {
				fmt.Println("Please provide a Command")
				flag.Usage()
				return
			}

			byteCount, err := byteCounter(countBytesFileName)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s %d", countBytesFileName, byteCount)
		case "l":
			numberOfLinesFile := *flags["l"].(*string)

			if numberOfLinesFile == "" {
				fmt.Println("Please provide a Command")
				flag.Usage()
				return
			}

			linesNumber, err := lineCounter(numberOfLinesFile)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s %d", numberOfLinesFile, linesNumber)
		case "w":
			numberOfWordsFile := *flags["w"].(*string)

			if numberOfWordsFile == "" {
				fmt.Println("Please provide a Command")
				flag.Usage()
				return
			}

			wordCount, err := wordCounter(numberOfWordsFile)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s %d", numberOfWordsFile, wordCount)
		case "m":
			numberOfCharactersFile := *flags["m"].(*string)

			if numberOfCharactersFile == "" {
				fmt.Println("Please provide a Command")
				flag.Usage()
				return
			}

			characterCount, err := characterCounter(numberOfCharactersFile)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s %d", numberOfCharactersFile, characterCount)

		default:
			fmt.Printf("Unknown flag: %s\n", f.Name)
		}
	})
}

func byteCounter(fileName string) (int64, error) {
	// Check if file exists
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist")
			return 0, err
		} else {
			fmt.Println("Error checking file:", err)
			return 0, err
		}
	}

	return fileInfo.Size(), nil
}

func wordCounter(fileName string) (int, error) {
	dat, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	words := strings.Fields(string(dat))

	return len(words), nil
}

func characterCounter(fileName string) (int, error) {
	dat, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	characters := utf8.RuneCountInString(string(dat))

	return characters, nil
}

func lineCounter(numberOfLinesFile string) (int, error) {

	f, err := os.Open(numberOfLinesFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var r io.Reader
	r = f

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

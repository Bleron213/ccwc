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
		"c": flag.Bool("c", false, "Count bytes in a file"),
		"l": flag.Bool("l", false, "Count lines in a file"),
		"w": flag.Bool("w", false, "Count words in a file"),
		"m": flag.Bool("m", false, "Count characters in a file"),
	}

	flag.Parse()

	if flag.NFlag() == 0 {
		defaultFile := flag.Args()[0]

		// Perform the default action: count bytes
		reader, _ := openStdinOrFile()
		byteCount, err := byteCounter(reader)
		if err != nil {
			fmt.Println("Error checking file:", err)
			return
		}

		reader, _ = openStdinOrFile()
		wordCount, err := wordCounter(reader)
		if err != nil {
			fmt.Println("Error checking file:", err)
			return
		}

		reader, _ = openStdinOrFile()
		lineCount, err := lineCounter(reader)
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
			reader, fileName := openStdinOrFile()

			countBytesFileName := *flags["c"].(*bool)

			if !countBytesFileName {
				fmt.Println("Please provide a Command")
				flag.Usage()
				return
			}

			byteCount, err := byteCounter(reader)
			if err != nil {
				panic(err)
			}

			if fileName != "" {
				fmt.Printf("%s %d", fileName, byteCount)
			} else {
				fmt.Printf("%d", byteCount)
			}
		case "l":
			reader, fileName := openStdinOrFile()
			numberOfLinesFile := *flags["l"].(*bool)

			if !numberOfLinesFile {
				fmt.Println("Please provide a Command")
				flag.Usage()
				return
			}

			linesNumber, err := lineCounter(reader)
			if err != nil {
				panic(err)
			}

			if fileName != "" {
				fmt.Printf("%s %d", fileName, linesNumber)
			} else {
				fmt.Printf("%d", linesNumber)
			}
		case "w":
			reader, fileName := openStdinOrFile()

			numberOfWordsFile := *flags["w"].(*bool)

			if !numberOfWordsFile {
				fmt.Println("Please provide a Command")
				flag.Usage()
				return
			}

			wordCount, err := wordCounter(reader)
			if err != nil {
				panic(err)
			}

			if fileName != "" {
				fmt.Printf("%s %d", fileName, wordCount)
			} else {
				fmt.Printf("%d", wordCount)
			}

		case "m":
			reader, fileName := openStdinOrFile()

			numberOfCharactersFile := *flags["m"].(*bool)

			if !numberOfCharactersFile {
				fmt.Println("Please provide a Command")
				flag.Usage()
				return
			}

			characterCount, err := characterCounter(reader)
			if err != nil {
				panic(err)
			}

			if fileName != "" {
				fmt.Printf("%s %d", fileName, characterCount)
			} else {
				fmt.Printf("%d", characterCount)
			}

		default:
			fmt.Printf("Unknown flag: %s\n", f.Name)
		}
	})
}

func byteCounter(reader io.Reader) (int64, error) {
	var totalBytes int64
	buffer := make([]byte, 1024) // Define a buffer to hold chunks of data

	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// End of data reached
				break
			}
			return 0, err
		}
		totalBytes += int64(bytesRead)
	}

	return totalBytes, nil
}

func wordCounter(reader io.Reader) (int, error) {
	dat, err := io.ReadAll(reader)
	if err != nil {
		return 0, err
	}

	words := strings.Fields(string(dat))

	return len(words), nil
}

func characterCounter(reader io.Reader) (int, error) {
	dat, err := io.ReadAll(reader)
	if err != nil {
		return 0, err
	}
	characters := utf8.RuneCountInString(string(dat))

	return characters, nil
}

func lineCounter(reader io.Reader) (int, error) {

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := reader.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func openStdinOrFile() (io.Reader, string) {
	var err error
	var fileName string
	r := os.Stdin

	if len(os.Args) == 3 {
		fileName = os.Args[2]
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			return r, ""
		}

		r, err = os.Open(fileName)
		if err != nil {
			panic(err)
		}
	} else if len(os.Args) == 2 {
		fileName = os.Args[1]

		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			return r, ""
		}

		r, err = os.Open(fileName)
		if err != nil {
			panic(err)
		}
	}
	return r, fileName
}

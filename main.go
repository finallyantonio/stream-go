package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// Stream represents a data stream with a reader and a writer.
type Stream struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

// Pipe applies a transformation function to the stream.
func (s *Stream) Pipe(transform func(io.Reader, io.Writer)) *Stream {
	// Create a new pipe
	reader, writer := io.Pipe()

	// Launch a goroutine to apply the transformation
	go func() {
		transform(s.reader, writer)
		writer.Close()
	}()

	// Return a new Stream with the new reader and writer
	return &Stream{
		reader: reader,
		writer: writer,
	}
}

// NewStream initializes a new stream with input data.
func NewStream(input string) *Stream {
	reader, writer := io.Pipe()

	go func() {
		// Write input data to the stream
		writer.Write([]byte(input))
		writer.Close()
	}()

	return &Stream{
		reader: reader,
		writer: writer,
	}
}

func main() {
	input := "hello\nworld\nthis is a test\n"

	// Create a new stream with input data
	stream := NewStream(input)

	// Define transformation functions
	upperCaseTransform := func(r io.Reader, w io.Writer) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			w.Write([]byte(strings.ToUpper(scanner.Text() + "\n")))
		}
	}

	addPrefixTransform := func(r io.Reader, w io.Writer) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			w.Write([]byte("PREFIX: " + scanner.Text() + "\n"))
		}
	}

	// Build the stream pipeline
	stream = stream.Pipe(upperCaseTransform).Pipe(addPrefixTransform)

	// Read the final output
	finalScanner := bufio.NewScanner(stream.reader)
	for finalScanner.Scan() {
		fmt.Println(finalScanner.Text())
	}

	if err := finalScanner.Err(); err != nil {
		log.Fatalf("Error reading from final stream: %v", err)
	}
}

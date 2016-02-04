package main

import (
	"bufio"
	encoding "encoding/csv"
	"flag"
	"fmt"
	"github.com/wildducktheories/go-csv"
	"github.com/wildducktheories/influx-log-analyser/record"
	"io"
	"os"
	"strings"
)

func parse(input io.ReadCloser, builder csv.WriterBuilder) {
	output := builder(record.CsvHeaders)
	lines := bufio.NewReader(input)
	for {
		if line, err := lines.ReadString('\n'); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		} else {
			line = strings.TrimSpace(line)
			if r, err := record.ParseInfluxLogLine(line); err != nil {
				if err != record.ErrNotHttp && err != record.ErrInvalidNumberOfTokens {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
			} else {
				cr := output.Blank()
				r.Encode(cr)
				output.Write(cr)
			}
		}
	}
	output.Close(nil)
}

func main() {

	tabs := false
	comma := ","
	flag.BoolVar(&tabs, "tabs", false, "Use tabs as the output delimiter.")
	flag.Parse()

	if tabs {
		comma = "\t"
	}

	csvEncoder := encoding.NewWriter(os.Stdout)
	csvEncoder.Comma = rune(comma[0])
	parse(os.Stdin, csv.WithCsvWriter(csvEncoder, os.Stdout))
}

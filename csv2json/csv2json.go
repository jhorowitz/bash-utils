package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

var asArray = flag.Bool("as-array", false, "format as json array")

func init() {
	flag.Parse()
}

func main() {
	reader := csv.NewReader(os.Stdin)
	var headers []string
	idx := -1
	for {
		idx++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.WithError(err).Fatal("failed to convert csv")
		}
		if headers == nil {
			headers, err = makeHeaders(row)
			if err != nil {
				logrus.WithError(err).Fatal("failed to convert headers")
			}
			if *asArray {
				_, _ = fmt.Fprint(os.Stdout, "[\n")
			}
			continue
		}
		if len(headers) != len(row) {
			logrus.WithField("headers", len(headers)).WithField("row", len(row)).WithField("idx", idx).Fatal("header/row count mismatch")
		}
		m := make(map[string]string)
		for i := range headers {
			m[headers[i]] = row[i]
		}

		b, err := json.Marshal(m)
		if err != nil {
			logrus.WithError(err).Fatal("failed to marshal json")
		}
		if *asArray {
			if idx == 1 {
				_, _ = fmt.Fprintf(os.Stdout, "%s", b)
			} else {
				_, _ = fmt.Fprintf(os.Stdout, ",\n%s", b)
			}
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "%s\n", b)
		}
	}
	if *asArray {
		_, _ = fmt.Fprint(os.Stdout, "\n]\n")
	}
}

func makeHeaders(h []string) ([]string, error) {
	var result []string
	seen := make(map[string]bool)
	for _, v := range h {
		v = strcase.ToSnake(v)
		if seen[v] {
			return nil, fmt.Errorf("duplicate headers can't be converted to json: %s", strconv.Quote(v))
		}
		seen[v] = true
		result = append(result, v)
	}

	return result, nil
}

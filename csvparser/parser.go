package csvparser

import (
	"errors"
	"io"
)

var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

type CSVParser interface {
	ReadLine(r io.Reader) (string, error)
	GetField(n int) (string, error)
	GetNumberOfFields() int
}

type parser struct {
	buffer []byte
	fields []string
}

func NewParser() CSVParser {
	return &parser{}
}

func (p *parser) ReadLine(r io.Reader) (string, error) {
	p.buffer = p.buffer[:0]
	p.fields = nil
	tmp := make([]byte, 1)
	inQuotes := false

	for {
		n, err := r.Read(tmp)
		if n > 0 {
			c := tmp[0]

			// \r — пропускаем
			if c == '\r' {
				continue
			}

			if c == '"' {
				inQuotes = !inQuotes
			}

			if c == '\n' && !inQuotes {
				break
			}

			p.buffer = append(p.buffer, c)
		}

		if err == io.EOF {
			if inQuotes {
				return "", ErrQuote
			}
			if len(p.buffer) == 0 {
				return "", io.EOF
			}
			break
		}

		if err != nil {
			return "", err
		}
	}

	var err error
	p.fields, err = parseFields(p.buffer)
	if err != nil {
		return "", err
	}

	return string(p.buffer), nil
}

func (p *parser) GetField(n int) (string, error) {
	if n < 0 || n >= len(p.fields) {
		return "", ErrFieldCount
	}
	return p.fields[n], nil
}

func (p *parser) GetNumberOfFields() int {
	return len(p.fields)
}

// parseFields parses a single CSV line into fields
func parseFields(line []byte) ([]string, error) {
	var fields []string
	var field []byte
	inQuotes := false

	for i := 0; i < len(line); i++ {
		c := line[i]

		if inQuotes {
			if c == '"' {
				// Check for escaped quote
				if i+1 < len(line) && line[i+1] == '"' {
					field = append(field, '"')
					i++
				} else {
					inQuotes = false
				}
			} else {
				field = append(field, c)
			}
		} else {
			if c == ',' {
				fields = append(fields, string(field))
				field = field[:0]
			} else if c == '"' {
				if len(field) == 0 {
					inQuotes = true
				} else {
					return nil, ErrQuote
				}
			} else {
				field = append(field, c)
			}
		}
	}

	if inQuotes {
		return nil, ErrQuote
	}

	fields = append(fields, string(field))
	return fields, nil
}

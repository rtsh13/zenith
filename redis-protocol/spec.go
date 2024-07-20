package resp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/zenith"
)

var (
	ErrScanningInput         = zenith.ScannerError{Message: "error : [%v] in scanning input"}
	ErrInvalidStartCharacter = zenith.ProtocolError{Message: "invalid start character : %v"}
	ErrInvalidBulkString     = zenith.ProtocolError{Message: "invalid bulk string prefix. expected : %v, got : %v"}
	ErrInvalidBulkLength     = zenith.ProtocolError{Message: "invalid bulk length"}
	ErrUnexpectedEndOfStream = zenith.ProtocolError{Message: "unexpected end of stream"}
)

type resp struct{}

type Protocol interface {
	Serialize([]string) string
	Deserialize(string) (strings.Builder, error)
}

func New() Protocol {
	return &resp{}
}

func (r *resp) Serialize(input []string) string {
	arr := make([]string, 0)

	for _, value := range input {
		arr = append(arr, fmt.Sprintf("$%d%s%v", len(value), zenith.CRLF, value))
	}

	response := strings.Join(arr, zenith.CRLF)

	return fmt.Sprintf("*%d%s%v%s", len(input), zenith.CRLF, response, zenith.CRLF)
}

func (r *resp) Deserialize(input string) (strings.Builder, error) {
	parsedCommand, err := r.parseCMD(input)
	if err != nil {
		return strings.Builder{}, err
	}

	instructions, err := r.insBuilder(parsedCommand.String())
	if err != nil || instructions.Len() <= 0 {
		return strings.Builder{}, err
	}

	return instructions, nil
}

func (r *resp) parseCMD(input string) (strings.Builder, error) {
	var response strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)

	if !scanner.Scan() || !strings.HasPrefix(scanner.Text(), "*") {
		return response, fmt.Errorf(ErrInvalidStartCharacter.Message, scanner.Text())
	}

	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if !strings.HasPrefix(line, "$") {
			return strings.Builder{}, fmt.Errorf(ErrInvalidBulkString.Message, "$", scanner.Text())
		}

		response.WriteString(strings.TrimPrefix(line, "$"))

		if !scanner.Scan() {
			return strings.Builder{}, ErrScanningInput
		}

		nextLine := scanner.Text()
		response.WriteString(nextLine)
	}

	return response, nil
}

func (r *resp) insBuilder(instruction string) (strings.Builder, error) {
	cursor := 0

	cmd := strings.Builder{}
	for cursor < len(instruction) {
		cFootprint := cursor

		for instruction[cursor] >= '0' && instruction[cursor] <= '9' {
			cursor++
		}

		if cFootprint == cursor {
			return cmd, ErrInvalidBulkLength
		}

		size, err := strconv.Atoi(instruction[cFootprint:cursor])
		if err != nil {
			return cmd, ErrInvalidBulkLength
		}

		if size+cursor > len(instruction) {
			return cmd, ErrUnexpectedEndOfStream
		}

		if cmd.Len() > 0 {
			cmd.WriteString(" ")
		}

		cmd.WriteString(instruction[cursor : size+cursor])
		cursor += size
	}

	return cmd, nil
}

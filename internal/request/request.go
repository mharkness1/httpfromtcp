package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	strReq := string(input)
	strParts := strings.Split(strReq, "\r\n")
	requestLine, err := parseRequestLine(strParts[0])
	if err != nil {
		return nil, err
	}

	fullRequest := Request{
		RequestLine: requestLine,
	}

	return &fullRequest, nil
}

func parseRequestLine(req string) (RequestLine, error) {
	parts := strings.Split(req, " ")
	if len(parts) != 3 {
		return RequestLine{}, errors.New("Incorrect number of request line parts")
	}
	if strings.ToUpper(parts[0]) != parts[0] || !IsLetter(parts[0]) {
		return RequestLine{}, errors.New("Method not capitalised.")
	}
	if parts[2] != "HTTP/1.1" {
		return RequestLine{}, errors.New("HTTP version not supported, must be HTTP/1.1")
	}
	httpVersion := strings.Split(parts[2], "/")

	request := RequestLine{
		HttpVersion:   httpVersion[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}

	return request, nil
}

func IsLetter(s string) bool {
	return !strings.ContainsFunc(s, func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z')
	})
}

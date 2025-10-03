package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

const crlf = "\r\n"
const bufferSize = 8

type requestStatus int

const (
	requestInitialized    requestStatus = iota //0
	requestParsingHeaders                      //1
	requestParsingBody                         //2
	requestDone                                //3
)

type Request struct {
	RequestLine   RequestLine
	Headers       headers.Headers
	Body          []byte
	RequestStatus requestStatus
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	parsedRequest := &Request{
		Headers:       headers.NewHeaders(),
		RequestStatus: requestInitialized,
		Body:          make([]byte, 0),
	}

	buffer := make([]byte, bufferSize)

	bufferIdx := 0

	for parsedRequest.RequestStatus != requestDone {

		if bufferIdx >= len(buffer) {
			tmpBuffer := make([]byte, len(buffer)*2)
			copy(tmpBuffer, buffer[:bufferIdx])
			buffer = tmpBuffer
		}

		readBytes, err := reader.Read(buffer[bufferIdx:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				if parsedRequest.RequestStatus != requestDone {
					return nil, fmt.Errorf("incomplete request")
				}
				break
			}
			return nil, err
		}

		bufferIdx += readBytes

		parsedBytes, err := parsedRequest.parse(buffer[:bufferIdx])
		if err != nil {
			return nil, err
		}

		if parsedBytes > 0 {
			copy(buffer, buffer[parsedBytes:])
			bufferIdx -= parsedBytes
		}
	}
	return parsedRequest, nil
}

func (r *Request) parse(data []byte) (int, error) {

	parsedBytes := 0

	for r.RequestStatus != requestDone {
		n, err := r.parseSingle(data[parsedBytes:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return parsedBytes, nil
		}
		parsedBytes += n
	}

	return parsedBytes, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.RequestStatus {
	case requestInitialized:
		requestLine, parsedBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if parsedBytes == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.RequestStatus = requestParsingHeaders

		return parsedBytes, nil

	case requestParsingHeaders:
		parsedBytes, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if parsedBytes == 0 {
			return 0, nil
		}
		if done {
			r.RequestStatus = requestParsingBody
		}
		return parsedBytes, nil
	case requestParsingBody:
		contentLength, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.RequestStatus = requestDone
			return 0, nil
		}

		contentLengthNum, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, fmt.Errorf("Error converting to int %s", contentLength)
		}

		r.Body = append(r.Body, data...)

		if len(r.Body) > contentLengthNum {
			return 0, fmt.Errorf("Content is greater than provided length %d", contentLengthNum)
		}
		if len(r.Body) == contentLengthNum {
			r.RequestStatus = requestDone
			return len(data), nil
		}

		return len(data), nil
	case requestDone:
		return 0, fmt.Errorf("Request has already been processed")
	default:
		return 0, fmt.Errorf("Incorrect status %q", r.RequestStatus)
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return nil, 0, nil
	}

	requestLine := string(data[:idx])

	parsedRequestLine, err := requestLineFromString(requestLine)
	if err != nil {
		return nil, 0, err
	}
	return parsedRequestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {

	requestLineParts := strings.Split(str, " ")
	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("Request line does not consist of 3 parts: %s", str)
	}

	method := requestLineParts[0]

	for _, char := range method {
		if char > 'Z' || char < 'A' {
			fmt.Println(char)
			return nil, fmt.Errorf("Request method contains other than capital letters: %s", method)
		}
	}

	requestTarget := requestLineParts[1]

	httpVersionParts := strings.Split(requestLineParts[2], "/")
	if len(httpVersionParts) != 2 {
		return nil, fmt.Errorf("Http version does not consist of 2 parts: %s", requestLineParts[2])
	}

	httpVersionNumber := httpVersionParts[1]

	if httpVersionNumber != "1.1" {
		return nil, fmt.Errorf("Http version not supported: %s", httpVersionNumber)
	}

	return &RequestLine{
		HttpVersion:   httpVersionNumber,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}

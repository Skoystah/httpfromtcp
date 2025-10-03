package handlers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const yourproblem = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

const myproblem = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

const banger = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

func NewHandler(w *response.Writer, req *request.Request) {

	switch {
	case req.RequestLine.RequestTarget == "/yourproblem":
		handler400(w, req)
		return
	case req.RequestLine.RequestTarget == "/myproblem":
		handler500(w, req)
		return
	case req.RequestLine.RequestTarget == "/":
		handler200(w, req)
		return
	case req.RequestLine.RequestTarget == "/video":
		handlervideo(w, req)
		return
	case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/"):
		handlerbin(w, req)
		return
	}
}

func handlerbin(w *response.Writer, req *request.Request) {
	targetSuffix := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := fmt.Sprintf("https://httpbin.org/%s", targetSuffix)

	res, err := http.Get(url)
	if err != nil {
		log.Printf("error retrieving data from %s", url)
		return
	}

	err = w.WriteStatusLine(response.OK)
	if err != nil {
		log.Printf("error sending response status line: %s", err)
		return
	}

	headers := response.GetDefaultHeaders(0)
	headers.Delete("content-length")
	headers.Set("Transfer-Encoding", "chunked")
	headers.Set("Trailer", "X-Content-SHA256, X-Content-Length")

	err = w.WriteHeaders(headers)
	if err != nil {
		log.Printf("error sending headers: %s", err)
		return
	}

	fullBody := make([]byte, 0)

	for {
		body := make([]byte, 1024)
		n, err := res.Body.Read(body)

		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("error reading data from body")
				return
			}
		}

		log.Printf("Reading %d bytes from body\n", n)
		w.WriteChunkedBody(body)

		fullBody = append(fullBody, body[:n]...)

		if errors.Is(err, io.EOF) || n == 0 {
			break
		}
	}

	bodyHash := sha256.Sum256(fullBody)
	log.Printf("Bodyhash %d", bodyHash[:])

	w.WriteChunkedBodyDone()
	trailers := response.GetTrailersFromHeader(headers)
	trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", bodyHash))
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
	w.WriteTrailers(trailers)
}

func handler400(w *response.Writer, _ *request.Request) {

	err := w.WriteStatusLine(response.BadRequest)
	if err != nil {
		log.Printf("error sending response status line: %s", err)
		return
	}

	body := []byte(yourproblem)
	headers := response.GetDefaultHeaders(len(body))
	headers.Update("content-type", "text/html")

	err = w.WriteHeaders(headers)
	if err != nil {
		log.Printf("error sending headers: %s", err)
		return
	}

	_, err = w.WriteBody(body)
	if err != nil {
		log.Printf("error writing body: %s", err)
		return
	}
}

func handler500(w *response.Writer, _ *request.Request) {

	err := w.WriteStatusLine(response.InternalServerError)
	if err != nil {
		log.Printf("error sending response status line: %s", err)
		return
	}

	body := []byte(myproblem)
	headers := response.GetDefaultHeaders(len(body))
	headers.Update("content-type", "text/html")

	err = w.WriteHeaders(headers)
	if err != nil {
		log.Printf("error sending headers: %s", err)
		return
	}

	_, err = w.WriteBody(body)
	if err != nil {
		log.Printf("error writing body: %s", err)
		return
	}
}

func handler200(w *response.Writer, _ *request.Request) {

	err := w.WriteStatusLine(response.OK)
	if err != nil {
		log.Printf("error sending response status line: %s", err)
		return
	}

	body := []byte(banger)
	headers := response.GetDefaultHeaders(len(body))
	headers.Update("content-type", "text/html")

	err = w.WriteHeaders(headers)
	if err != nil {
		log.Printf("error sending headers: %s", err)
		return
	}

	_, err = w.WriteBody(body)
	if err != nil {
		log.Printf("error writing body: %s", err)
		return
	}
}

func handlervideo(w *response.Writer, req *request.Request) {
	video, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		return
	}

	err = w.WriteStatusLine(response.OK)
	if err != nil {
		log.Printf("error sending response status line: %s", err)
		return
	}

	body := video
	headers := response.GetDefaultHeaders(len(body))
	headers.Update("content-type", "video/mp4")

	err = w.WriteHeaders(headers)
	if err != nil {
		log.Printf("error sending headers: %s", err)
		return
	}

	_, err = w.WriteBody(body)
	if err != nil {
		log.Printf("error writing body: %s", err)
		return
	}
}

package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"log"
)

type Writer struct {
	Writer       io.Writer
	WriterStatus writerStatus
}

type writerStatus int

const (
	writeStatusLine writerStatus = iota //0
	writeHeaders                        //2
	writeBody                           //3
	writeDone                           //4
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer:       w,
		WriterStatus: writeStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.WriterStatus != writeStatusLine {
		return fmt.Errorf("Incorrect status %q", w.WriterStatus)
	}

	reason, err := statusToString(statusCode)
	if err != nil {
		return err
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reason)
	_, err = w.Writer.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("Error writing status line %s: %s", statusLine, err)
	}

	w.WriterStatus = writeHeaders

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterStatus != writeHeaders {
		return fmt.Errorf("Incorrect status %q", w.WriterStatus)
	}

	for key, value := range headers {
		header := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Writer.Write([]byte(header))

		if err != nil {
			return fmt.Errorf("Error writing header %s: %s", header, err)
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))

	if err != nil {
		return fmt.Errorf("Error writing header \\r\\n: %s", err)
	}

	w.WriterStatus = writeBody

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterStatus != writeBody {
		return 0, fmt.Errorf("Incorrect status %q", w.WriterStatus)
	}

	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("Error writing body: %s", err)
	}

	w.WriterStatus = writeDone

	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.WriterStatus != writeBody {
		return 0, fmt.Errorf("Incorrect status %q", w.WriterStatus)
	}

	length := fmt.Sprintf("%x\r\n", len(p))
	_, err := w.Writer.Write([]byte(length))
	if err != nil {
		return 0, fmt.Errorf("Error writing body: %s", err)
	}

	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("Error writing body: %s", err)
	}
	_, err = w.Writer.Write([]byte("\r\n"))
	return n, nil
}
func (w *Writer) WriteChunkedBodyDone() (int, error) {
	zeroLength := fmt.Sprintf("%x\r\n", 0)
	n, err := w.Writer.Write([]byte(zeroLength))
	if err != nil {
		return 0, fmt.Errorf("Error writing body: %s", err)
	}

	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {

	for key, value := range h {
		header := fmt.Sprintf("%s: %s\r\n", key, value)
		log.Printf("Adding trailer %s\n", header)
		_, err := w.Writer.Write([]byte(header))

		if err != nil {
			return fmt.Errorf("Error writing header %s: %s", header, err)
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))

	if err != nil {
		return fmt.Errorf("Error writing header \\r\\n: %s", err)
	}

	return nil
}

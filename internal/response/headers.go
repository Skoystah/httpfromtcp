package response

import (
	"httpfromtcp/internal/headers"
	"strconv"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("content-length", strconv.Itoa(contentLen))
	headers.Set("connection", "close")
	headers.Set("content-type", "text/plain")

	return headers
}

func GetTrailersFromHeader(h headers.Headers) headers.Headers {
	trailers := headers.NewHeaders()

	// addedTrailers, exists := h.Get("Trailer")
	// if !exists {
	// 	return nil
	// }
	//
	// for trailerKey := range strings.Split(addedTrailers, ",") {
	// 	trailers.Set(trailerKey, "")
	// }
	// trailers.Set("content-length", strconv.Itoa(contentLen))
	// trailers.Set("connection", "close")

	return trailers
}

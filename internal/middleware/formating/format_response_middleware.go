package middleware

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"

	general_reponse_dto "github.com/SilentPlaces/basicauth.git/internal/dto/general"
)

// ResponseFormattingMiddleware intercepts the response from the next handler,
// wraps it into a standardized JSON response
// Error responses (>= 400) and 200 and empty responses without body are reformatted.
// Other responses (e.g. 301 redirects) are passed through as it is.
func ResponseFormattingMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cw := &CustomResponseWriter{ResponseWriter: w}
		next(cw, r, ps)

		//if the next handler did not set status code, set it to 200
		if cw.statusCode == 0 {
			cw.statusCode = http.StatusOK
		}

		// Handle error responses (>= 400)
		if cw.statusCode >= http.StatusBadRequest {
			resp := general_reponse_dto.Response{
				Status:  "error",
				Code:    cw.statusCode,
				Message: strings.TrimSpace(cw.bodyBuffer.String()),
			}
			cw.bodyBuffer.Reset()
			writeJSON(w, cw.statusCode, resp)
			return
		}

		// Handle 200 OK responses.
		if cw.statusCode == http.StatusOK {
			// If there's a body, attempt to parse it.
			if cw.bodyBuffer.Len() > 0 {
				var data interface{}
				if err := json.Unmarshal(cw.bodyBuffer.Bytes(), &data); err == nil {
					resp := general_reponse_dto.Response{
						Status:  "success",
						Code:    cw.statusCode,
						Message: "success",
						Data:    data,
					}
					writeJSON(w, cw.statusCode, resp)
					return
				}
				// If unmarshalling fails, fall back to sending the original response.
				w.WriteHeader(cw.statusCode)
				_, _ = w.Write(cw.bodyBuffer.Bytes())
				return
			}

			// send a standardized success message.
			resp := general_reponse_dto.Response{
				Status:  "success",
				Code:    cw.statusCode,
				Message: "success",
			}
			writeJSON(w, cw.statusCode, resp)
			return
		}

		// Default: pass through the original response codes including 301
		w.WriteHeader(cw.statusCode)
		_, _ = w.Write(cw.bodyBuffer.Bytes())
	}
}

// writeJSON sets the Content-Type header, writes the status code, and encodes the response as JSON.
func writeJSON(w http.ResponseWriter, statusCode int, resp general_reponse_dto.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

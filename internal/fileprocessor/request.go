package fileprocessor

import (
	"bytes"
	"io"
	"net/http"
)

type Response struct {
	Code int //status code
	Body []byte
}

func (w *FileWriter) makeRequest(method, url string, body []byte) (*Response, error) {
	log := w.logger.With().Str("method", method).Str("url", url).Logger()
	log.Debug().
		Str("body", string(body)).
		Msg("making http request")

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Info().Err(err).Msg("failed to create request")
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Info().Err(err).Msg("failed to submit request")
		return nil, err
	}

	return &Response{
		Code: resp.StatusCode,
		Body: b,
	}, nil
}

func (w *FileWriter) get(endpoint string) (*Response, error) {
	return w.makeRequest("GET", endpoint, nil)
}

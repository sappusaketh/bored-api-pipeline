package fileprocessor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/sappusaketh/bored-api-pipeline/internal/config"
)

const (
	filePerm   = os.FileMode(0664)
	fileFormat = "json"
)

type FileWriter struct {
	cfg              *config.Config
	file             *os.File
	bytesWritten     *int64
	fileCreationTime *time.Time
	logger           zerolog.Logger
}

func New(cfg *config.Config, logger zerolog.Logger) (*FileWriter, error) {
	// Check if the directory exists
	if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
		// Directory does not exist, create it
		if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			logger.Error().Err(err).Msg("Error creating directory")
			return nil, err
		}
		logger.Debug().Msg("Directory created successfully")
	} else if err != nil {
		logger.Error().Err(err).Msg("Error checking directory existence")
		return nil, err
	} else {
		logger.Debug().Msg("Directory exists")
	}
	return &FileWriter{
		cfg:    cfg,
		logger: logger,
	}, nil
}

// writes data to new file or appends to existing file
func (w *FileWriter) Write(data []byte) (int, error) {

	if w.file == nil || w.shouldRotate() {
		if err := w.rotateFile(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(append(data, []byte("\n")...))
	if err != nil {
		return n, err
	}

	w.bytesWritten = Ptr(*w.bytesWritten + int64(n))

	if *w.bytesWritten >= w.cfg.Rotate.Size {
		if err := w.rotateFile(); err != nil {
			return n, err
		}
	}

	return n, nil
}

// tells if current file should be closed
func (w *FileWriter) shouldRotate() bool {
	return time.Since(*w.fileCreationTime) >= w.cfg.Rotate.Interval
}

// closes current file if exists
func (w *FileWriter) close() error {
	if w.file != nil {
		return w.file.Close()
	}

	return nil
}

// closes the current file and rotates to a new file
func (w *FileWriter) rotateFile() error {
	if err := w.close(); err != nil {
		return err
	}

	creationTime := time.Now()
	timestamp := creationTime.Format("20060102150405")
	newFilePath := fmt.Sprintf("%s/%s.%s", w.cfg.OutputDir, timestamp, fileFormat)

	file, err := os.OpenFile(newFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePerm)
	if err != nil {
		return err
	}

	w.file = file
	w.bytesWritten = Ptr(int64(0))
	w.fileCreationTime = Ptr(creationTime)

	return nil
}

func (w *FileWriter) Start(startTime time.Time) {

	defer w.close()
	for time.Since(startTime) < w.cfg.MaxPollTime {
		// fetch data
		resp, err := w.get(w.cfg.Endpoint)
		if err != nil {
			w.logger.Fatal().Err(err).Msg("Error fetching data")
			return
		}
		if resp.Code != 200 {
			w.logger.Error().Int("statusCode", resp.Code).Msg("Request failed")
			continue
		}
		var data map[string]interface{}

		if err := json.Unmarshal(resp.Body, &data); err != nil {
			w.logger.Fatal().Err(err).Msg("Error parsing JSON:")
			return
		}
		w.logger.Debug().Any("response", data).Msg("Fetch successful")

		// Write response body to file
		if _, err := w.Write(resp.Body); err != nil {
			w.logger.Fatal().Err(err).Msg("Error writing response body")
			continue
		}

	}
}

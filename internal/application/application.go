package application

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/sappusaketh/bored-api-pipeline/internal/config"
	"github.com/sappusaketh/bored-api-pipeline/internal/fileprocessor"
)

type application struct {
	env   string
	runId string
}

type Application interface {
	Run()
}

func New(env string) Application {
	runId := uuid.New().String()
	if os.Getenv("RUN_ID") != "" {
		runId = os.Getenv("RUN_ID")
	}
	return &application{
		env:   env,
		runId: runId,
	}
}

func (a *application) Run() {
	startTime := time.Now()
	// ctx := context.Background()

	var logger = zerolog.New(os.Stdout).With().Str("runId", a.runId).Timestamp().Logger()
	if os.Getenv("LOG_LEVEL") == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	c, err := config.New(a.env, a.runId)
	if err != nil {
		logger.Fatal().Err(err).Msg(a.env + " config load failed, Makesure config yml file exists")
	}
	logger.Info().Any("config", c).Msg("Loaded config succesfully 🚀")
	fileProcessor, err := fileprocessor.New(c, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error initializing file processor")
	}
	fileProcessor.Start(startTime)
	logger.Info().Msg("finished processing files, exiting...")
}

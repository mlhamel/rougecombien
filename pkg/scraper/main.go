package scraper

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mlhamel/rougecombien/pkg/config"
	"github.com/mlhamel/rougecombien/pkg/gcloud"
)

const url = "https://www.cehq.gouv.qc.ca/suivihydro/fichier_donnees.asp?NoStation=040204"

type Scraper struct {
	cfg    *config.Config
	client *http.Client
}

type Result struct {
	ScrapedAt time.Time
	TakenAt   time.Time
	Outflow   float64
}

// Data return Result as an array of bytes
func (r *Result) Data() []byte {
	return []byte(fmt.Sprintf(
		"scrapedAt:%d,takenAt:%d,outflow:%f",
		r.ScrapedAt.Unix(),
		r.TakenAt.Unix(),
		r.Outflow,
	))
}

func NewScraper(cfg *config.Config) *Scraper {
	return &Scraper{cfg: cfg, client: &http.Client{}}
}

func (s *Scraper) Run(ctx context.Context) error {
	emiter, err := gcloud.NewPubSubEmission(ctx, s.cfg)

	if err != nil {
		return err
	}

	response, err := s.client.Get(url)

	if err != nil {
		return fmt.Errorf("Error while downloading: %w", err)
	}

	s.cfg.Logger().Info().
		Int("statusCode", response.StatusCode).
		Int64("ContentLength", response.ContentLength).
		Msgf("Response from %s", url)

	defer response.Body.Close()

	reader := csv.NewReader(response.Body)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()

	var record Result

	if err != nil {
		return fmt.Errorf("Error while parsing: %w", err)
	}

	for _, each := range data {
		if strings.TrimSpace(each[0]) == "Date" {
			continue
		}

		if strings.TrimSpace(each[0]) == "" {
			continue
		}

		rawDate := fmt.Sprintf("%s %s", each[0], strings.TrimRight(each[1], "\\"))

		takenAt, err := time.Parse("2006-01-02 15:04", rawDate)

		if err != nil {
			s.cfg.Logger().Error().Err(err).Msg("Error parsing date")
			continue
		}

		rawFlow := strings.Replace(each[2], ",", ".", 1)

		outflow, err := strconv.ParseFloat(rawFlow, 64)

		if err != nil {
			s.cfg.Logger().Error().Err(err).Msg("Error parsing flow")
			continue
		}

		record.ScrapedAt = time.Now().UTC()
		record.TakenAt = takenAt.UTC()
		record.Outflow = outflow

		s.cfg.Logger().
			Info().
			Time("ScrapedAt", record.ScrapedAt).
			Time("TakenAt", record.TakenAt).
			Float64("Outflow", record.Outflow).
			Msg("Record parsed")

		if err = emiter.Publish(ctx, s.cfg.TopicName(), record.Data()); err != nil {
			return err
		}
	}

	return nil
}

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
)

type Manager struct {
	storageClient *storage.Client
}

func NewManager(sc *storage.Client) Manager {
	return Manager{
		storageClient: sc,
	}
}

const GCS_BUCKET_NAME = "weather-data-responses"

func (m Manager) FetchWeatherData(input FetchWeatherDataRequest) error {
	latitudeStr := strconv.FormatFloat(input.Latitude, 'f', 2, 64)
	longitudeStr := strconv.FormatFloat(input.Longitude, 'f', 2, 64)

	// Contruct Open-Meteo Historical Weather API URL
	url := "https://archive-api.open-meteo.com/v1/archive"
	url += "?latitude=" + latitudeStr + "&longitude=" + longitudeStr + "&start_date=" + input.StartDate + "&end_date=" + input.EndDate
	url += "&daily=temperature_2m_max,temperature_2m_min,temperature_2m_mean,apparent_temperature_max,apparent_temperature_min,apparent_temperature_mean"

	log.Printf("Calling URL: %s", url)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the response body: %v", err)
	}

	fileName := "response-" + time.Now().Format("2006-01-02-15-04-05") + ".json"
	err = m.uploadToGCS(context.Background(), GCS_BUCKET_NAME, fileName, body)
	if err != nil {
		return err
	}

	return nil
}

// uploads response to GCS Bucket
func (m Manager) uploadToGCS(ctx context.Context, bucketName, objectName string, data []byte) error {
	client := m.storageClient
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	writer := obj.NewWriter(ctx)
	defer writer.Close()

	// Write data to GCS
	_, err := writer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to GCS: %w", err)
	}

	log.Printf("Data successfully uploaded to bucket %s as object %s\n", bucketName, objectName)
	return nil
}

// List all the files inside a GCS Bucket
func (m Manager) ListFiles() (ListWeatherFilesResponse, error) {
	ctx := context.Background()
	client := m.storageClient
	defer client.Close()

	bucket := client.Bucket(GCS_BUCKET_NAME)

	var fileNames []string
	it := bucket.Objects(ctx, nil) // List all objects in the bucket with no filtering

	for {
		objAttrs, err := it.Next()
		if err != nil {
			if err.Error() == "no more items in iterator" {
				break
			}
			return ListWeatherFilesResponse{}, fmt.Errorf("failed to list objects in bucket: %w", err)
		}

		fileNames = append(fileNames, objAttrs.Name)
	}

	return ListWeatherFilesResponse{
		Files: fileNames,
	}, nil
}

// Reads the content of a file stored in GCS
func (m Manager) FetchWeatherFileContent(fileName string) (interface{}, error) {
	client := m.storageClient
	defer client.Close()

	bucket := client.Bucket(GCS_BUCKET_NAME)
	obj := bucket.Object(fileName)

	ctx := context.Background()
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader for GCS object: %w", err)
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	var weatherData map[string]interface{}
	err = json.Unmarshal(data, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	return weatherData, nil
}

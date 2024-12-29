# Weather App Documentation

## Setup and Installation

### Prerequisites

- **Go Version**: `go1.23` or higher
- **Environment Variable**: Set up the required GCP credentials.

#### Setting the Environment Variable

**Windows (PowerShell)**:

```powershell
$env:GOOGLE_APPLICATION_CREDENTIALS="<path-to-secret>\meta-me-credentials.json"
```

#### Running the server locally

```
go run cmd/main.go
```

The server will listen for requests on port 8080 (default).

---

## Deployment Steps for Google Cloud Run

1. **Build and Push Docker Image**  
   Build the Docker image using the Dockerfile and push it to the Google Artifact Registry.

2. **Enter Container Image URL**  
   Provide the URL of the container image in the Cloud Run setup.

3. **Allow Unauthenticated Invocations**  
   Select the option to **Allow unauthenticated invocations** to enable public access to the service.

4. **Note the Generated URL**  
   After deployment, note the Cloud Run service URL:  
   `https://weather-app-754744729721.us-central1.run.app`

# API Documentation

This document provides the details of the implemented API endpoints, including endpoint URLs, methods, input/output formats, and required parameters.

---

## Endpoints

### 1. Store Weather Data

- **URL**: `https://weather-app-754744729721.us-central1.run.app/store-weather-data`
- **Method**: POST
- **Description**: Fetches weather data for a specified location and date range and uploads it to a Google Cloud Storage (GCS) bucket.
- **Request Body**:

  ```json
  {
    "latitude": 23.46,
    "longitude": 44.499,
    "start_date": "2024-12-25",
    "end_date": "2024-12-26"
  }
  ```

  - **latitude**: Latitude of the location (float).
  - **longitude**: Longitude of the location (float).
  - **start_date**: Start date in `YYYY-MM-DD` format.
  - **end_date**: End date in `YYYY-MM-DD` format.

- **Response**:

  - On success: Returns a status message indicating the file was uploaded to the GCS bucket.
  - On error: Returns an error message describing the failure.

- **cURL Command**:
  ```bash
  curl --location --request POST 'https://weather-app-754744729721.us-central1.run.app/store-weather-data' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "latitude": 23.46,
      "longitude": 44.499,
      "start_date": "2024-12-25",
      "end_date": "2024-12-26"
  }'
  ```

---

### 2. List Weather Files

- **URL**: `https://weather-app-754744729721.us-central1.run.app/list-weather-files`
- **Method**: GET
- **Description**: Lists all the weather data files stored in the GCS bucket.

- **Response**:

  - A JSON array containing the names of the weather data files.
    ```json
    ["response-2024-12-25-14-30-45.json", "response-2024-12-26-15-40-10.json"]
    ```
  - On error: Returns an error message describing the failure.

- **cURL Command**:
  ```bash
  curl --location --request GET 'https://weather-app-754744729721.us-central1.run.app/list-weather-files'
  ```

---

### 3. Get Weather File Content

- **URL**: `https://weather-app-754744729721.us-central1.run.app/weather-file-content/{fileName}`
- **Method**: GET
- **Description**: Fetches and displays the content of a specific weather data file stored in the GCS bucket.
- **Path Parameter**:

  - `fileName`: The name of the file to retrieve (e.g., `response-2024-12-29-19-38-44.json`).

- **Response**:

  - The content of the specified weather data file in JSON format.
    ```json
    {
      "latitude": 23.44464,
      "longitude": 44.479496,
      "generationtime_ms": 0.10704994201660156,
      "utc_offset_seconds": 0,
      "timezone": "GMT",
      "timezone_abbreviation": "GMT",
      "elevation": 933.0,
      "daily_units": {
        "time": "iso8601",
        "temperature_2m_max": "°C",
        "temperature_2m_min": "°C",
        "temperature_2m_mean": "°C",
        "apparent_temperature_max": "°C",
        "apparent_temperature_min": "°C",
        "apparent_temperature_mean": "°C"
      },
      "daily": {
        "time": ["2024-12-25", "2024-12-26"],
        "temperature_2m_max": [23.1, 23.1],
        "temperature_2m_min": [12.6, 12.2],
        "temperature_2m_mean": [17.3, 17.3],
        "apparent_temperature_max": [18.7, 17.3],
        "apparent_temperature_min": [8.4, 7.3],
        "apparent_temperature_mean": [13.2, 12.1]
      }
    }
    ```
  - On error: Returns an error message describing the failure.

- **cURL Command**:
  ```bash
  curl --location --request GET 'https://weather-app-754744729721.us-central1.run.app/weather-file-content/response-2024-12-29-19-38-44.json'
  ```

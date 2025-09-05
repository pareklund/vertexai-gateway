# Vertex AI Gateway

A minimal Go application that acts as an inference gateway to Google Cloud Vertex AI.

## Features

- REST API gateway for Vertex AI text generation
- Built with Gin web framework
- Configurable model parameters
- Health check endpoint
- Proper error handling and logging

## Requirements

- Go 1.25.1 or later
- Google Cloud Project with Vertex AI API enabled
- Service account credentials or Application Default Credentials

## Environment Variables

- `GOOGLE_CLOUD_PROJECT`: Your Google Cloud project ID (required)
- `GOOGLE_CLOUD_LOCATION`: Vertex AI location (default: us-central1)
- `GOOGLE_APPLICATION_CREDENTIALS`: Path to service account key file (optional if using ADC)
- `PORT`: Server port (default: 8080)

## Installation

1. Clone or download the application
2. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

1. Set your environment variables:
   ```bash
   export GOOGLE_CLOUD_PROJECT="your-project-id"
   export GOOGLE_APPLICATION_CREDENTIALS="path/to/service-account.json"
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

## API Endpoints

### Health Check
```
GET /health
```

### Inference
```
POST /v1/inference
```

**Request Body:**
```json
{
  "prompt": "Your text prompt here",
  "model": "text-bison@001",
  "temperature": 0.2,
  "max_tokens": 256,
  "parameters": {
    "topP": 0.8,
    "topK": 40
  }
}
```

**Response:**
```json
{
  "text": "Generated response text",
  "model": "text-bison@001",
  "usage": {
    "prompt_tokens": 15,
    "completion_tokens": 42,
    "total_tokens": 57
  }
}
```

## Example

```bash
curl -X POST http://localhost:8080/v1/inference \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Write a short poem about artificial intelligence",
    "temperature": 0.7,
    "max_tokens": 100
  }'
```
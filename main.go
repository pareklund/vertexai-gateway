package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/genai"

	"github.com/joho/godotenv"
)

type Gateway struct {
	client    *genai.Client
	projectID string
	location  string
}

func NewGateway() (*Gateway, error) {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable is required")
	}

	location := os.Getenv("GOOGLE_CLOUD_LOCATION")
	if location == "" {
		location = "us-central1" // Default location
	}

	var opts []option.ClientOption
	if keyFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); keyFile != "" {
		opts = append(opts, option.WithCredentialsFile(keyFile))
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:      os.Getenv("GEMINI_API_KEY"),
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Gateway{
		client:    client,
		projectID: projectID,
		location:  location,
	}, nil
}

type InferenceRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int32   `json:"max_tokens"`
}

func (g *Gateway) Close() error {
	return nil
}

func (g *Gateway) handleInference(c *gin.Context) {
	ctx := context.Background()

	var req InferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := g.client.Models.GenerateContent(ctx,
		req.Model,
		genai.Text(req.Prompt),
		&genai.GenerateContentConfig{
			Temperature:     &req.Temperature,
			MaxOutputTokens: req.MaxTokens,
		},
	)
	if err != nil {
		c.Error(fmt.Errorf("failed to generate content: %w", err))
	}

	respText := resp.Text()
	c.JSON(http.StatusOK, gin.H{
		"text": respText,
	})
}

func (g *Gateway) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "vertexai-gateway",
	})
}

func main() {
	_ = godotenv.Load()
	gateway, err := NewGateway()
	if err != nil {
		log.Fatalf("Failed to initialize gateway: %v", err)
	}
	defer gateway.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/health", gateway.healthCheck)
	router.POST("/v1/inference", gateway.handleInference)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Vertex AI Gateway on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

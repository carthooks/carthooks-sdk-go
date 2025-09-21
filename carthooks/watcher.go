package carthooks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// WatcherConfig holds configuration for the watcher
type WatcherConfig struct {
	Client       *Client
	WatcherID    string
	AppID        uint
	CollectionID uint
	SQSQueueURL  string
	AWSRegion    string
	Filters      map[string]interface{}
	Handler      func(ctx interface{}, record map[string]interface{})
}

// Watcher represents a data change watcher
type Watcher struct {
	config    *WatcherConfig
	sqsClient *sqs.Client
	running   bool
	stopChan  chan bool
}

// SQSMessageBody represents the expected SQS message structure
type SQSMessageBody struct {
	Meta    map[string]interface{} `json:"meta"`
	Payload map[string]interface{} `json:"payload"`
	Version string                 `json:"version"`
}

// NewWatcher creates a new watcher instance
func NewWatcher(config *WatcherConfig) (*Watcher, error) {
	// Load AWS configuration
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(config.AWSRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	return &Watcher{
		config:    config,
		sqsClient: sqsClient,
		running:   false,
		stopChan:  make(chan bool),
	}, nil
}

// Subscribe sets up the watch data subscription
func (w *Watcher) Subscribe() error {
	// Start watch data
	watchName := fmt.Sprintf("watch-%d-%d", w.config.AppID, w.config.CollectionID)

	options := &WatchDataOptions{
		EndpointURL:    w.config.SQSQueueURL,
		EndpointType:   "sqs",
		Name:           watchName,
		AppID:          w.config.AppID,
		CollectionID:   w.config.CollectionID,
		Filters:        w.config.Filters,
		Age:            432000, // 5 days
		WatchStartTime: 0,
	}

	result := w.config.Client.StartWatchData(options)
	if !result.Success {
		return fmt.Errorf("failed to start watch data: %s", result.Error)
	}

	log.Printf("✅ Monitoring task registered successfully: %s", watchName)
	return nil
}

// Run starts the watcher and begins listening for messages
func (w *Watcher) Run() error {
	if w.running {
		return fmt.Errorf("watcher is already running")
	}

	// Subscribe to watch data
	if err := w.Subscribe(); err != nil {
		return err
	}

	w.running = true
	log.Printf("🎯 SQS mode running...")

	// Start SQS message polling
	go w.pollSQSMessages()

	// Wait for stop signal
	<-w.stopChan
	w.running = false
	log.Printf("🛑 Watcher stopped")

	return nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	if w.running {
		w.stopChan <- true
	}
}

// pollSQSMessages continuously polls SQS for messages
func (w *Watcher) pollSQSMessages() {
	for w.running {
		// Receive messages from SQS
		result, err := w.sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(w.config.SQSQueueURL),
			MaxNumberOfMessages: 5,
			VisibilityTimeout:   300, // 5 minutes
			WaitTimeSeconds:     20,  // Long polling
		})

		if err != nil {
			log.Printf("❌ Error receiving SQS messages: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Process each message
		for _, message := range result.Messages {
			if err := w.processMessage(message); err != nil {
				log.Printf("⚠️ Message processing failed: %v", err)
				continue
			}

			// Delete message after successful processing
			_, err := w.sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(w.config.SQSQueueURL),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("⚠️ Failed to delete message: %v", err)
			}
		}

		// Short sleep to prevent excessive polling
		if len(result.Messages) == 0 {
			time.Sleep(1 * time.Second)
		}
	}
}

// processMessage processes a single SQS message
func (w *Watcher) processMessage(message types.Message) error {
	if message.Body == nil {
		return fmt.Errorf("message body is nil")
	}

	// Parse message body
	var messageBody SQSMessageBody
	if err := json.Unmarshal([]byte(*message.Body), &messageBody); err != nil {
		return fmt.Errorf("failed to parse message body: %w", err)
	}

	// Validate message format
	if messageBody.Payload == nil {
		return fmt.Errorf("message payload is nil")
	}

	// Check if payload has ID
	if _, exists := messageBody.Payload["id"]; !exists {
		return fmt.Errorf("incorrect message format, missing payload.id")
	}

	// Call user handler
	if w.config.Handler != nil {
		w.config.Handler(nil, messageBody.Payload)
	}

	return nil
}

// WatcherBuilder provides a fluent interface for building watchers
type WatcherBuilder struct {
	config *WatcherConfig
}

// NewWatcherBuilder creates a new watcher builder
func NewWatcherBuilder(client *Client, watcherID string) *WatcherBuilder {
	return &WatcherBuilder{
		config: &WatcherConfig{
			Client:    client,
			WatcherID: watcherID,
			AWSRegion: "ap-southeast-1", // Default region
		},
	}
}

// WithApp sets the app and collection IDs
func (wb *WatcherBuilder) WithApp(appID, collectionID uint) *WatcherBuilder {
	wb.config.AppID = appID
	wb.config.CollectionID = collectionID
	return wb
}

// WithSQS sets the SQS configuration
func (wb *WatcherBuilder) WithSQS(queueURL, region string) *WatcherBuilder {
	wb.config.SQSQueueURL = queueURL
	wb.config.AWSRegion = region
	return wb
}

// WithFilters sets the data filters
func (wb *WatcherBuilder) WithFilters(filters map[string]interface{}) *WatcherBuilder {
	wb.config.Filters = filters
	return wb
}

// WithHandler sets the message handler
func (wb *WatcherBuilder) WithHandler(handler func(ctx interface{}, record map[string]interface{})) *WatcherBuilder {
	wb.config.Handler = handler
	return wb
}

// Build creates the watcher
func (wb *WatcherBuilder) Build() (*Watcher, error) {
	return NewWatcher(wb.config)
}

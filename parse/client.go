package parse

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("vim-go")
}

// ParseClient represents a client that communicates with the Parse Server.
type ParseClient struct {
	// BaseURL is the base url of the Parse Server.
	BaseURL string
	// AppID is the application id of the Parse Server.
	AppID string
	// APIKey is the REST API key of the Parse Server.
	APIKey string
	// HTTPClient is the underlying HTTP client that performs the requests.
	HTTPClient *http.Client
	// Database is the database client that performs the database operations.
	DB interface{}
}

// NewParseClient creates a new ParseClient with the given base url, app id and api key.
func NewParseClient(baseURL, appID, apiKey string, db interface{}) *ParseClient {
	return &ParseClient{
		BaseURL:    baseURL,
		AppID:      appID,
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
		DB:         db,
	}
}

// Do sends an HTTP request and returns an HTTP response.
func (c *ParseClient) Do(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

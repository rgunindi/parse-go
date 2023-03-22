package parse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ParseQuery represents a query object on the Parse Server.
type ParseQuery struct {
	// ClassName is the name of the class to query.
	ClassName string
	// Where is a map of constraints for the query.
	Where map[string]interface{}
	// Limit is the maximum number of results to return.
	Limit int
	// Skip is the number of results to skip before returning any results.
	Skip int
	// Order is a comma-separated list of field names that specify the sorting order of the results.
	Order string
	// Keys is a comma-separated list of field names that specify which fields to include in the returned objects.
	Keys string
	// Include is a comma-separated list of related objects to include in the returned objects.
	Include string
}

// NewParseQuery creates a new ParseQuery with the given class name and an empty Where map.
func NewParseQuery(className string) *ParseQuery {
	return &ParseQuery{
		ClassName: className,
		Where:     make(map[string]interface{}),
	}
}

// Find finds all objects that match this query and returns a slice of ParseObjects.
func (q *ParseQuery) Find(client *ParseClient) ([]*ParseObjectServer, error) {
	// Create a new HTTP request with the query url and GET method
	url := client.BaseURL + "/classes/" + q.ClassName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Add the required headers to the request
	req.Header.Add("X-Parse-Application-Id", client.AppID)
	req.Header.Add("X-Parse-REST-API-Key", client.APIKey)
	// Add the query parameters to the request
	qs := req.URL.Query()
	if len(q.Where) > 0 {
		where, err := json.Marshal(q.Where)
		if err != nil {
			return nil, err
		}
		qs.Add("where", string(where))
	}
	if q.Limit > 0 {
		qs.Add("limit", strconv.Itoa(q.Limit))
	}
	if q.Skip > 0 {
		qs.Add("skip", strconv.Itoa(q.Skip))
	}
	if q.Order != "" {
		qs.Add("order", q.Order)
	}
	if q.Keys != "" {
		qs.Add("keys", q.Keys)
	}
	if q.Include != "" {
		qs.Add("include", q.Include)
	}
	req.URL.RawQuery = qs.Encode()
	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		// If the response is successful, create a slice of ParseObjects from the result
		var objects []*ParseObjectServer
		if results, ok := result["results"].([]interface{}); ok {
			for _, r := range results {
				if data, ok := r.(map[string]interface{}); ok {
					obj := NewParseObject(q.ClassName, data)
					objects = append(objects, obj)
				}
			}
		}
		return objects, nil
	} else {
		// Otherwise return an error with the message from the server
		return nil, fmt.Errorf("%v: %v", resp.Status, result["error"])
	}
}

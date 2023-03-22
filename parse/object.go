package parse

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// ParseObject interface defines the methods for accessing and manipulating Parse objects
type ParseObject interface {
	GetObjectID() string
	GetClassName() string
	GetData() map[string]interface{}
	Save(client interface{}) error
	Delete(client interface{}) error
}

// ParseObject represents a data object on the Parse Server.
type ParseObjectServer struct {
	// ClassName is the name of the class that this object belongs to.
	ClassName string `json:"className,omitempty"`
	// ObjectID is the unique identifier of this object on the server.
	ObjectID string `json:"objectId,omitempty"`
	// CreatedAt is the time when this object was created on the server.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt is the time when this object was last updated on the server.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// ACL is the access control list that determines who can read or write this object.
	ACL map[string]interface{} `json:"ACL,omitempty"`
	// Data is a map that holds the custom fields and values of this object.
	Data map[string]interface{} `json:"-"`
}

// NewParseObject creates a new ParseObject with the given class name and data.
func NewParseObject(className string, data map[string]interface{}) *ParseObjectServer {
	return &ParseObjectServer{
		ClassName: className,
		Data:      data,
	}
}

// GetObjectID returns the object id of this object
func (p *ParseObjectServer) GetObjectID() string {
	return p.ObjectID
}

// GetClassName returns the class name of this object
func (p *ParseObjectServer) GetClassName() string {
	return p.ClassName
}

// GetData returns the data of this object as a map
func (p *ParseObjectServer) GetData() map[string]interface{} {
	return p.Data
}

// Save with interface
func (p *ParseObjectServer) Save(client interface{}) error {
	return p.Save(client)
}

// Delete deletes this object from the server.
func (p *ParseObjectServer) Delete(client interface{}) error {
	return p.Delete(client)
}

//// Save inserts or updates this object in MongoDB database using the given client
//func (p *ParseObjectServer) Save(client *mongo.Client) error {
//	return p.Save(client)
//}
//
//// Save inserts or updates this object in Redis database using the given client
//func (p *ParseObjectServer) Save(client *redis.Client) error {
//	return p.Save(client)
//}

// Save saves a given ParseObject to the database using the client's DB field
func (c *ParseClient) Save(obj *ParseObjectServer) error {
	// Check the type of the DB field and call the corresponding Save method
	switch db := c.DB.(type) {
	case *mongo.Client:
		return obj.Save(db) // Call the Save method for MongoDB
	case *redis.Client:
		return obj.Save(db) // Call the Save method for Redis
	default:
		return fmt.Errorf("Unsupported database type: %T", db) // Return an error if the database type is not supported
	}
}
func (c *ParseClient) Delete(obj *ParseObjectServer) error {
	// Check the type of the DB field and call the corresponding Save method
	switch db := c.DB.(type) {
	case *mongo.Client:
		return obj.Delete(db) // Call the Save method for MongoDB
	case *redis.Client:
		return obj.Delete(db) // Call the Save method for Redis
	default:
		return fmt.Errorf("Unsupported database type: %T", db) // Return an error if the database type is not supported
	}
}

// Save saves this object to the server and updates its fields with the server response.
//func (p *ParseObject) Save(client *ParseClient) error {
//	// Create a JSON body with the data of this object
//	body, err := json.Marshal(p.Data)
//	if err != nil {
//		return err
//	}
//	var url string
//	var method string
//	// If this object has an ObjectID, use PUT method and append the ObjectID to the url
//	if p.ObjectID != "" {
//		url = client.BaseURL + "/classes/" + p.ClassName + "/" + p.ObjectID
//		method = "PUT"
//	} else {
//		// Otherwise, use POST method and generate a new ObjectID on the server
//		url = client.BaseURL + "/classes/" + p.ClassName
//		method = "POST"
//	}
//	// Create a new HTTP request with the given url, method and body
//	req, err := http.NewRequest(method, url, bytes.NewReader(body))
//	if err != nil {
//		return err
//	}
//	// Add the required headers to the request
//	req.Header.Add("X-Parse-Application-Id", client.AppID)
//	req.Header.Add("X-Parse-REST-API-Key", client.APIKey)
//	req.Header.Add("Content-Type", "application/json")
//	// Send the request and get the response
//	resp, err := client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	var result map[string]interface{}
//	err = json.NewDecoder(resp.Body).Decode(&result)
//	if err != nil {
//		return err
//	}
//	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
//		// If the response is successful, update this object's fields with the result
//		if oid, ok := result["objectId"].(string); ok {
//			p.ObjectID = oid
//		}
//		if createdAtStr, ok := result["createdAt"].(string); ok {
//			p.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAtStr)
//		}
//		if updatedAtStr, ok := result["updatedAt"].(string); ok {
//			p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updatedAtStr)
//		}
//	} else {
//		// Otherwise return an error with the message from the server
//		return fmt.Errorf("%v: %v", resp.Status, result["error"])
//	}
//	return nil
//}

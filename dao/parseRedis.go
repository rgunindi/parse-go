package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

type ParseObjectRedis struct {
	ObjectID  string                 `json:"objectId"`
	ClassName string                 `json:"className"`
	Data      map[string]interface{} `json:"data"`
}

// GetObjectID returns the object id of this object
func (p *ParseObjectRedis) GetObjectID() string {
	return p.ObjectID
}

// GetClassName returns the class name of this object
func (p *ParseObjectRedis) GetClassName() string {
	return p.ClassName
}

// GetData returns the data of this object as a map
func (p *ParseObjectRedis) GetData() map[string]interface{} {
	return p.Data
}

// Save inserts or updates this object in Redis database using the given client
func (p *ParseObjectRedis) Save(client *redis.Client) error {
	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if p.ObjectID == "" {
		// If this object has no ObjectID, generate a new one and insert it into the database with a unique key
		p.ObjectID = uuid.New().String()
		key := p.ClassName + ":" + p.ObjectID  // The key format is ClassName:ObjectId
		err = client.Set(ctx, key, p, 0).Err() // Set the value to the JSON representation of this object with no expiration
	} else {
		// Otherwise, update the existing value with this ObjectID in the database
		key := p.ClassName + ":" + p.ObjectID  // The key format is ClassName:ObjectId
		err = client.Set(ctx, key, p, 0).Err() // Set the value to the JSON representation of this object with no expiration
	}
	return err // Return any error that occurred during the operation

}

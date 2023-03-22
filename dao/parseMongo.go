package dao

import (
	"context"
	"github.com/rgunindi/parse-go/parse"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ParseClient struct {
	ParseClient *parse.ParseClient
}

// ParseObjectMongoDB is a concrete implementation of ParseObject interface for MongoDB database
type ParseObjectMongoDB struct {
	ObjectID  string                 `bson:"_id"`
	ClassName string                 `bson:"className"`
	Data      map[string]interface{} `bson:"data"`
}

// GetObjectID returns the object id of this object
func (p *ParseObjectMongoDB) GetObjectID() string {
	return p.ObjectID
}

// GetClassName returns the class name of this object
func (p *ParseObjectMongoDB) GetClassName() string {
	return p.ClassName
}

// GetData returns the data of this object as a map
func (p *ParseObjectMongoDB) GetData() map[string]interface{} {
	return p.Data
}

// Save inserts or updates this object in MongoDB database using the given client
func (p *ParseObjectMongoDB) Save(client *mongo.Client) error {
	// Get the collection for this class name from the client's database
	collection := client.Database("parse").Collection(p.ClassName)
	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	if p.ObjectID == "" {
		// If this object has no ObjectID, generate a new one and insert it into the collection
		p.ObjectID = primitive.NewObjectID().Hex()
		_, err = collection.InsertOne(ctx, p)
	} else {
		// Otherwise, update the existing document with this ObjectID in the collection
		filter := bson.M{"_id": p.ObjectID}
		update := bson.M{"$set": bson.M{"data": p.Data}}
		_, err = collection.UpdateOne(ctx, filter, update)
	}
	return err // Return any error that occurred during the operation

}

// Delete deletes this object from MongoDB database using the given client
func (p *ParseObjectMongoDB) Delete(client *mongo.Client) error {
	// Get the collection for this class name from the client's database
	collection := client.Database("parse").Collection(p.ClassName)
	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Delete the document with this ObjectID from the collection
	filter := bson.M{"_id": p.ObjectID}
	_, err := collection.DeleteOne(ctx, filter)
	return err // Return any error that occurred during the operation
}

package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ParseUser represents a user object on the Parse Server.
type ParseUser struct {
	// Embed the ParseObject type to inherit its fields and methods
	*ParseObjectServer
	// Username is the unique username of this user.
	Username string `json:"username,omitempty"`
	// Password is the password of this user.
	Password string `json:"password,omitempty"`
	// Email is the email address of this user.
	Email string `json:"email,omitempty"`
	// SessionToken is the session token of this user after logging in or signing up.
	SessionToken string `json:"sessionToken,omitempty"`
}

// NewParseUser creates a new ParseUser with the given username, password and email.
func NewParseUser(username, password, email string) *ParseUser {
	return &ParseUser{
		ParseObjectServer: NewParseObject("_User", nil),
		Username:          username,
		Password:          password,
		Email:             email,
	}
}

// SignUp signs up this user to the server and updates its fields with the server response.
func (u *ParseUser) SignUp(client *ParseClient) error {
	// Create a JSON body with the data of this user
	body, err := json.Marshal(u)
	if err != nil {
		return err
	}
	// Create a new HTTP request with the sign up url and POST method
	url := client.BaseURL + "/users"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	// Add the required headers to the request
	req.Header.Add("X-Parse-Application-Id", client.AppID)
	req.Header.Add("X-Parse-REST-API-Key", client.APIKey)
	req.Header.Add("Content-Type", "application/json")
	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusCreated {
		// If the response is successful, update this user's fields with the result
		if oid, ok := result["objectId"].(string); ok {
			u.ObjectID = oid
		}
		if createdAtStr, ok := result["createdAt"].(string); ok {
			u.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAtStr)
		}
		if sessionToken, ok := result["sessionToken"].(string); ok {
			u.SessionToken = sessionToken
		}
	} else {
		// Otherwise return an error with the message from the server
		return fmt.Errorf("%v: %v", resp.Status, result["error"])
	}
	return nil
}

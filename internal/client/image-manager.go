package client

// Interface defines the methods
type ImageManagerClient interface {
	ValidateImageID(imageId string) error
}

// Concrete struct
type imageManagerClient struct {
	// fields like http.Client, config, etc.
	serviceName string
}

// Constructor
func NewImageManagerClient() ImageManagerClient {
	return &imageManagerClient{}
}

// Implementation of interface
func (c *imageManagerClient) ValidateImageID(imageId string) error {
	// Logic to validate image
	return nil
}

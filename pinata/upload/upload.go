// Package upload provides functionality for uploading content to Pinata
package upload

// Service provides upload operations for Pinata
type Service struct {
	config interface{}
	Public *PublicService
	Private *PrivateService
}

// New creates a new upload service with the provided configuration
func New(config interface{}) *Service {
	service := &Service{
		config: config,
	}

	// Initialize public and private services
	service.Public = NewPublicService(config)
	service.Private = NewPrivateService(config)

	return service
}

// Config returns the service configuration
func (s *Service) Config() interface{} {
	return s.config
}

package files

// Service provides access to files operations
type Service struct {
	Config  interface{}
	Public  *PublicService
	Private *PrivateService
}

// New creates a new files service
func New(config interface{}) *Service {
	service := &Service{
		Config: config,
	}

	service.Public = NewPublicService(config)
	service.Private = NewPrivateService(config)

	return service
}

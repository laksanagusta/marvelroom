package cryptography

// CryptoServiceProvider provides crypto service instances
type ServiceProvider struct{}

// NewServiceProvider creates a new CryptoServiceProvider
func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

// NewCryptoService creates a new crypto service instance
func (p *ServiceProvider) NewCryptoService() (Service, error) {
	return NewCryptoService()
}

// NewCryptoServiceWithKeys creates a new crypto service instance with existing keys
func (p *ServiceProvider) NewCryptoServiceWithKeys(publicKey, privateKey []byte) (Service, error) {
	return NewCryptoServiceWithKeys(publicKey, privateKey)
}

// Ensure ServiceProvider implements the interface
var _ CryptoServiceProvider = (*ServiceProvider)(nil)

// CryptoServiceProvider defines the interface for creating crypto service instances
type CryptoServiceProvider interface {
	NewCryptoService() (Service, error)
	NewCryptoServiceWithKeys(publicKey, privateKey []byte) (Service, error)
}
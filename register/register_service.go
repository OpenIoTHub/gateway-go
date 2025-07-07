package register

import (
	"github.com/OpenIoTHub/utils/models"
	"net"
	"sync"
)

var registeredServices = make([]models.MDNSResult, 0)
var registeredServicesLock sync.RWMutex

func RegisterService(instance, service, domain, hostname string, port int, text []string, TTL uint32, AddrIPv4, AddrIPv6 []net.IP) (err error) {
	registeredServicesLock.Lock()
	defer registeredServicesLock.Unlock()
	registeredServices = append(registeredServices, models.MDNSResult{
		Instance: instance,
		Service:  service,
		Domain:   domain,
		HostName: hostname,
		Port:     port,
		Text:     text,
		TTL:      TTL,
		AddrIPv4: AddrIPv4,
		AddrIPv6: AddrIPv6,
	})
	return
}

func GetRegisteredServices() (services []models.MDNSResult) {
	registeredServicesLock.RLock()
	defer registeredServicesLock.RUnlock()
	return registeredServices
}

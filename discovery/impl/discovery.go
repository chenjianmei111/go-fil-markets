package discoveryimpl

import (
	"github.com/chenjianmei111/go-fil-markets/discovery"
)

func Multi(r discovery.PeerResolver) discovery.PeerResolver { // TODO: actually support multiple mechanisms
	return r
}

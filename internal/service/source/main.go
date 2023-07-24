package source

import (
	"github.com/opal-project-dev/oracle/internal/service/types"
)

type Source interface {
	GetCurrentPrice(priceChan chan<- types.Data) error
}

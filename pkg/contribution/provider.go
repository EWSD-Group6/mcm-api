package contribution

import "github.com/google/wire"

var Set = wire.NewSet(InitializeRepository, InitializeService)

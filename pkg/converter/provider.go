package converter

import "github.com/google/wire"

var Set = wire.NewSet(NewGotenbergDocumentConverter)

package radx

// block size (doesn't fit in 1024 because 256*4=1024 but we also need block header)
const RadxBlockSize = 2048

const (
	RadxTypeRoot    uint8 = iota // root node (which can only be a full index)
	RadxTypeFull          = iota // full index
	RadxTypeCompact       = iota // compact index
	RadxTypeEdge          = iota
)

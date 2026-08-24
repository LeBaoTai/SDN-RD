package builder

import "github.com/LeBaoTai/myco-controller/internal/oc"

type Builder struct {
	device *oc.Device
}

func NewBuilder() *Builder {
	return &Builder{
		device: &oc.Device{},
	}
}

func (b *Builder) Build() *oc.Device {
	return b.device
}

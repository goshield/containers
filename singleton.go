package containers

var (
	c Container
)

func ContainerInstance() Container {
	if c == nil {
		c = NewContainer()
	}

	return c
}

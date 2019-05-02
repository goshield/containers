package containers_test

import (
	"testing"
)

func TestContainers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Containers Suite")
}

package containers

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container", func() {
	It("NewContainer should return Container", func() {
		Expect(NewContainer()).NotTo(BeNil())
	})
})

type sampleInvalidError struct{}
type facErr struct {
	code string
	msg  string
}
type InjectErrorResolveErrorNotExistAbstract struct {
	Err error `inject:"*"`
}
type InjectOk struct {
	Err error `inject:"*"`
}
type InjectRecursiveOk struct {
	Err                          error       `inject:"*"`
	Foo                          InjectFooer `inject:"*"`
	NotInjectableNonInterface    string      `inject:"*"`
	notInjectablePrivateProperty error       `inject:"*"`
}
type InjectFooer interface {
	Foo() string
}
type InjectFoo struct {
	Baz InjectBazer `inject:"*"`
}

func (f *InjectFoo) Foo() string { return "Foo.." }

type InjectBazer interface {
	Baz() string
}
type InjectBaz struct{}

func (b *InjectBaz) Baz() string { return "Baz.." }

var _ = Describe("FactoryContainer", func() {
	It("Bind should return error code ERR_BIND_INVALID_ARGUMENTS", func() {
		c := NewContainer()
		err := c.Bind("string", errors.New("some-error"))
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(ERR_BIND_INVALID_ARGUMENTS))
	})

	It("Bind should return error code ERR_BIND_NOT_IMPLEMENT_INTERFACE", func() {
		c := NewContainer()
		err := c.Bind((*InjectFooer)(nil), errors.New("some-error"))
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(ERR_BIND_NOT_IMPLEMENT_INTERFACE, "*errors.errorString", "containers.InjectFooer")))
	})

	//It("Bind should return error code ERR_BIND_INVALID_CONCRETE", func() {
	//	c := NewContainer()
	//	a := make([]interface{}, 2)
	//	a[0] = "a_string"
	//	a[1] = sampleInvalidError{}
	//	for _, i := range a {
	//		err := c.Bind((*InjectFooer)(nil), i)
	//		Expect(err).NotTo(BeNil())
	//		Expect(err.Error()).To(Equal(fmt.Sprintf(ERR_BIND_INVALID_CONCRETE)))
	//	}
	//})

	It("Bind should return nil when binding struct", func() {
		c := NewContainer()
		err := c.Bind(InjectBaz{}, &InjectBaz{})
		Expect(err).To(BeNil())
	})

	It("Bind should return error code ERR_BIND_INVALID_STRUCT", func() {
		c := NewContainer()
		err := c.Bind(InjectBaz{}, "not_a_struct")
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(ERR_BIND_INVALID_STRUCT))
	})

	It("Bind should return error code ERR_BIND_INVALID_STRUCT_CONCRETE", func() {
		c := NewContainer()
		err := c.Bind(InjectBaz{}, errors.New("some-error"))
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(ERR_BIND_INVALID_STRUCT_CONCRETE, "containers.InjectBaz", "errors.errorString")))
	})

	It("Bind should return nil when binding function", func() {
		c := NewContainer()
		err := c.Bind((*error)(nil), func(message string) error {
			return errors.New(message)
		})
		Expect(err).To(BeNil())
	})

	It("Bind should return nil when binding a pointer", func() {
		c := NewContainer()
		err := c.Bind((*error)(nil), errors.New("my_message"))
		Expect(err).To(BeNil())
	})

	It("Resolve should return error code ERR_RESOLVE_NOT_EXIST_ABSTRACT", func() {
		c := NewContainer()
		_, err := c.Resolve(&InjectBaz{})
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(ERR_RESOLVE_NOT_EXIST_ABSTRACT, "containers.InjectBaz")))
	})

	It("Resolve should return error code ERR_RESOLVE_INSUFFICIENT_ARGUMENTS", func() {
		c := NewContainer()
		_ = c.Bind((*error)(nil), func(message string) error {
			return errors.New(message)
		})
		_, err := c.Resolve((*error)(nil))
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(ERR_RESOLVE_INSUFFICIENT_ARGUMENTS, 1, 0)))
	})

	It("Resolve should return error code ERR_RESOLVE_NON_VALUES_RETURNED", func() {
		c := NewContainer()
		_ = c.Bind((*error)(nil), func(message string) {})
		_, err := c.Resolve((*error)(nil), "my_message")
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(ERR_RESOLVE_NON_VALUES_RETURNED))
	})

	It("Resolve should return nil", func() {
		c := NewContainer()
		_ = c.Bind((*error)(nil), func(message string) error {
			return errors.New(message)
		})
		e, err := c.Resolve((*error)(nil), "my_message")
		Expect(err).To(BeNil())
		Expect(e.(error).Error()).To(Equal("my_message"))
	})

	It("Resolve should return nil when resolving pointer", func() {
		c := NewContainer()
		_ = c.Bind((*error)(nil), errors.New("my_message"))
		e, err := c.Resolve((*error)(nil))
		Expect(err).To(BeNil())
		Expect(e.(error).Error()).To(Equal("my_message"))
	})

	It("Resolve should return nil when resolving struct", func() {
		c := NewContainer()
		_ = c.Bind(InjectBaz{}, &InjectBaz{})
		v, err := c.Resolve(InjectBaz{})
		Expect(err).To(BeNil())
		Expect(v).NotTo(BeNil())
	})

	It("Resolve should return error code ERR_RESOLVE_INVALID_ARGUMENTS", func() {
		c := NewContainer()
		_, err := c.Resolve("string")
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(ERR_RESOLVE_INVALID_ARGUMENTS))
	})

	It("Inject should return error code ERR_INJECT_INVALID_TARGET_TYPE", func() {
		c := NewContainer()
		e := facErr{"my_code", "my_message"}
		err := c.Inject(e)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(ERR_INJECT_INVALID_TARGET_TYPE, "struct")))
	})

	It("Inject should return error code ERR_RESOLVE_NOT_EXIST_ABSTRACT", func() {
		c := NewContainer()
		in := &InjectErrorResolveErrorNotExistAbstract{}
		err := c.Inject(in)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal(fmt.Sprintf(ERR_RESOLVE_NOT_EXIST_ABSTRACT, "error")))
	})

	It("Inject should return nil", func() {
		c := NewContainer()
		_ = c.Bind((*error)(nil), errors.New("my_message"))
		in := &InjectOk{}
		err := c.Inject(in)
		Expect(err).To(BeNil())
		Expect(in.Err.Error()).To(Equal("my_message"))
	})

	It("Inject should do a recursive injection", func() {
		// InjectRecursiveOk -> Err
		// InjectRecursiveOk -> Foo -> Baz
		c := NewContainer()
		_ = c.Bind((*error)(nil), errors.New("my_message"))
		_ = c.Bind((*InjectFooer)(nil), &InjectFoo{})
		_ = c.Bind((*InjectBazer)(nil), &InjectBaz{})
		in := &InjectRecursiveOk{}
		err := c.Inject(in)
		Expect(err).To(BeNil())

		// Asserting for in.Err
		Expect(in.Err.Error()).To(Equal("my_message"))

		// Asserting for in.Foo
		Expect(in.Foo.Foo()).To(Equal("Foo.."))
		Expect(in.Foo.(*InjectFoo).Baz.Baz()).To(Equal("Baz.."))
	})
})

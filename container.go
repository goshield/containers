package containers

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Container acts as a dependency-injection manager
type Container interface {
	// Bind stores a concrete of an abstract, as default sharing is enable
	Bind(abstract interface{}, concrete interface{}) error

	// Resolve processes and returns a concrete of proposed abstract
	Resolve(abstract interface{}, args ...interface{}) (concrete interface{}, err error)

	// Inject resolves target's dependencies
	Inject(target interface{}) error
}

func NewContainer() Container {
	return &factoryContainer{items: new(sync.Map)}
}

type factoryContainer struct {
	items *sync.Map
}

func (c *factoryContainer) Bind(abstract interface{}, concrete interface{}) error {
	at, isInterface := c.interfaceOf(abstract)
	if isInterface == nil {
		return c.bindInterface(at, concrete)
	}

	at, isStruct := c.structOf(abstract)
	if isStruct == nil {
		return c.bindStruct(at, concrete)
	}

	return errors.New(ERR_BIND_INVALID_ARGUMENTS)
}

func (c *factoryContainer) Resolve(abstract interface{}, args ...interface{}) (concrete interface{}, err error) {
	at, isInterface := c.interfaceOf(abstract)
	if isInterface == nil {
		return c.resolveInterface(at, args...)
	}

	at, isStruct := c.structOf(abstract)
	if isStruct == nil {
		return c.resolveStruct(at, args...)
	}

	return nil, errors.New(ERR_RESOLVE_INVALID_ARGUMENTS)
}

func (c *factoryContainer) Inject(target interface{}) error {
	t := reflect.TypeOf(target)
	switch t.Kind() {
	case reflect.Ptr:
	default:
		return errors.New(fmt.Sprintf(ERR_INJECT_INVALID_TARGET_TYPE, t.Kind()))
	}

	s := t.Elem()
	n := s.NumField()
	if n == 0 {
		return nil
	}
	v := reflect.ValueOf(target).Elem()
	for i := 0; i < n; i++ {
		sf := s.Field(i)
		if _, ok := sf.Tag.Lookup("inject"); ok == false {
			continue
		}

		if sf.Type.Kind() != reflect.Interface &&
			sf.Type.Kind() != reflect.Struct && sf.Type.Kind() != reflect.Ptr {
			continue
		}

		f := v.Field(i)
		if f.CanSet() == false {
			continue
		}

		o, err := c.Resolve(sf.Type)
		if err != nil {
			return err
		}

		err = c.Inject(o)
		if err != nil {
			return err
		}
		f.Set(reflect.ValueOf(o))
	}
	return nil
}

func (c *factoryContainer) bindInterface(at reflect.Type, concrete interface{}) error {
	ct := reflect.TypeOf(concrete)
	switch ct.Kind() {
	case reflect.Func:
	case reflect.Ptr:
		if c.instanceOf(at, ct) == false {
			return errors.New(fmt.Sprintf(ERR_BIND_NOT_IMPLEMENT_INTERFACE, ct, at))
		}
	default:
		return errors.New(fmt.Sprintf(ERR_BIND_INVALID_CONCRETE, ct.Kind()))
	}

	c.items.Store(at.String(), reflect.ValueOf(concrete))
	return nil
}

func (c *factoryContainer) bindStruct(at reflect.Type, concrete interface{}) error {
	ct, err := c.structOf(concrete)
	if err != nil {
		return err
	}

	if at.String() != ct.String() {
		return errors.New(fmt.Sprintf(ERR_BIND_INVALID_STRUCT_CONCRETE, at.String(), ct.String()))
	}

	c.items.Store(at.String(), reflect.ValueOf(concrete))
	return nil
}

func (c *factoryContainer) resolveInterface(at reflect.Type, args ...interface{}) (interface{}, error) {
	v, ok := c.items.Load(at.String())
	if ok == false {
		return nil, errors.New(fmt.Sprintf(ERR_RESOLVE_NOT_EXIST_ABSTRACT, at))
	}
	value := v.(reflect.Value)

	switch value.Kind() {
	case reflect.Func:
		return c.resolveFunc(value, args...)
	case reflect.Ptr:
		return value.Interface(), nil
	default:
		return nil, errors.New(fmt.Sprintf(ERR_RESOLVE_INVALID_CONCRETE, value.Kind()))
	}
}

func (c *factoryContainer) resolveStruct(at reflect.Type, args ...interface{}) (interface{}, error) {
	v, ok := c.items.Load(at.String())
	if ok == false {
		return nil, errors.New(fmt.Sprintf(ERR_RESOLVE_NOT_EXIST_ABSTRACT, at))
	}
	value := v.(reflect.Value)

	switch value.Kind() {
	case reflect.Struct, reflect.Ptr:
		return value.Interface(), nil
	default:
		return nil, errors.New(fmt.Sprintf(ERR_RESOLVE_INVALID_CONCRETE, value.Kind()))
	}
}

func (c *factoryContainer) structOf(value interface{}) (reflect.Type, error) {
	if t, ok := value.(reflect.Type); ok == true {
		return c.structOfType(t)
	}

	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New(ERR_BIND_INVALID_STRUCT)
	}

	return t, nil
}

func (c *factoryContainer) structOfType(t reflect.Type) (reflect.Type, error) {
	switch t.Kind() {
	case reflect.Struct:
		return t, nil
	case reflect.Ptr:
		return t.Elem(), nil
	default:
		return nil, errors.New(ERR_BIND_INVALID_STRUCT_TYPE)
	}
}

func (c *factoryContainer) interfaceOf(value interface{}) (reflect.Type, error) {
	if t, ok := value.(reflect.Type); ok == true && t.Kind() == reflect.Interface {
		return t, nil
	}
	t := reflect.TypeOf(value)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		return nil, errors.New(ERR_BIND_INVALID_INTERFACE)
	}

	return t, nil
}

func (c *factoryContainer) instanceOf(abstract reflect.Type, concrete reflect.Type) bool {
	if abstract.Kind() != reflect.Interface {
		return false
	}

	switch concrete.Kind() {
	case reflect.Struct, reflect.Ptr:
		return concrete.Implements(abstract)
	default:
		return false
	}
}

func (c *factoryContainer) resolveFunc(value reflect.Value, args ...interface{}) (interface{}, error) {
	t := value.Type()
	if len(args) != t.NumIn() {
		return nil, errors.New(fmt.Sprintf(ERR_RESOLVE_INSUFFICIENT_ARGUMENTS, t.NumIn(), len(args)))
	}

	in := make([]reflect.Value, t.NumIn())
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	out := value.Call(in)
	if len(out) == 0 {
		return nil, errors.New(ERR_RESOLVE_NON_VALUES_RETURNED)
	}
	return out[0].Interface(), nil
}

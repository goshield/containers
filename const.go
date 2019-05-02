package containers

const (
	ERR_BIND_INVALID_ARGUMENTS         = "binding error! Invalid arguments"
	ERR_RESOLVE_INVALID_ARGUMENTS      = "resolving error! Invalid arguments"
	ERR_INJECT_INVALID_TARGET_TYPE     = "injecting to %v is not supported"
	ERR_BIND_NOT_IMPLEMENT_INTERFACE   = "%v is not an instance of %v"
	ERR_BIND_INVALID_CONCRETE          = "non-supported kind of concrete. Got %v"
	ERR_BIND_INVALID_STRUCT_CONCRETE   = "expects %s. Got %s"
	ERR_RESOLVE_NOT_EXIST_ABSTRACT     = "%v is not bound yet"
	ERR_RESOLVE_INVALID_CONCRETE       = "type %v is not supported"
	ERR_BIND_INVALID_STRUCT            = "called structOf with a value that is not a pointer to a struct. (*MyStruct)(nil)"
	ERR_BIND_INVALID_STRUCT_TYPE       = "called structOfType with a value that is not a pointer to a struct. (*MyStruct)(nil)"
	ERR_BIND_INVALID_INTERFACE         = "called interfaceOf with a value that is not a pointer to an interface. (*MyInterface)(nil)"
	ERR_RESOLVE_INSUFFICIENT_ARGUMENTS = "expects to have %v input arguments. Got %v"
	ERR_RESOLVE_NON_VALUES_RETURNED    = "expects to have at least 1 value returned. Got 0"
)

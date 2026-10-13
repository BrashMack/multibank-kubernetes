package reflector

import "errors"

var (
	ErrNotSupportType                 = errors.New("not support function type")
	ErrFoundedNilInterfaceInStructure = errors.New("founded nil interfaces")
)

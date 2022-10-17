package internal

import (
	cer "github.com/kaatinga/const-errs"
)

const (
	ErrUnsupportedFieldType cer.Error = "unsupported field type"
	ErrModelHasEmptyStruct  cer.Error = "an input struct has no fields"
	ErrNotAStruct           cer.Error = "the configuration must be a struct"
	ErrNotAddressable       cer.Error = "the main struct must be pointed out via pointer"
	ErrNotAddressableField  cer.Error = "the value is not addressable"
)

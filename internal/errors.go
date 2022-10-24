package internal

import (
	"fmt"
	cer "github.com/kaatinga/const-errs"
	"github.com/pkg/errors"
)

const (
	ErrUnsupportedFieldType cer.Error = "unsupported field type"
	ErrModelHasEmptyStruct  cer.Error = "an input struct has no fields"
	ErrNotAStruct           cer.Error = "the configuration must be a struct"
	ErrNotAddressable       cer.Error = "the main struct must be pointed out via pointer"
	ErrNotAddressableField  cer.Error = "the value is not addressable"
	ErrHandle               cer.Error = "unknown error"
	ErrBaseTypeNotMatch     cer.Error = "types not match"
	ErrInvalidOptions       cer.Error = "invalid value options"
)

func errBaserTypeNotMatch(tag, outType, foundedType string) error {
	return errors.Wrap(ErrBaseTypeNotMatch,
		fmt.Sprintf("field tag %s have type: %s; find: %s",
			tag, outType, foundedType))
}

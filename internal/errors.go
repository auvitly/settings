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
	ErrProcessing           cer.Error = "processing error"
	ErrBaseTypeNotMatch     cer.Error = "types not match"
	ErrInvalidOptions       cer.Error = "invalid value options"
	ErrInvalidOptionsType   cer.Error = "invalid options type"
	ErrExceedingExpectValue cer.Error = "exceeding the expected value"
	ErrUnsupportedFieldTag  cer.Error = "unsupported tag"
)

func errBaserTypeNotMatch(tag, outType, foundedType string) error {
	return errors.Wrap(ErrBaseTypeNotMatch,
		fmt.Sprintf("field tag %s have type: %s; find: %s",
			tag, outType, foundedType))
}

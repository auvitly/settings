package tests

type TestENVModelWithTags struct {
	ValueString string `env:"STRING"`
	ValueBool   string `env:"BOOL"`
}

type TestENVModelEmptyTags struct {
	ValueString string
	ValueBool   string
}

type TestComplexModel struct {
	ValueStructure1 *TestENVModelWithTags
	ValueStructure2 TestENVModelWithTags
}

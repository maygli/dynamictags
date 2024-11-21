package dynamictags

// Create processor to process 'env' tag.
// This processor replace structure field with 'env' tag
// by value of environment variable. The environment variable
// name is tag value
// Returns:
//   - Env tag processor.
func NewEnvProcessor() *DynamicTagProcessor {
	processor := DynamicTagProcessor{}
	processor.InitProcessor()
	converter := NewEnvTagConverter()
	processor.AddTagConverter(converter)
	return &processor
}

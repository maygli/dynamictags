package dynamictags

// Create processor to process 'default' tag.
// This processor replace structure field with 'default' tag
// by tag value.
// Returns:
//   - Default tag processor.
func NewDefaultProcessor() *DynamicTagProcessor {
	processor := DynamicTagProcessor{}
	processor.InitProcessor()
	converter := NewDefaultTagConverter()
	processor.AddTagConverter(converter)
	return &processor
}

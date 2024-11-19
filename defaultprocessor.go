package dynamictags

func NewDefaultProcessor() *DynamicTagProcessor {
	processor := DynamicTagProcessor{}
	processor.InitProcessor()
	converter := NewDefaultTagConverter()
	processor.AddTagConverter(converter)
	return &processor
}

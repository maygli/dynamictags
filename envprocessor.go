package dynamictags

func NewEnvProcessor() *DynamicTagProcessor {
	processor := DynamicTagProcessor{}
	processor.InitProcessor()
	converter := NewEnvTagConverter()
	processor.AddTagConverter(converter)
	return &processor
}

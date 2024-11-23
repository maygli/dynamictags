package dynamictags

// Create configuration processor.
// Returns:
//   - Json configuration processor if success.
//   - error if error occured during processor creation
func NewConfigurationProcessor(content any, rootPath string) (*DynamicTagProcessor, error) {
	processor := DynamicTagProcessor{}
	processor.InitProcessor()
	jsonconv, err := NewJsonTagConverter(content, rootPath)
	if err == nil {
		processor.AddTagConverter(jsonconv)
	}
	envconv := NewEnvTagConverter()
	processor.AddTagConverter(envconv)
	defaultconv := NewDefaultTagConverter()
	processor.AddTagConverter(defaultconv)
	return &processor, nil
}

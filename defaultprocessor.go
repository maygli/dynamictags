package dynamictags

const (
	DEFAULT_TAG = "default"
)

func NewDefaultProcessor() *DynamicTagProcessor {
	processor := DynamicTagProcessor{
		Tag: DEFAULT_TAG,
	}
	processor.InitProcessor()
	return &processor
}

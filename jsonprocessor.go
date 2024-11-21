package dynamictags

// Create processor to process 'json' tag.
// This processor replace structure field with 'json' tag
// by value get from json. Example usage:
//     jsonFile, err := os.Open("serverconfiguration.json")
//     if err != nil
//       return err
//     defer jsonFile.close()
//     val, _ := ioutil.ReadAll(jsonFile)
//     var res any
//     err = json.Unmarshal([]byte(val), &res)
//     if err != nil {
//       return err
//     }
//     processor := NewJsonProcessor(res, "$.database")
//     processor.Process(&databaseConfiguration, nil)
// Returns:
//   - Json tag processor if success.
//   - error if error occured during processor creation
func NewJsonProcessor(content any, rootPath string) (*DynamicTagProcessor, error) {
	converter, err := NewJsonTagConverter(content, rootPath)
	if err != nil || converter == nil {
		return nil, err
	}
	processor := DynamicTagProcessor{}
	processor.InitProcessor()
	processor.AddTagConverter(converter)
	return &processor, nil
}

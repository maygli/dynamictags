# dynamictags
Golang package allow to have dynamic generated tags inside structure:
For example for structure:

```go
type TestStruct struct {
  Mode string  `env:`${SERVER_NAME}_MODE`
  Port int     `env:`${SERVER_NAME}_PORT`   
}
```
If environment variable 'SERVER_NAME' has value 'TEST' field 'Mode' will be processed in the follow way:
1) Environment variable name will be calculated as TEST_MODE. 
2) Value of field 'Mode' will get from environment variable 'TEST_MODE' if defined

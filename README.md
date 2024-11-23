# dynamictags
Configuration reader. 
Configuration reader can read configuration from three sources. 1) Default value,
2) Environment variable 3) From json configuration file
Configuration reader allows to have dynaic tags. I.e. tags which value depends on environment variable or dictionary value

For example for structure:

```go
type TestStruct struct {
  Mode string  `env:`${SERVER_NAME}_MODE,defualt:https,json:mode`
  Port int     `env:`${SERVER_NAME}_PORT,default:443,json:server_port`   
}
```
If environment variable 'SERVER_NAME' has value 'TEST' field 'Mode' will be processed in the follow way:
1) Environment variable name will be calculated as TEST_MODE. 
2) Value of field 'Mode' will get from environment variable 'TEST_MODE' if defined

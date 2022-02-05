# Goblin

A go program which will respond with data to mock a server; mainly useful while developing a frontend application, whose backend is not running but API documentation is ready, duh.

### Usage

```
Usage of goblin:
  -a string
        address to listen on (default ":9090")
  -d uint
        delay to induce before each response in milliseconds
  -f string
        path to file with data to mock

Example:
  goblin -f sample.json -a :8080 -d 200
```

### File Format

Goblin requires a file in which the mocked data is to be specified. It should follow the format given below:

```js
{
  "paths": [
    {
      "path": "/todo/{id}/foo",
      "get": {
        "response": {
          // response object - any valid JSON
        }
      },
      "post": {
        "responseFromFile": "relativePath/from/this/file/to/actual/responseFile.json"
      },
      // patch, put, delete are also supported
    },
    // ...
  ]
}
```
The response can be either directly specified in the file or it can be moved to a different json file and point to it using `responseFromFile`.
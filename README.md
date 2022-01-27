# Mocker
A go program which will respond with data to mock a server; mainly useful while developing a frontend application, whose backend is not running but API documentation is ready, duh.

### Usage
```
Usage of mocker:
  -a string
        address to listen on (default ":9090")
  -f string
        path to file with data to mock

Example:
  mocker -f sample.json -a :9090
```

### File Format
Mocker required a file in which the mocked data is to be specified. It should follow the format given below:

```js
{
  "paths": [
    {
      "path": "/todo/{id}/foo",
      "get": {
        "response": {
          // response object
        }
      },
      "post": {
        "response": {
          // response object
        }
      }
    }
  ]
}
```
`response object` can be any valid JSON.
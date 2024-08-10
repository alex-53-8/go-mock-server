# Mock server

## CLI arguments
`-help` - print all arguments
`file` - specify absolute or relative path from a model in a YAML file
`caching-enabled` - enable caching for a responses stored in files in a storage
`cached-item-max-size` - max item's size can be stored in a cache; if a file's size exceeds the number then the response will be read from a storage when requested, otherwise a file is read only once and stored in memory to avoid reading from storage.

## Model
A model of endpoints is defined in a YAML file and has following format

### Top level

- `port` numeric value for port to run a server on.
- `endpoint` list of endpoints, see definition in `Endpoint` section.

```yaml
port: 8081
endpoints:
    - path: /{$}
      ...
```

### Endpoint 

- `path` endpoints path.
- `methods` optional element, if omited or empty then mapped on all method, otherwise specify list of methods to map on
            possible values: get, post, put, delete, patch, options, connect, trace, head
- `headers` a map of headers to be sent in response, key of a map is a header name and value of a map is a header's value
- `response-body` can be just a string a or reference to a file with response
- `status-code` a status code to be returned to a callee

```yaml
#string
response: > 
    {"test": 1}

# file
response: file:sample/response/root/index.html
```

```yaml
endpoints:
  - path: /download
    method: [get]
    response: file:model/responses/download/file.txt
    headers:
      content-type: application/octet-stream
```
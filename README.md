# Mock server
A mock http server for daily work for mocking endpoints with a predefined path, status, headers and response. The project is useful when you want create a mock for one of your services to ease your development. That can be easily done with creating a YAML file with a model. The main benefit of the project - you do not need to write code, just create a configuration in a simple format and run a server in Docker or in terminal.

## CLI arguments

| Item                    | Description                                              |
| ----------------------- | -------------------------------------------------------- |
| `-help`                 | numeric value for port to run a server on                |
| `-file`                 | specify absolute or relative path to a model in a YAML file |
| `-caching-enabled`      | enable caching for a responses stored in files in a storage |
| `-cached-item-max-size` | max item's size can be stored in a cache; if a file's size exceeds the number then the response will be read from a storage when requested, otherwise a file is read only once and stored in memory to avoid reading from storage |

## Run in Docker
Use a `Dockerfile` in a project to run a service in a Docker.

It is expected that a directory with a server's model is mounted on a path `/app/model` in a container as a volume. `/app/model` must contain a server.yml file. 

```shell
# build an image
docker build -t mock-server .
# run a server
docker run -v ./model:/app/model -p 8081:8081 -t mock-server:latest
```

## Run in terminal
Assume you have Go installed on a target machine with minimal version 1.22.4. Then there are two options how to run a server in a terminal. 
1) Just run with Go
```shell
go run . -file=model/server.yml
```

2) Build an application and then run it
```shell
go build -o mock-server && ./mock-server -file=model/server.yml
```

## Model
A model is a definition of endpoints of your server stored in a YAML format. Repository contains a sample configuration in a folder `model`, refer its data as an example. Below you can find description of fields.§

### Root level

| Item                | Mandatory;<br/>Default value | Description                                              |
| ------------------- | ---------------------------- | -------------------------------------------------------- |
| `port`              | Yes                          | Numeric value for port to run a server on.               |
| `endpoints`         | Yes                          | List of endpoints, see definition in `Endpoint` section. |

**Full example**
```yaml
port: 8081
endpoints:
    - path: /{$}
      ...
```

### Endpoint 
| Item                | Mandatory;<br/>Default value | Description            |
| ------------------- | ----------------------------- | --------------------- |
| `path`              | Yes                           | Endpoints path.        |
| `methods`           | No; all methods               | If omited or empty then mapped on all methods. Otherwise specify list of methods to map on. Possible values: get, post, put, delete, patch, options, connect, trace, head. |
| `headers`           | No                            | A map of headers sent in a response, key of a map is a header name and value of a map is a header's value. |
| `status-code`       | No; 200                       | An HTTP status code to be returned to a callee. |
| `response-body`     | No; empty string              | Can be just a string a or reference to a file with response, see description below. |


#### Response-body
The item is not mandatory, can be omitted. If omitted, then no response body is returned

Can be configured as:

**Predefined string in a configuration file**
```yaml
response-body: > 
    {"test": 1}
```

**File stored in a file system**
```yaml
response-body: file:sample/response/root/index.html
```

**Full example**
A response is stored in a file `model/responses/download/file.txt`. We specify in a property `response-body` a string which starts with a prefix `file:` and then a path to a file with a response.

```yaml
endpoints:
  - path: /download
    method: [get]
    response-body: file:model/responses/download/file.txt
    status-code: 201
    headers:
      content-type: application/octet-stream
```

## http_client: A Configurable HTTP Client for Go

http_client is a highly configurable HTTP client for Go, designed for ease of use, extendability, and robustness. It supports various features like automatic retries and request rate limiting, making it ideal for both simple and advanced HTTP requests.

### Features
- **Flexible Request Options**: Easily configure headers, URLs, and authentication.
- **Retry Mechanism**: Automatic retries with configurable policies.
- **Rate Limiting**: Client-side rate limiting to avoid server-side limits.

### Installation
To use Clink in your Go project, install it using `go get`:

```bash
go get -u github.com/apus-run/sea-kit/http_client
```

### Usage
Here is a basic example of how to use http_client:

```go
package main

import (
	"fmt"
	"github.com/apus-run/sea-kit/http_client"
	"net/http"
)

func main() {
	// Create a new client with default options.
	client := http_client.NewClient()

	// Create a new request with default options.
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/anything", nil)

	// Send the request and get the response.
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Hydrate the response body into a safemap.
	var target map[string]any
	err = http_client.ResponseToJson(resp, &target)

	// Print the target safemap.
	fmt.Println(target)
}
```

*HTTP Methods (HEAD, OPTIONS, GET, HEAD, POST, PATCH, DELETE)* are also supported 
```go
package main

import (
	"github.com/apus-run/sea-kit/http_client"
	"encoding/json"
)

func main() {
    client := http_client.NewClient()
    resp, err := client.Get("https://httpbin.org/get")
    // ....
    payload, err := json.Marshal(map[string]string{"username": "yumi"})
    resp, err := client.Post("https://httpbin.org/post", payload)
}
```

### Examples
For more examples, see the [examples](https://github.com/apus-run/sea-kit/http_client/examples) directory.

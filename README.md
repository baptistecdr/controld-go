<h3 align="center">controld-go</h3>
<p align="center">
    controld-go allows you to interact with <a href="https://docs.controld.com/reference/get-started">Control D's API</a> in Go. (obviously)
    <br>
    <a href="https://github.com/baptistecdr/controld-go/issues/new">Report bug</a>
    ·
    <a href="https://github.com/baptistecdr/controld-go/issues/new">Request feature</a>
</p>

<div align="center">

![GitHub Tag](https://img.shields.io/github/v/tag/baptistecdr/controld-go?label=Latest%20version)
[![Go Reference](https://pkg.go.dev/badge/baptistecdr/controld-go.svg)](https://pkg.go.dev/baptistecdr/controld-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/baptistecdr/controld-go?style=flat-square)](https://goreportcard.com/report/github.com/baptistecdr/controld-go)

</div>

## Quick start

```bash
go get github.com/baptistecdr/controld-go
```

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/baptistecdr/controld-go"
)

func main() {
	// Construct a new API object using an API token
	api, err := controld.New(os.Getenv("CONTROLD_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch details of the account
	u, err := api.ListUser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Print user details
	fmt.Println(u)
}
```

## Bugs and feature requests

Have a bug or a feature request? Please first search for existing and closed issues. If your problem or idea is not
addressed yet, [please open a new issue](https://github.com/baptistecdr/controld-go/issues/new).

## Contributing

Contributions are welcome!

## Thanks to

- https://github.com/cloudflare/cloudflare-go

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/korbjt/router)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/korbjt/router/blob/master/LICENSE)
# Router
The router package implements an adapter for the Mux interface from
github.com/korbjt/relay. It uses github.com/julinschmidt/httprouter as the
underlying http multiplexer.

##Usage
```go
package main

import (
    "log"
   	"net/http"
    
	"github.com/korbjt/relay"
	"github.com/korbjt/router"
)

func main() {
    app := &relay.App{
        Mux      : router.New(),
        PreMatch : relay.Chain(),
        PostMatch: relay.Chain(),    
    }

    //add some routes to the app
    log.Panic(http.ListenAndServe(":8080", app.Handler()))
}
```
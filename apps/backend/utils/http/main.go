package http 

import (
	"fmt"
	"net/http"
)

func ProtectEndpoint() {}

func LogRequestReceipt(request *http.Request) {
    var requestPath string = request.URL.Path

    if (len(request.URL.RawQuery) != 0) {
        requestPath += "?" + request.URL.RawQuery
    }

    fmt.Printf(
        "received from %s - \"%s %s %s\"\n", 
        request.RemoteAddr, 
        request.Method, 
        requestPath,
        request.Proto,
    )
}

func LogRequestFailure() {}

func LogRequestCompletion(request *http.Request) {
    var requestPath string = request.URL.Path

    if (len(request.URL.RawQuery) != 0) {
        requestPath += "?" + request.URL.RawQuery
    }

    fmt.Printf(
        "completed for %s - \"%s %s %s %d\"\n", 
        request.RemoteAddr, 
        request.Method, 
        requestPath,
        request.Proto,
        http.StatusOK,
    )
}

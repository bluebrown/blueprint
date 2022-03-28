package values

import "net/http"

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

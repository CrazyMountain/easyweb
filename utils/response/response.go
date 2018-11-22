package response

import "net/http"

var Msg = make(map[int]string)

func init() {
	Msg[http.StatusOK] = "Success."
	Msg[http.StatusBadRequest] = "Illegal arguments."
	Msg[http.StatusInternalServerError] = "Operation failed."
}

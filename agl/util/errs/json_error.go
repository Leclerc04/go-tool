package errs

import "encoding/json"

type ErrorInJSON struct {
	Status       string `json:"status,omitempty"`
	ErrorMsg     string `json:"error_msg,omitempty"`
	ErrorKind    string `json:"error_kind,omitempty"`
	ErrorSubKind string `json:"error_sub_kind,omitempty"`
}

func NewErrorInJSON(err error) ErrorInJSON {
	var ret ErrorInJSON
	if err == nil {
		ret.Status = "ok"
		return ret
	}
	ret.Status = "error"
	e := Wrap(err).(*Error)
	ret.ErrorMsg = e.ToStringNoStack()
	ret.ErrorKind = e.Kind.String()
	ret.ErrorSubKind = e.SubKind
	return ret
}

func (e ErrorInJSON) Error() string {
	bs, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

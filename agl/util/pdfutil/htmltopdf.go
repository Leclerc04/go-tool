package pdfutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/leclecr04/go-tool/agl/util/errs"
)

// HTMLToPDF converts html string to pdf bytes, opts is options for wkhtmltopdf
func HTMLToPDF(htmlStr string, options map[string]string) ([]byte, error) {
	url := "https://e1s3idgl7f.execute-api.cn-north-1.amazonaws.com.cn/prod/html"

	contents := struct {
		HTMLBase64 string            `json:"html_base64"`
		Options    map[string]string `json:"options"`
	}{
		HTMLBase64: base64.StdEncoding.EncodeToString([]byte(htmlStr)),
		Options:    options,
	}

	j, err := json.Marshal(contents)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "hF25hsdFEv4NBfv95jwdV9WZptnOTk919pppVEi6")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := struct {
		PDFBase64    string `json:"pdf_base64"`
		ErrorMessage string `json:"errorMessage"`
	}{}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	if response.ErrorMessage != "" {
		return nil, errs.InvalidArgument.Newf(response.ErrorMessage)
	}
	data, err = base64.StdEncoding.DecodeString(response.PDFBase64)
	if err != nil {
		return nil, err
	}
	return data, nil
}

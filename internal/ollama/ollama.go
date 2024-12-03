package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Ollama struct {
	model       string
	url         string
	stream      bool
	contentType string
}

func (o *Ollama) Generate(prompt string) (string, error) {
	body := fmt.Sprintf(`{
  "model": "%v",
  "prompt": "%v",
  "stream": %v
}`, o.model, prompt, o.stream)
	request, err := http.Post(o.url, o.contentType, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	if request.StatusCode != 200 {
		return "", fmt.Errorf("OllamaRequest returned HTTP status %d", request.StatusCode)
	}
	defer request.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, request.Body)
	if err != nil {
		return "", err
	}
	res := struct {
		Result string `json:"response"`
	}{}
	err = json.Unmarshal(buf.Bytes(), &res)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

func NewOllama(model string, url string, stream bool, contentType string) *Ollama {
	return &Ollama{
		model:       model,
		url:         url,
		stream:      stream,
		contentType: contentType,
	}
}

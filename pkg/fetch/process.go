package fetch

import (
	"path"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/pkg/errors"

	"cactbot_importer/pkg/http"
)

var ErrTransform = errors.New("failed to transform")

func process(resp http.Response) (string, error) {
	ext := path.Ext(resp.URL())

	var s string
	var err error

	switch ext {
	case ".ts":
		s, err = processTS(resp.String())
	default:
		s, err = processDefault(resp.String())
	}

	return s, err
}

func processTS(content string) (string, error) {
	result := api.Transform(content, api.TransformOptions{
		Format:  api.FormatIIFE,
		Charset: api.CharsetUTF8,
		Loader:  api.LoaderTS,
	})

	if len(result.Errors) != 0 {
		s := ""
		for _, message := range result.Errors {
			s += message.Text + "\n"
			for _, note := range message.Notes {
				s += note.Text + "\n"
			}
		}
		return s, ErrTransform
	}

	return string(result.Code), nil
}

func processDefault(content string) (string, error) {
	result := api.Transform(content, api.TransformOptions{
		Format:  api.FormatIIFE,
		Charset: api.CharsetUTF8,
		Loader:  api.LoaderJSX,
	})

	if len(result.Errors) != 0 {
		s := ""
		for _, message := range result.Errors {
			for _, note := range message.Notes {
				s += note.Text
			}
		}
		return s, ErrTransform
	}

	return string(result.Code), nil
}

package repo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/pkg/errors"

	"cactbot_importer/pkg/http"
)

var ErrNestedRepo = errors.New("can't refer a repo in another repo")

func Fetch(remote string) ([]string, error) {
	u, err := url.Parse(remote)
	if err != nil {
		return nil, errors.Wrapf(err, "不是合法的HTTP链接 %s", remote)
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("无法访问 %s", remote)
	}

	v := new(Manifest)
	if err := json.Unmarshal(resp.Body(), v); err != nil {
		return nil, errors.Wrapf(err, "无法解析 %s, 不是合法的json", remote)
	}

	p := path.Dir(u.Path)

	for i, file := range v.Files {
		if path.Ext(file) == ".json" {
			return nil, ErrNestedRepo
		}
		uu := &url.URL{
			Scheme: u.Scheme,
			Host:   u.Host,
			Path:   path.Join(p, file),
		}
		v.Files[i] = uu.String()
	}

	return v.Files, nil
}

type Manifest struct {
	Files []string `json:"files"`
}

package syncer

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Config struct {
	EndPoint string
}

var config Config

func Init(cfg Config) {
	config = cfg
}

func Lock(id string) error {
	r, err := http.NewRequest("GET", config.EndPoint+"/lock?id="+id, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to lock status : %d", resp.StatusCode)
	}

	return nil
}

func Unlock(id string) {
	req, err := http.NewRequest("GET", config.EndPoint+"/unlock?id="+id, nil)
	if err != nil {
		return
	}
	http.DefaultClient.Do(req)
}

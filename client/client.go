package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/clabland/go-homelab-cable/domain"
	"github.com/pkg/errors"
)

type Client struct {
	Server  string
	network string
}

func Connect(host, port string) (*Client, error) {
	c := &Client{
		Server: fmt.Sprintf("%s:%s/api/", host, port),
	}

	resp, err := http.Get(c.Server + "networks")
	if err != nil {
		return c, err
	}

	if resp.StatusCode != http.StatusOK {
		return c, errors.Errorf("non-200: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}

	var networks []domain.Network
	err = json.Unmarshal(body, &networks)
	if err != nil {
		return c, err
	}

	if len(networks) == 0 {
		return c, errors.New("no networks")
	}

	c.network = networks[0].CallSign

	return c, nil
}

func (c *Client) CurrentChannel() (domain.Channel, error) {
	var channel domain.Channel

	resp, err := http.Get(c.Server + "networks/" + c.network + "/live")
	if err != nil {
		return channel, err
	}

	if resp.StatusCode != http.StatusOK {
		return channel, errors.Errorf("non-200: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return channel, err
	}

	err = json.Unmarshal(body, &channel)
	if err != nil {
		return channel, err
	}

	return channel, nil
}

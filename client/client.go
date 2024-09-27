package client

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/loft-sh/vcluster/cmd/vclusterctl/cmd/find"
	"github.com/mheers/vcluster-operator/helpers"
)

type Client struct {
	config *ClientConfig
}

func NewClient(url string) *Client {
	return &Client{
		config: &ClientConfig{
			URL: url,
		},
	}
}

func (c *Client) LoadConfig() error {
	return c.config.Load()
}

type LoginResponse struct {
	Code   int       `json:"code"`
	Expire time.Time `json:"expire"`
	Token  string    `json:"token"`
}

func (c *Client) GetLoginResponse(username, password string) (*LoginResponse, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password must be set")
	}

	// makes a formular request
	data := map[string]string{
		"username": username,
		"password": password,
	}
	urlValues := helpers.MapToUrlValues(data)
	resp, err := http.PostForm(c.config.URL+"/login", urlValues)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	loginResp := &LoginResponse{}

	err = helpers.ReadJSON(resp.Body, loginResp)
	if err != nil {
		return nil, err
	}

	if loginResp.Code != 200 {
		return nil, fmt.Errorf("login failed with code %d", loginResp.Code)
	}
	return loginResp, nil

}

func (c *Client) Login(username, password string) error {
	loginResp, err := c.GetLoginResponse(username, password)
	if err != nil {
		return err
	}

	// get the token
	c.config.Token = loginResp.Token
	c.config.TokenExpire = loginResp.Expire

	err = c.config.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Logout() error {
	return c.config.Delete()
}

func (c *Client) authorizeRequest(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.Token))
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.config.URL+url, nil)
	if err != nil {
		return nil, err
	}
	c.authorizeRequest(req)
	return http.DefaultClient.Do(req)
}

func (c *Client) post(url string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.config.URL+url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	c.authorizeRequest(req)
	return http.DefaultClient.Do(req)
}

func (c *Client) delete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.config.URL+url, nil)
	if err != nil {
		return nil, err
	}
	c.authorizeRequest(req)
	return http.DefaultClient.Do(req)
}

func (c *Client) List() ([]find.VCluster, error) {
	resp, err := c.get("/api/vclusters")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	clusters := []find.VCluster{}
	err = helpers.ReadJSON(resp.Body, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func (c *Client) Get(name string) (*find.VCluster, error) {
	resp, err := c.get("/api/vclusters/" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	cluster := &find.VCluster{}
	err = helpers.ReadJSON(resp.Body, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (c *Client) ClusterToken(name string) (string, error) {
	resp, err := c.get("/api/vclusters/" + name + "/token")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := helpers.ReadBytes(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *Client) Kubeconfig(name string) (string, error) {
	resp, err := c.get("/api/vclusters/" + name + "/kubeconfig")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := helpers.ReadBytes(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *Client) Create(name string) (*find.VCluster, error) {
	resp, err := c.post("/api/vclusters/"+name, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	cluster := &find.VCluster{}
	err = helpers.ReadJSON(resp.Body, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (c *Client) Delete(name string) error {
	resp, err := c.delete("/api/vclusters/" + name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

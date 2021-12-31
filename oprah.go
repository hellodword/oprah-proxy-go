package oprah

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	digest_auth_client "github.com/xinsnake/go-http-digest-auth-client"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func Transport(user, pass, ip, host string, dialer net.Dialer) (*http.Transport, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s:%s@%s",
		user, pass, host))
	if err != nil {
		return nil, err
	}
	return &http.Transport{
		Proxy: http.ProxyURL(u),
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == fmt.Sprintf("%s:443", host) {
				addr = fmt.Sprintf("%s:443", ip)
			}
			return dialer.DialContext(ctx, network, addr)
		},
	}, nil
}

const (
	ClientType = "se0316"
	ClientKey  = "SILrMEPBmJuhomxWkfm3JalqHX2Eheg1YhlEZiMh8II"
)

type Oprah struct {
	c *http.Client
}

type Response struct {
	ReturnCode map[string]string `json:"return_code"`
	Data       struct {
		ClientType     string `json:"client_type"`
		DeviceId       string `json:"device_id"`
		DevicePassword string `json:"device_password"`
		Geos           []struct {
			CountryCode string `json:"country_code"`
			Country     string `json:"country"`
		} `json:"geos"`
		Ips []struct {
			Ip  string `json:"ip"`
			Geo struct {
				CountryCode string `json:"country_code"`
			} `json:"geo"`
			Ports []int `json:"ports"`
		} `json:"ips"`
	} `json:"data"`
}

func (r Response) Status() string {
	if r.ReturnCode == nil {
		return ""
	}
	for _, v := range r.ReturnCode {
		return v
	}
	return ""
}

func New(timeout time.Duration, tr *http.Transport) *Oprah {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar:     jar,
		Timeout: timeout,
	}
	if tr != nil {
		c.Transport = tr
	}

	return &Oprah{
		c: c}
}

func (o *Oprah) post(path, data string) (text string, response Response, err error) {

	req := digest_auth_client.NewRequest(ClientType, ClientKey, http.MethodPost, fmt.Sprintf("https://api.sec-tunnel.com%s", path), data)
	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	header.Set("SE-Client-Version", "beta 60.0.3255.103")
	header.Set("SE-Operating-System", "Windows")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36 OPR/60.0.3255.103 (Edition beta)")
	req.Header = header
	req.HTTPClient = o.c
	r, err := req.Execute()
	if err != nil {
		return
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	text = string(b)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &response)
	return
}

func (o *Oprah) RegisterSubscriber() (string, Response, error) {
	email := fmt.Sprintf("%s@%s.surfeasy.vpn", uuid.New().String(), ClientType)
	m := sha1.New()
	m.Write([]byte(email))
	password := strings.ToUpper(hex.EncodeToString(m.Sum(nil)))

	values := &url.Values{}
	values.Set("email", email)
	values.Set("password", password)
	return o.post("/v4/register_subscriber", values.Encode())
}

func (o *Oprah) RegisterDevice() (string, Response, error) {
	values := &url.Values{}
	values.Set("client_type", ClientType)
	values.Set("device_hash", "4BE7D6F1BD040DE45A371FD831167BC108554111")
	values.Set("device_name", "Opera-Browser-Client")
	return o.post("/v4/register_device", values.Encode())
}

func (o *Oprah) GeoList(device_id string) (string, Response, error) {
	m := sha1.New()
	m.Write([]byte(device_id))
	device_id_hash := strings.ToUpper(hex.EncodeToString(m.Sum(nil)))

	values := &url.Values{}
	values.Set("device_id", device_id_hash)
	return o.post("/v4/geo_list", values.Encode())
}

func (o *Oprah) Discover(device_id, requested_geo string) (string, Response, error) {
	m := sha1.New()
	m.Write([]byte(device_id))
	device_id_hash := strings.ToUpper(hex.EncodeToString(m.Sum(nil)))

	values := &url.Values{}
	values.Set("serial_no", device_id_hash)
	values.Set("requested_geo", fmt.Sprintf(`"%s"`, requested_geo))
	return o.post("/v4/discover", values.Encode())
}

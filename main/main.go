package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

func transport(user, pass, ip string, base net.Dialer) (*http.Transport, error) {
	host := "eu0.sec-tunnel.com:443"
	u, err := url.Parse(fmt.Sprintf("https://%s:%s@%s",
		user, pass, host))
	if err != nil {
		return nil, err
	}
	return &http.Transport{
		Proxy: http.ProxyURL(u),
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == host {
				addr = fmt.Sprintf("%s:443", ip)
			}
			return base.DialContext(ctx, network, addr)
		},
	}, nil
}

func main() {
	t, err := transport("se0316-cbqa9jf2rzaddge8f48q",

		"eyJhbGciOiJFQ0RILUVTK0EyNTZLVyIsImN0eSI6IkpXVCIsImVuYyI6IkEyNTZHQ00iLCJlcGsiOnsia3R5IjoiRUMiLCJjcnYiOiJQLTI1NiIsIngiOiJGV01BenlDeU9KMTI3a1V6alQxcWxsWmgtMmpXM242enBJVE1McUZma2FvIiwieSI6IkYtZHVXVERDaTBoOUFFUU8zRnI4UVd6aWtsR3haX3huR1NISEduQzBDeTAifX0.HQBrCbwdY1uLjEm-bClbARgCmTayErNB70NtD0IOpjJqNUAHDiNfOg.eWjmt0RoMxhP22aV.8FvMIfp4bYyaEcQOvli-5IsnFA2zQ8RyoSoThlsCTnA2UcZ9bOAXOoVU3m-20uA9RGdL0GhNeQt0anxmXRGwv0SLK8ZvfE5Pjkh7FFHQDE-YFf8-m6VYU9huLgSg5Yq6Zvul-IDPbSpUOoUtT9kJ8fEJQM62J2TyasoQUuvsF2zVWGCBHRJOhRaibR2lOYEtI59KwjTkSBU7z0gfEwHpCuK0HdN2f9RrzRKWe1DcFKDzPGPC2loOhRqVaCunhD-7S1ebNiENqjqvQVzAp0P8Okf5Ybz0BoY_OPGhPfvR6rY3lcHmTAF-GfFm_r4psfnkTXvR4JbCo-epElPrBOYWUVctdLuJDmlJIMzDKp_oeZoTVOSyHw8Qd-aeBwL4MpckeaZsfoXYj1xVwRgOppdmDW42h0zwb71vRpKEEv0J2spWxpVkf6BpdZyp5HJr4VcOBg.AaUJkN85Zg0rdFduE-w9ZQ",

		"77.111.245.11",

		net.Dialer{})
	if err != nil {
		panic(err)
	}

	c := &http.Client{
		Transport: t,
	}

	r, err := c.Get("https://httpbin.org/anything")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
}

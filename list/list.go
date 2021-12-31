package main

import (
	"fmt"
	"github.com/hellodword/oprah-proxy-go"
	"net/http"
	"time"
)

func main() {
	var tr *http.Transport

	o := oprah.New(time.Minute, tr)
	_, _, err := o.RegisterSubscriber()
	if err != nil {
		panic(err)
	}

	_, device, err := o.RegisterDevice()
	if err != nil || device.Data.DeviceId == "" {
		panic(err)
	}

	_, geo, err := o.GeoList(device.Data.DeviceId)
	if err != nil || len(geo.Data.Geos) == 0 {
		panic(err)
	}

	for i := range geo.Data.Geos {
		_, ip, err := o.Discover(device.Data.DeviceId, geo.Data.Geos[i].CountryCode)
		if err != nil {
			panic(err)
		}
		for j := range ip.Data.Ips {
			fmt.Println(fmt.Sprintf(`curl -s --proxy "https://%s:%s@eu0.sec-tunnel.com" --resolve eu0.sec-tunnel.com:443:%s http://httpbin.org/anything`,
				device.Data.DeviceId,
				device.Data.DevicePassword,
				ip.Data.Ips[j].Ip,
			))
		}
	}

}

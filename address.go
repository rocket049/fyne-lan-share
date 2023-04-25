package main

import (
	"fmt"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	//	"fyne.io/fyne/v2/canvas"

	qrcode "github.com/skip2/go-qrcode"
)

type QrAddr struct {
	Text  string
	Image fyne.Resource
}

func getAddress() []QrAddr {
	ifs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	res := []QrAddr{}
	ips := []string{}
	for n, if1 := range ifs {
		addrs, err := if1.Addrs()
		if err != nil {
			panic(err)
		}

		for i, addr := range addrs {

			if strings.HasPrefix(addr.String(), "127.") {
				continue
			}

			if strings.Contains(addr.String(), ":") {
				continue
			}

			vs := strings.Split(addr.String(), "/")

			if err != nil {
				panic(err)
			}
			png := fmt.Sprintf("fileserver-%d-%d.png", n, i)
			var addr string
			if strings.Contains(vs[0], ":") {
				addr = fmt.Sprintf("http://[%s]:6868/index", vs[0])
			} else {
				addr = fmt.Sprintf("http://%s:6868/index", vs[0])
			}
			ips = append(ips, vs[0])
			consoleAppend(fmt.Sprintf("Access URL: %s\n", addr))
			data, err := qrcode.Encode(addr, qrcode.Highest, 400)
			//showImg(window, png, addr)
			if err == nil {
				//img := canvas.NewImageFromResource(fyne.NewStaticResource(png, data))
				//img.Resize(fyne.NewSize(400, 400))
				res = append(res, QrAddr{Text: addr, Image: fyne.NewStaticResource(png, data)})
			}

		}
	}
	runDiscover(ips)
	return res
}

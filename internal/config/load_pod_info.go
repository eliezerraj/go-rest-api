package config

import (
	"net"
	"os"
	"log"
	"github.com/go-rest-api/internal/model"
	"strconv"

)

func PodInfo(app *model.ManagerInfo) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error to get the POD IP address !!!", err)
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				app.App.IpAdress = ipnet.IP.String()
			}
		}
	}
	app.App.OSPID = strconv.Itoa(os.Getpid())
}

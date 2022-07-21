package controller

import (
	"log"
	"time"

	//"github.com/go-rest-api/internal/model"
)

type metrics struct {
}

func NewMetricsService() *metrics {
	return &metrics{
	}
}

func (p *metrics) Health() bool {
	log.Printf("Health")
	
	return true
}

func (p *metrics) StressCPU(count int) string {
	log.Printf("StressCPU")
	start := time.Now()

	for n := 0; n <= count; n++ {
		f := make([]int, count+1, count+2)
		if count < 2 {
			f = f[0:2]
		}
		f[0] = 0
		f[1] = 1
		for i := 2; i <= count; i++ {
			f[i] = f[i-1] + f[i-2]
		}
    }

	t := time.Now()
	elapsed := t.Sub(start)

	return "Done in " + elapsed.String()
}
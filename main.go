package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

type Profile struct {
	HeartBeat          uint64
	StatusNotification uint64
	MeterValues        uint64
}

type Borne struct {
	CBI string
	Profile
	Communication
}

func GetTicker(val uint64) *time.Ticker {
	if val != 0 {
		return time.NewTicker(time.Duration(val) * time.Second)
	} else {
		return time.NewTicker(time.Hour * time.Duration(30))
	}
}

func Run(borne Borne) {

	heartBeat := GetTicker(borne.Profile.HeartBeat)
	statusNotification := GetTicker(borne.Profile.StatusNotification)
	meterValues := GetTicker(borne.Profile.MeterValues)

	defer func() {
		heartBeat.Stop()
		statusNotification.Stop()
		meterValues.Stop()
	}()

	for {
		select {
		case now := <-heartBeat.C:
			borne.Communication.Send(fmt.Sprintf("bip...bip... from %s at %v", borne.CBI, now))
		case now := <-statusNotification.C:
			borne.Communication.Send(fmt.Sprintf("status from %s at %v", borne.CBI, now))
		case now := <-meterValues.C:
			borne.Communication.Send(fmt.Sprintf("mesure from %s at %v", borne.CBI, now))
		}
	}

}

type Communication interface {
	Send(msg string)
}

type Console struct{}

func (c Console) Send(msg string) {
	fmt.Println(msg)
}

type ProfileFactory struct{}

func (p ProfileFactory) Create() Profile {

	return Profile{
		HeartBeat:          10,
		StatusNotification: uint64(rnd.UintN(100)),
		MeterValues:        uint64(rnd.UintN(100))}
}

var nbBorne = 10
var rnd *rand.Rand

func init() {

	rnd = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	flag.IntVar(&nbBorne, "n", 10, "Nombre de borne")
	flag.Parse()
}

func main() {

	profileFactory := ProfileFactory{}

	for i := 0; i < nbBorne; i++ {
		cbi := fmt.Sprintf("Borne %d", i)
		go Run(Borne{CBI: cbi,
			Profile:       profileFactory.Create(),
			Communication: Console{}})
	}

	fmt.Println("Appuyez sur Entrée pour arrêter l'application...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Println("L'application est arrêtée.")
}

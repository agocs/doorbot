package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

type directory map[string]string

func (d directory) contains(s string) bool {
	_, ok := d[s]
	return ok
}

func main() {

	dir := readLocalCache("localcache.txt")
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	//NewDirectPinDriver returns a pointer - this wasn't immediately obvious to me
	splate := gpio.NewDirectPinDriver(beagleboneAdaptor, "splate", "P9_12")

	for {
		input := getKBInput()

		if authenticate(input, dir) {
			openDoor(*splate)
		}
	}
}

// getKBInput blocks and waits for someone to scan their rfid card
func getKBInput() string {
	var input string
	fmt.Scanln(&input) // strips newline
	return input
}

func readLocalCache(filename string) directory {
	rawDirBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Unable to read %s. Reason: %v", filename, err)
	}
	rawDirString := string(rawDirBytes)

	dirLinesSlice := strings.Split(rawDirString, "\n")

	dir := directory{}

	for _, l := range dirLinesSlice {
		lineBits := strings.SplitN(l, " ", 2)
		if len(lineBits) != 2 {
			log.Println("Malformatted local cache line: %s", l)
			continue
		}
		dir[lineBits[0]] = lineBits[1]
	}
	return dir
}

func authenticate(code string, dir directory) bool {
	if dir.contains(code) {
		name := dir[code]
		log.Printf("%s scanned in", name)
		return true
	} else {
		log.Println("Refused to admit %s", code)
		return false
	}
}

func openDoor(sp gpio.DirectPinDriver) {
	sp.DigitalWrite(1)
	gobot.After(5*time.Second, func() {
		sp.DigitalWrite(0)
	})
}

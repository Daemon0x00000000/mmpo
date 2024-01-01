package main

import (
	"bufio"
	"go.bug.st/serial"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Printer struct {
	name         string
	port         *serial.Port
	temperatures map[string]float64
	//ReadBuffer   []byte
	mutex sync.Mutex
}

func NewPrinter() (p Printer, err error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open("/dev/ttyACM0", mode)
	if err != nil {
		log.Fatal(err)
		return Printer{}, err
	}

	return Printer{
		name:         "Artillery X2",
		port:         &port,
		temperatures: make(map[string]float64),
	}, nil
}

func (p *Printer) Close() {
	port := *p.port
	err := port.Close()
	if err != nil {
		log.Fatalf("Cannot close port: %v\n", err)
	}
}

/*
No locking here as read method is supposed to be only called from Send method
*/
func (p *Printer) read() (response string, ok bool) {

	port := *p.port
	response = ""
	var nbBytes int
	var err error
	readBuffer := make([]byte, 64)
	for {
		nbBytes, err = port.Read(readBuffer)
		if err != nil {
			log.Fatal(err)
		}
		if nbBytes == 0 {
			break
		}

		response += string(readBuffer[:nbBytes])
		// EOF on LF
		if readBuffer[nbBytes-1] == 10 {
			break
		}
	}
	if strings.HasPrefix(response, "ok") {
		return strings.TrimSpace(
			strings.TrimPrefix(response, "ok")), true
	}
	return response, false

}

func (p *Printer) Send(cmd string) (data string, ok bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	port := *p.port
	_, err := port.Write([]byte(cmd + "\n\r"))
	if err != nil {
		log.Fatalln(err)
	}
	return p.read()

}

func (p *Printer) AutoTemperaturesJob() {

	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				data, ok := p.Send("M105")
				if ok {
					tPattern := regexp.MustCompile(`T:([0-9.]+)`)
					bPattern := regexp.MustCompile(`B:([0-9.]+)`)

					tValueStr := tPattern.FindStringSubmatch(data)[1]
					bValueStr := bPattern.FindStringSubmatch(data)[1]

					parseFloat := func(val string) (f float64) {
						var err error
						f, err = strconv.ParseFloat(val, 64)
						if err != nil {
							return 0
						}
						return f
					}

					p.temperatures["Extruder"] = parseFloat(tValueStr)
					p.temperatures["Bed"] = parseFloat(bValueStr)

					log.Println(p.temperatures)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (p *Printer) Print(name string) (ok bool) {
	file, err := os.Open("./prints/" + name + ".gcode")
	defer file.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	reader := bufio.NewReader(file)
	fi, err := file.Stat()
	if err != nil {
		// Could not obtain stat, handle error
	}
	log.Println(fi.Size())
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if line[0] == 59 {

		}
		//data, ok := p.Send(string(line))

	}
	return true
}

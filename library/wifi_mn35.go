package main

import (
	"machine"
	"strconv"
	"time"
)

type PinInfo struct {
	Num uint8
	Pin machine.Pin
}

type WIFI_MN35 struct {
	Exit_Pin      *PinInfo
	UART_Instance machine.UART
	RX_Pin        machine.Pin
	TX_Pin        machine.Pin
	BaudRate      uint16
}

/*
GET /hello.htm HTTP/1.1
User-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)
Host: www.tutorialspoint.com
Accept-Language: en-us
Accept-Encoding: gzip, deflate
Connection: Keep-Alive


POST /cgi-bin/birthday.pl HTTP/1.0
User-Agent; Mozilla/4.05 (WinNT; 1)
Host: www.ora.com
Content-type: application/x-www-form-urlencoded
Content-Length: 20

nionth=august&date=24
*/

type HttpPackage struct {
	Method        string
	Path          string
	HttpVersion   string
	Host          string
	ContentLength uint16
	Headers       map[string]string
	Body          string
}

const (
	ACCESS_POINT_NAME = "Test"
	ACCESS_POINT_PASS = "12345"
)

var definedBaudRate = []uint32{1200, 4800, 9600, 16384, 19200, 32000, 48000, 72000, 96000, 100000, 115200}

var commandMapping = map[string]string{
	"ping":        "at",
	"reboot":      "at+Rb=1\r",
	"getVersion":  "at+ver=?\r",
	"getMac":      "at+mac=?\r",
	"transparent": "at+TS",
	"save":        "at+WC=1\r",
	//for AP connection
	"wifiMode":          "at+WA", // STA - client, AP - server | 1 - AP, 2..etc - STA
	"pointByMac":        "at+Sbssid",
	"pointByName":       "at+Sssid",
	"pointByNameLength": "at+Sssidl",
	"encryption":        "at+Sam", // 9 - Wpa/Wpa2_aes
	"pointPassword":     "at+Spw",
	"passLength":        "at+Spwl",
	"type":              "at+UType",
	"dhcp":              "at+dhcp",
	"mask":              "at+mask",
	"dns":               "at+dns",
	"gateway":           "at+gw",
	"ip":                "at+ip",

	//Remote
	"remoteServer": "at+UIp",
	"remoteType":   "at+UType",
	"remotePort":   "at+URPort",

	//Buffer
	"messageLength":   "at+UPL",
	"messageTimeout":  "at+UPT",
	"messageInterval": "at+UPT2",

	//Socket
	"sOpen":  "at+SO",
	"sClose": "at+SC",
	"sCheck": "at+SL",
	"sRead":  "at+SR",
	"sSend":  "at+SW",

	//Utils
	"getIP": "at+DR",
}

const (
	fallbackNum = 30
	minBaud     = 1200
	maxBaud     = 115200
	specialChar = "at+"
)

func MN35_New(gpio4 machine.Pin, uart machine.UART, rx, tx machine.Pin) *WIFI_MN35 {
	println("start initializing wifi")

	instance := WIFI_MN35{
		Exit_Pin: &PinInfo{
			Num: uint8(gpio4),
			Pin: gpio4,
		},
		UART_Instance: uart,
		RX_Pin:        rx,
		TX_Pin:        tx,
		BaudRate:      9600,
	}

	instance.Exit_Pin.Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	instance.setCommandMode(true)
	time.Sleep(time.Microsecond * 50)
	instance.findBaudRate()
	return &instance
}

func (wifi *WIFI_MN35) findBaudRate() {
	for _, rate := range definedBaudRate {
		wifi.UART_Instance.Configure(machine.UARTConfig{ // init UART
			BaudRate: rate,
			TX:       wifi.TX_Pin,
			RX:       wifi.RX_Pin,
		})
		data := wifi.listenImmediate([]byte(commandMapping["getVersion"]), true)
		if string(data) == string(commandMapping["ping"]) {
			println("Module found baudrate:", rate)
			wifi.BaudRate = uint16(rate)
			break
		}
		println(rate)
	}
}

func (wifi *WIFI_MN35) Try(command string, withLogs bool) {
	var result []byte
	wifi.UART_Instance.Write([]byte(command))
	time.Sleep(time.Millisecond * 10)
	if withLogs {
		println("Current Command:", string(command))
	}
	if wifi.UART_Instance.Buffered() > 0 {
		result = make([]byte, 64)
		_, _ = wifi.UART_Instance.Read(result)
		if withLogs {
			println("listenImmediate::", string(result))
		}
	}
	defer func() {
		result = nil
	}()
}

func (wifi *WIFI_MN35) listenCommand(command []byte, isWait bool) ([]byte, bool) {
	var result []byte
	isGotBytes := false
	currentFallback := 0
	wifi.UART_Instance.Write(command)
	println(string(command))
	for {
		time.Sleep(time.Millisecond * 50)
		if currentFallback >= fallbackNum || !isWait {
			break
		}
		if wifi.UART_Instance.Buffered() > 0 {
			result = make([]byte, 30)
			_, _ = wifi.UART_Instance.Read(result)
			println("listenCommand-for::", string(result))
			isGotBytes = true
			break
		}
		currentFallback++
	}
	return result, isGotBytes
}

func (wifi *WIFI_MN35) listenImmediate(command []byte, withLogs bool) []byte {
	var result []byte
	wifi.UART_Instance.Write(command)
	time.Sleep(time.Millisecond * 10)
	if withLogs {
		println("Current Command:", string(command))
	}
	if wifi.UART_Instance.Buffered() > 0 {
		result = make([]byte, 30)
		_, _ = wifi.UART_Instance.Read(result)
		if withLogs {
			println("listenImmediate::", string(result))
		}
	}
	return result
}

func (wifi *WIFI_MN35) setCommandMode(mode bool) {
	if mode {
		wifi.Exit_Pin.Pin.Set(false)
		time.Sleep(time.Second * 2)
		wifi.Exit_Pin.Pin.Set(true)
	} else {
		wifi.listenCommand(handleMapCommands(commandMapping["reboot"], "1"), false)
	}
}

func (wifi *WIFI_MN35) ConnectToPoint(name, password string) {
	wifi.listenImmediate(handleMapCommands(commandMapping["wifiMode"], "2"), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["dhcp"], "1"), true)
	// wifi.listenImmediate(handleMapCommands(commandMapping["ip"], "192,168,1,250"), true)
	// wifi.listenImmediate(handleMapCommands(commandMapping["dns"], "0,0,0,0"), true)
	// wifi.listenImmediate(handleMapCommands(commandMapping["mask"], "255,255,255,0"), true)
	// wifi.listenImmediate(handleMapCommands(commandMapping["gateway"], "192,168,1,1"), true)

	wifi.listenImmediate(handleMapCommands(commandMapping["pointByName"], ACCESS_POINT_NAME), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["pointByNameLength"], strconv.Itoa(len(ACCESS_POINT_NAME))), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["encryption"], "9"), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["pointPassword"], ACCESS_POINT_PASS), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["passLength"], strconv.Itoa(len(ACCESS_POINT_PASS))), true)

	wifi.listenImmediate(handleMapCommands(commandMapping["messageLength"], "128"), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["messageInterval"], "1000"), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["messageTimeout"], "1000"), true)

	wifi.listenImmediate(handleMapCommands(commandMapping["remoteType"], "2"), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["remoteServer"], "192.168.1.90"), true)
	wifi.listenImmediate(handleMapCommands(commandMapping["remotePort"], "8080"), true)

	// wifi.listenImmediate([]byte(commandMapping["save"]), true)
	wifi.listenImmediate([]byte(commandMapping["reboot"]), true)
}

func handleMapCommands(key, value string) []byte {
	return []byte(key + "=" + value + "\r")
}

func (wifi *WIFI_MN35) GetIPAddress(url string) []byte {
	wifi.listenImmediate(handleMapCommands(commandMapping["getIP"], url), false)
	time.Sleep(time.Second)
	return wifi.listenImmediate(handleMapCommands(commandMapping["getIP"], url), false)
}

//http

func (wifi *WIFI_MN35) generateHTTPackage(pkg *HttpPackage) string {
	httpHeaderLine := pkg.Method + " " + pkg.Path + " " + pkg.HttpVersion + "\n"
	httpHostLine := "Host: " + pkg.Host + "\n"
	httpHeaders := ""
	for key, value := range pkg.Headers {
		httpHeaders += key + ": " + value + "\n"
	}
	httpBody := pkg.Body + "\n"
	return httpHeaderLine + httpHostLine + httpHeaders + "\n" + httpBody
}

func (wifi *WIFI_MN35) GetMac() {
	response := wifi.listenImmediate([]byte(commandMapping["getMac"]), false)
	for _, v := range response {
		hex := strconv.FormatInt(int64(v), 16)
		print(hex, ":")
	}
}

/*
WA Wifi mode,ap/sta
WM Wifista method:manual or smartconfig
Sbssid set target ap bssid
Sssid set target ap ssid
Sssidl set target ap ssid length
Sam set target ap encryption method
Spw set target ap key
Spwl set length of target ap key
WC calculation PMK
dhcp set dhcp or static
ip static ip
mask Static mask
dns Static DNS
gw Static gateway
Ub Set uart bandrate
Ud Set uart datalength
Up Serial parity bit

Us Serial stop bit length
UType Set TCP or UDP
UIp Set remote ip address
URPort Set remote port
ULPort Set local port
UPL Set or query data length of automatic framing
UPT Set or query period of automatic framing
UPT2 Set or query Interval period of automatic framing
DP Prefix data for UDP/988 port executes the at command
DE UDP/988 port executes the at command enable or disable
Rb Reboot the module
ver version
Df Back to default setting
SO Socket open
SC Socket close
SL Socket check
SW Socket send
SR Socket read
DR Domain name resolution
GW GPIO write
GR GPIO read
TS Transparent ransmission change
mac Get mac address
Assid Softap SSID
Assidl Softap SSID length
Achan Softap wifi channel
Aam Softap encryption method
Apw Softap key
Apwl Softap key length
Aip Softap the module’s ip address
*/

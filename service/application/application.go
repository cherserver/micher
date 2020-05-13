package application

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/ipv4"

	"github.com/cherserver/micher/service/interfaces"
)

func devicesBroadcastAddressIp() net.IP {
	return net.IPv4(224, 0, 0, 50)
}

//func devicesBroadcastAddressIp() net.IP {
//	return net.IPv4(224, 0, 0, 50)
//}

// const udpListen = "0.0.0.0:9898"
const udpListen = "224.0.0.50:9898"

func devicesBroadcastAddressUdp() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 50),
		Port: 9898,
	}
}

const (
	maxDataGramSize = 8192
)

type application struct {
	env        interfaces.Environment
	httpServer http.Server
}

func New(environment interfaces.Environment) *application {
	return &application{
		env: environment,
		httpServer: http.Server{
			Addr:              "",
			Handler:           nil,
			TLSConfig:         nil,
			ReadTimeout:       0,
			ReadHeaderTimeout: 0,
			WriteTimeout:      0,
			IdleTimeout:       0,
			MaxHeaderBytes:    0,
			TLSNextProto:      nil,
			ConnState:         nil,
			ErrorLog:          nil,
			BaseContext:       nil,
			ConnContext:       nil,
		},
	}
}

func (a *application) Init() error {
	// Multicast()
	// ListenBroadcast()

	MulticastWithJoin()

	// client, err := mi.NewRouter("")
	// if err != nil {
	//	return err
	// }

	// err = client.Connect()
	// if err != nil {
	//  return err
	// }
	return nil
}

type Payload struct {
	Cmd     string
	Mode    string
	Sid     string
	ShortId uint32
	Token   string
	Data    string
}

func Multicast() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("net.Interfaces: %v", err)
	}

	devicesConn, err := net.ListenPacket("udp4", "0.0.0.0:9898")
	if err != nil {
		log.Fatalf("net.ListenPacket: %v", err)
	}

	defer func() { _ = devicesConn.Close() }()

	packetConn := ipv4.NewPacketConn(devicesConn)

	for _, netInterface := range netInterfaces {
		currentInterface := &netInterface
		interfaceName := currentInterface.Name

		if currentInterface.Flags&(net.FlagMulticast) == 0 {
			log.Printf("Interface '%v' have no multicast option, skip it", interfaceName)
			continue
		}

		log.Printf("Interface '%v' has multicast option, use it", interfaceName)

		multicastAddresses, err := currentInterface.MulticastAddrs()
		if err != nil {
			log.Printf("Failed to get interface '%v' multicast addresses, skip it: %v", interfaceName, err)
		}

		for _, address := range multicastAddresses {
			log.Printf("Interface '%v' multicast address '%v'", interfaceName, address)
		}

		miHubLinkLocal := net.UDPAddr{
			IP: devicesBroadcastAddressIp(),
		}

		if err := packetConn.JoinGroup(currentInterface, &miHubLinkLocal); err != nil {
			log.Printf("Interface '%v' JoinGroup: %v", interfaceName, err)
			continue
		}

		log.Printf("Interface '%v': joined %v group", interfaceName, miHubLinkLocal)

		/*go func() {
			defer func() { _ = packetConn.LeaveGroup(currentInterface, &miHubLinkLocal) }()

			b := make([]byte, 4096)
			for {
				cnt, _, peer, err := packetConn.ReadFrom(b)
				if err != nil {
					log.Printf("Interface '%v' ReadFrom error: %v", interfaceName, err)
					return
				}

				log.Printf("Interface '%v' Read '%v' bytes from '%s' ( network '%s')", interfaceName, cnt, peer.String(), peer.Network())

				if cnt <= 0 {
					log.Printf("Interface '%v': zero bytes read, continue", interfaceName)
					continue
				}

				rawData := b[:cnt]

				log.Printf("Raw payload:\n%s", string(rawData))

				var payload Payload
				err = json.Unmarshal(rawData, &payload)
				if err != nil {
					log.Fatalf("Failed to unmarshal received payload: %v", err)
				}

				log.Printf("Incoming data:\n%v", payload)
			}
		}() */
	}
}

func ListenBroadcast() {
	addr, err := net.ResolveUDPAddr("udp4", devicesBroadcastAddressUdp().String())
	if err != nil {
		log.Fatalf("ResolveUDPAddr: %v", err)
	}

	log.Print("Address resolved")

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("ListenMulticastUDP: %v", err)
	}

	log.Print("Listen set")

	_ = conn.SetReadBuffer(maxDataGramSize)

	log.Print("Buffer set")

	// Loop forever reading from the socket
	for {
		log.Print("In for")
		buffer := make([]byte, maxDataGramSize)
		numBytes, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		log.Printf("%v, %v, %v", src, numBytes, buffer)
	}
}

func MulticastWithJoin() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("net.Interfaces: %v", err)
	}

	devicesConn, err := net.ListenPacket("udp4", udpListen)
	if err != nil {
		log.Fatalf("net.ListenPacket: %v", err)
	}

	// defer func() { _ = devicesConn.Close() }()

	packetConn := ipv4.NewPacketConn(devicesConn)

	//created := false

	for _, netInterface := range netInterfaces {
		currentInterface := &netInterface
		interfaceName := currentInterface.Name

		if currentInterface.Flags&(net.FlagMulticast) == 0 {
			log.Printf("Interface '%v' have no multicast option, skip it", interfaceName)
			continue
		}

		log.Printf("Interface '%v' has multicast option, use it", interfaceName)

		multicastAddresses, err := currentInterface.MulticastAddrs()
		if err != nil {
			log.Printf("Failed to get interface '%v' multicast addresses, skip it: %v", interfaceName, err)
		}

		for _, address := range multicastAddresses {
			log.Printf("Interface '%v' multicast address '%v'", interfaceName, address)
		}

		miHubLinkLocal := &net.UDPAddr{
			IP: devicesBroadcastAddressIp(),
		}

		// miHubLinkLocal := devicesBroadcastAddressUdp()

		//if interfaceName != "Ethernet" {
		//	log.Printf("Interface '%v' is not Ethernet, skip", interfaceName)
		//	continue
		//}

		//if created {
		//	continue
		//}

		//created = true

		if err := packetConn.JoinGroup(currentInterface, miHubLinkLocal); err != nil {
			log.Printf("Interface '%v' JoinGroup: %v", interfaceName, err)
			continue
		}

		log.Printf("Interface '%v': joined %v group", interfaceName, miHubLinkLocal)

		go func() {
			defer func() { _ = packetConn.LeaveGroup(currentInterface, miHubLinkLocal) }()

			b := make([]byte, 4096)
			for {
				cnt, _, peer, err := packetConn.ReadFrom(b)
				if err != nil {
					log.Printf("Interface '%v' ReadFrom error: %v", interfaceName, err)
					return
				}

				log.Printf("Interface '%v' Read '%v' bytes from '%s' ( '%s' network)", interfaceName, cnt, peer.String(), peer.Network())

				if cnt <= 0 {
					log.Printf("Interface '%v': zero bytes read, continue", interfaceName)
					continue
				}

				rawData := b[:cnt]

				log.Printf("Raw payload:\n%s", string(rawData))

				var payload Payload
				err = json.Unmarshal(rawData, &payload)
				if err != nil {
					log.Fatalf("Failed to unmarshal received payload: %v", err)
				}

				log.Printf("Incoming data:\n%v", payload)
			}
		}()
	}
}

/*

func Multicast() {
	ifaces, err := net.Interfaces()
	log.Printf("ifaces: %v", ifaces)

	c, err := net.ListenPacket("udp4", "0.0.0.0:9898")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	p := ipv4.NewPacketConn(c)

	en0, err := net.InterfaceByName("Ethernet")
	if err != nil {
		log.Fatalf("InterfaceByName: %v", err)
	}
	miHubLinkLocal := net.UDPAddr{IP: net.IPv4(224, 0, 0, 50)}
	if err := p.JoinGroup(en0, &miHubLinkLocal); err != nil {
		log.Fatalf("JoinGroup: %v", err)
	}
	defer func() { _ = p.LeaveGroup(en0, &miHubLinkLocal) }()

	b := make([]byte, 4096)
	for {
		cnt, _, peer, err := p.ReadFrom(b)
		if err != nil {
			log.Fatalf("error ReadFrom: %v", err)
		}

		log.Printf("Read '%v' bytes from '%s' ( network '%s')", cnt, peer.String(), peer.Network())

		if cnt <= 0 {
			log.Print("Zero bytes read, continue")
			continue
		}

		rawData := b[:cnt]

		log.Printf("Raw payload:\n%s", string(rawData))

		var payload Payload
		err = json.Unmarshal(rawData, &payload)
		if err != nil {
			log.Fatalf("Failed to unmarshal received payload: %v", err)
		}

		log.Printf("Incoming data:\n%s", payload)
	}
}


*/

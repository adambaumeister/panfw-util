package pcaptest

import (
	"github.com/adambaumeister/panfw-util/panos/errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func ReadPcap(fn string) []Flow {
	var flows []Flow

	if handle, err := pcap.OpenOffline(fn); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			f := handlePacket(packet) // Do something with a packet here.
			flows = append(flows, f)
		}
	}
	return Conversations(flows)
}

func handlePacket(p gopacket.Packet) Flow {
	net := p.NetworkLayer()
	src, dst := net.NetworkFlow().Endpoints()
	proto := p.TransportLayer()
	if proto == nil {
		errors.LogDebug("IP packet with no transport header, cannot convert to flow.")
	}
	flow := Flow{}
	tcpLayer := p.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		flow.Protocol = 6
	} else {
		flow.Protocol = 17
	}
	srcport, dstport := proto.TransportFlow().Endpoints()
	flow.SourceIP = src.String()
	flow.DestinationIP = dst.String()
	flow.Port = dstport.String()
	flow.sourcePort = srcport.String()

	flow.nethash = net.NetworkFlow().FastHash()
	flow.transportHash = proto.TransportFlow().FastHash()
	return flow
}

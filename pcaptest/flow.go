package pcaptest

import "fmt"

type Flow struct {
	SourceIP      string `xml:"source"`
	DestinationIP string `xml:"destination"`
	Protocol      int    `xml:"protocol"`
	sourcePort    string
	Port          string `xml:"destination-port"`

	From *string `xml:"from,omitempty"`
	To   *string `xml:"to,omitempty"`

	nethash       uint64
	transportHash uint64
}

func (f *Flow) Print() {
	fmt.Printf("(%v) %v:%v %v:%v\n", f.Protocol, f.SourceIP, f.sourcePort, f.DestinationIP, f.Port)
}

func Conversations(flows []Flow) []Flow {
	c := []Flow{}
	var nethash = make(map[uint64]map[uint64]Flow)
	for _, f := range flows {
		// Build based on first packet seen in a conversation
		if nethash[f.nethash] == nil {
			nethash[f.nethash] = map[uint64]Flow{}
			nethash[f.nethash][f.transportHash] = f
		}
	}

	for nh, _ := range nethash {
		for _, f := range nethash[nh] {
			c = append(c, f)
		}
	}
	return c
}

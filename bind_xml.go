package main

import (
	"encoding/xml"
	"fmt"
	"time"
)

type bindXmlStats struct {
	Memory struct {
		Contexts struct {
			Context []struct {
				Blocksize   string `xml:"blocksize"`
				Hiwater     int    `xml:"hiwater"`
				ID          string `xml:"id"`
				Inuse       int    `xml:"inuse"`
				Lowater     int    `xml:"lowater"`
				Malloced    int    `xml:"malloced"`
				Maxinuse    int    `xml:"maxinuse"`
				Maxmalloced int    `xml:"maxmalloced"`
				Name        string `xml:"name"`
				Pools       int    `xml:"pools"`
				References  int    `xml:"references"`
				Total       int    `xml:"total"`
			} `xml:"context"`
		} `xml:"contexts"`
		Summary struct {
			BlockSize   int `xml:"BlockSize"`
			ContextSize int `xml:"ContextSize"`
			InUse       int `xml:"InUse"`
			Lost        int `xml:"Lost"`
			Malloced    int `xml:"Malloced"`
			TotalUse    int `xml:"TotalUse"`
		} `xml:"summary"`
	} `xml:"memory"`
	Server struct {
		BootTime   time.Time `xml:"boot-time"`
		ConfigTime time.Time `xml:"config-time"`
		Counters   []struct {
			Type    string `xml:"type,attr"`
			Counter []struct {
				Name     string `xml:"name,attr"`
				CharData string `xml:",chardata"`
			} `xml:"counter"`
		} `xml:"counters"`
		CurrentTime time.Time `xml:"current-time"`
		Version     string    `xml:"version"`
	} `xml:"server"`
	Socketmgr struct {
		Sockets struct {
			Socket []struct {
				ID           string  `xml:"id"`
				LocalAddress *string `xml:"local-address"`
				Name         *string `xml:"name"`
				PeerAddress  string  `xml:"peer-address"`
				References   int     `xml:"references"`
				States       struct {
					State []string `xml:"state"`
				} `xml:"states"`
				Type string `xml:"type"`
			} `xml:"socket"`
		} `xml:"sockets"`
	} `xml:"socketmgr"`
	Taskmgr struct {
		Tasks struct {
			Task []struct {
				Events     int     `xml:"events"`
				ID         string  `xml:"id"`
				Name       *string `xml:"name"`
				Quantum    int     `xml:"quantum"`
				References int     `xml:"references"`
				State      string  `xml:"state"`
			} `xml:"task"`
		} `xml:"tasks"`
		ThreadModel struct {
			DefaultQuantum int    `xml:"default-quantum"`
			Type           string `xml:"type"`
		} `xml:"thread-model"`
	} `xml:"taskmgr"`
	Traffic struct {
		Ipv4 struct {
			Tcp struct {
				Counters []struct {
					Type    string `xml:"type,attr"`
					Counter []struct {
						Name     string `xml:"name,attr"`
						CharData string `xml:",chardata"`
					} `xml:"counter"`
				} `xml:"counters"`
			} `xml:"tcp"`
			Udp struct {
				Counters []struct {
					Type    string `xml:"type,attr"`
					Counter []struct {
						Name     string `xml:"name,attr"`
						CharData string `xml:",chardata"`
					} `xml:"counter"`
				} `xml:"counters"`
			} `xml:"udp"`
		} `xml:"ipv4"`
		Ipv6 struct {
			Tcp struct {
				Counters []struct {
					Type    string `xml:"type,attr"`
					Counter []struct {
						Name     string `xml:"name,attr"`
						CharData string `xml:",chardata"`
					} `xml:"counter"`
				} `xml:"counters"`
			} `xml:"tcp"`
			Udp struct {
				Counters []struct {
					Type    string `xml:"type,attr"`
					Counter []struct {
						Name     string `xml:"name,attr"`
						CharData string `xml:",chardata"`
					} `xml:"counter"`
				} `xml:"counters"`
			} `xml:"udp"`
		} `xml:"ipv6"`
	} `xml:"traffic"`
	Views struct {
		View []struct {
			Name  string `xml:"name,attr"`
			Cache struct {
				Name  string `xml:"name,attr"`
				Rrset []struct {
					Counter int    `xml:"counter"`
					Name    string `xml:"name"`
				} `xml:"rrset"`
			} `xml:"cache"`
			Counters []struct {
				Type    string `xml:"type,attr"`
				Counter []struct {
					Name     string `xml:"name,attr"`
					CharData string `xml:",chardata"`
				} `xml:"counter"`
			} `xml:"counters"`
			Zones struct {
				Zone []struct {
					Name       string `xml:"name,attr"`
					Rdataclass string `xml:"rdataclass,attr"`
					Counters   []struct {
						Type    string `xml:"type,attr"`
						Counter []struct {
							Name     string `xml:"name,attr"`
							CharData string `xml:",chardata"`
						} `xml:"counter"`
					} `xml:"counters"`
					Expires *time.Time `xml:"expires"`
					Loaded  time.Time  `xml:"loaded"`
					Refresh *time.Time `xml:"refresh"`
					Serial  int        `xml:"serial"`
					Type    string     `xml:"type"`
				} `xml:"zone"`
			} `xml:"zones"`
		} `xml:"view"`
	} `xml:"views"`
}

func ReadXmlStats(statsData []byte) error {
	fmt.Printf("Read %d bytes of XML\n", len(statsData))

	var xmlStats bindXmlStats

	// Parse the XML statistics
	err := xml.Unmarshal(statsData, &xmlStats)
	if err != nil {
		fmt.Printf("Error parsing XML: %s\n", err)
		return err
	}

	return nil
}

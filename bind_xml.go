package main

import (
	"encoding/xml"
	"fmt"
)

type bindXmlStats struct {
	Stats []struct {
		Name  string `xml:"name,attr"`
		Value int    `xml:"value,attr"`
	} `xml:"statistics"`
}

func readXmlStats(statsData []byte) error {
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

package model

import "encoding/xml"

// Datafile represents the datomatic XML format stripped down to just the parts we need
type Datafile struct {
	XMLName xml.Name `xml:"datafile"`
	Games   []Game   `xml:"game"`
}

type Game struct {
	XMLName xml.Name `xml:"game"`
	Name    string   `xml:"name,attr"`
	ROM     Rom      `xml:"rom"`
}

type Rom struct {
	XMLName xml.Name `xml:"rom"`
	CRC32   string   `xml:"crc,attr"`
}

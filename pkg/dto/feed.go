package dto

import (
	"encoding/xml"
	"fmt"

	"github.com/bdrbt/stllc/internal/domain"
)

// Comment: i have no specs about that xml structure docs
// so i assume that all Fileds witl "List" postfix must be handled as arrays.

type SdnEntry struct {
	Text        string `xml:",chardata"`
	Uid         int64  `xml:"uid"`
	FirstName   string `xml:"firstName"`
	LastName    string `xml:"lastName"`
	Title       string `xml:"title"`
	SdnType     string `xml:"sdnType"`
	ProgramList struct {
		Text    string `xml:",chardata"`
		Program string `xml:"program"`
	} `xml:"programList"`

	AkaList struct {
		Text string `xml:",chardata"`
		Aka  []struct {
			Text      string `xml:",chardata"`
			Uid       string `xml:"uid"`
			Type      string `xml:"type"`
			Category  string `xml:"category"`
			LastName  string `xml:"lastName"`
			FirstName string `xml:"firstName"`
		} `xml:"aka"`
	} `xml:"akaList"`

	DateOfBirthList struct {
		Text            string `xml:",chardata"`
		DateOfBirthItem []struct {
			Text        string `xml:",chardata"`
			Uid         string `xml:"uid"`
			DateOfBirth string `xml:"dateOfBirth"`
			MainEntry   string `xml:"mainEntry"`
		} `xml:"dateOfBirthItem"`
	} `xml:"dateOfBirthList"`

	PlaceOfBirthList struct {
		Text             string `xml:",chardata"`
		PlaceOfBirthItem struct {
			Text         string `xml:",chardata"`
			Uid          string `xml:"uid"`
			PlaceOfBirth string `xml:"placeOfBirth"`
			MainEntry    string `xml:"mainEntry"`
		} `xml:"placeOfBirthItem"`
	} `xml:"placeOfBirthList"`

	IdList struct {
		Text string `xml:",chardata"`
		ID   []struct {
			Text      string `xml:",chardata"`
			Uid       string `xml:"uid"`
			IdType    string `xml:"idType"`
			IdNumber  string `xml:"idNumber"`
			IdCountry string `xml:"idCountry"`
		} `xml:"id"`
	} `xml:"idList"`

	AddressList struct {
		Text    string `xml:",chardata"`
		Address struct {
			Text    string `xml:",chardata"`
			Uid     string `xml:"uid"`
			Country string `xml:"country"`
		} `xml:"address"`
	} `xml:"addressList"`

	NationalityList struct {
		Text        string `xml:",chardata"`
		Nationality struct {
			Text      string `xml:",chardata"`
			Uid       string `xml:"uid"`
			Country   string `xml:"country"`
			MainEntry string `xml:"mainEntry"`
		} `xml:"nationality"`
	} `xml:"nationalityList"`
}

func (se *SdnEntry) Domain() domain.SDNRecord {
	return domain.SDNRecord{
		UID:       se.Uid,
		FirstName: se.FirstName,
		LastName:  se.LastName,
	}
}

func (se *SdnEntry) Pretty() string {
	format := `
	UID: %d
	First Name:%s
	Last Name %s
	`
	return fmt.Sprintf(format, se.Uid, se.FirstName, se.LastName)
}

type Info struct {
	Text        string `xml:",chardata"`
	PublishDate string `xml:"Publish_Date"`
	RecordCount string `xml:"Record_Count"`
}

type SdnResponse struct {
	XMLName           xml.Name   `xml:"sdnList"`
	Text              string     `xml:",chardata"`
	Xsi               string     `xml:"xsi,attr"`
	Xmlns             string     `xml:"xmlns,attr"`
	PublshInformation Info       `xml:"publshInformation"`
	SdnEntries        []SdnEntry `xml:"sdnEntry"`
}

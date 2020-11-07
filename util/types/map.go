package types

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type XMLMap map[string]interface{}

func (m XMLMap) GetString(key string) string {
	v := m[key]
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func (m XMLMap) GetInt(key string) int {
	v := m[key]
	if v == nil {
		return 0
	}
	if val, ok := v.(int); ok {
		return val
	}
	if val, err := strconv.ParseInt(fmt.Sprint(v), 10, 0); err == nil {
		return int(val)
	}
	return 0
}

func (m *XMLMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if len(*m) == 0 {
		*m = make(XMLMap)
	}

	var key, val string
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			key = t.Name.Local
		//case xml.EndElement:
		case xml.CharData:
			val = string(t)
		}
		//fmt.Printf("k:'%v' v:'%v' '%v'\n", key, val, []byte(val))

		if key != "" && !strings.HasPrefix(val, "\n") {
			(*m)[key] = val
		}
	}

	return nil
}

func (m XMLMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "xml" // 微信接口根为 xml
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for k, v := range m {
		err := e.EncodeElement(v, xml.StartElement{
			Name: xml.Name{
				Local: k,
				Space: start.Name.Space + start.Name.Space,
			},
			Attr: nil,
		})
		if err != nil {
			return err
		}
	}
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	return e.Flush()
}

type XMLCDATA string

func (cdata XMLCDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	type innerData struct {
		Value string `xml:",cdata"`
	}
	err = e.EncodeElement(innerData{Value: string(cdata)}, start)
	return
}

func XMLStartElement(local string) xml.StartElement {
	return xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: local,
		},
		Attr: nil,
	}
}

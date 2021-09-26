package types

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// XMLMap 微信 XML Map 类型
//
// 注意:
//   多层级如 <a><b>1</b></a> 会被解析为 map["a.b"] = 1
//   暂时不支持数组类别如 <a><b>1</b><b>2</b></a>, 2 会覆盖 1
type XMLMap map[string]interface{}

// GetString 获取 string 类型值
func (m XMLMap) GetString(key string) string {
	v := m[key]
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// GetInt 获取 int 类型值
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

// UnmarshalXML 解析 XML
func (m *XMLMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if len(*m) == 0 {
		*m = make(XMLMap)
	}

	var stack []string

	var key, val, last string
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
			stack = append(stack, key)
			last = "start"
		case xml.EndElement:
			key = t.Name.Local
			last = "end"
			if key != "" && key != start.Name.Local {
				mk := strings.Join(stack, ".")
				(*m)[mk] = val
			}
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			if last == "start" {
				val = string(t)
			} else {
				val = ""
			}
			last = "chardata"
		}
		//fmt.Println("->", reflect.TypeOf(token))
		//fmt.Printf("t: '%v' k:'%v' s:'%v' v:'%v' '%v'\n", token, key, stack, val, []byte(val))
	}

	return nil
}

// MarshalXML 序列化为 XML
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

// XMLCDATA CDATA 类型
type XMLCDATA string

// MarshalXML 序列化XML
func (cdata XMLCDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	type innerData struct {
		Value string `xml:",cdata"`
	}
	err = e.EncodeElement(innerData{Value: string(cdata)}, start)
	return
}

// XMLStartElement 返回一个XML起始节点
func XMLStartElement(local string) xml.StartElement {
	return xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: local,
		},
		Attr: nil,
	}
}

package types

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

type xmlCDATAStruct struct {
	Value XMLCDATA
}

type xmlCDATAStruct1 struct {
	Value XMLCDATA `xml:",cdata"`
}

func TestXMLMapCDATA(t *testing.T) {
	type s1 struct {
		Value XMLCDATA
	}
	type s2 struct {
		Value XMLCDATA `xml:",cdata"`
	}
	type s3 struct {
		Value XMLCDATA `xml:"tag"`
	}

	datas := [][]interface{}{
		// value to marshal, expect, expect, need error
		{XMLCDATA("test"), "<XMLCDATA><![CDATA[test]]></XMLCDATA>", XMLCDATA("test"), false},
		{s1{"test"}, "<s1><Value><![CDATA[test]]></Value></s1>", s1{"test"}, false},
		{s2{"test"}, "<s2><![CDATA[test]]></s2>", s2{"test"}, false},
		{s3{"test"}, "<s3><tag><![CDATA[test]]></tag></s3>", s3{"test"}, false},
	}

	for _, item := range datas {
		v := item[0]
		marshalExp, unmarshalExp := item[1], item[2]
		needErr := item[3].(bool)
		val, err := xml.Marshal(v)
		if needErr {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, marshalExp, string(val))

		err = xml.Unmarshal(val, &v)
		assert.NoError(t, err)
		assert.Equal(t, unmarshalExp, v)

		//var x xmlCDATAStruct
		//err = xml.Unmarshal(b, &x)
		//t.Log(x, err)
	}
}

func TestXMLMapMarshal(t *testing.T) {
	datas := [][]interface{}{
		{XMLMap{}, "<XMLMap></XMLMap>"},
		{XMLMap{"A": 123}, "<XMLMap><A>123</A></XMLMap>"},
		{XMLMap{"A": "abc"}, "<XMLMap><A>abc</A></XMLMap>"},
		{XMLMap{"1": "abc", "2": true}, "<XMLMap><1>abc</1><2>true</2></XMLMap>"},
	}

	for _, item := range datas {
		v, expect := item[0], item[1].(string)
		b, err := xml.Marshal(v)
		assert.NoError(t, err)
		assert.Equal(t, expect, string(b))
	}

	buf := bytes.NewBuffer(nil)
	enc := xml.NewEncoder(buf)
	err := XMLMap{"A": 123}.MarshalXML(enc, xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "root",
		},
		Attr: nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, "<root><A>123</A></root>", string(buf.Bytes()))

}

func TestXMLMapUnmarshal(t *testing.T) {
	datas := [][]interface{}{
		{"<root><a>123</a><b><![CDATA[test]]></b></root>"},
	}
	for _, item := range datas {
		str := item[0].(string)

		var m XMLMap
		xml.Unmarshal([]byte(str), &m)
	}
}

package pay

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ParseNotify 解析支付通知
func (s *PayCombineService) ParseNotify(req *http.Request) (*CombineOrder, error) {
	notify, err := s.client.ParseNotify(req, nil)
	if err != nil {
		return nil, err
	}

	if notify.EventType != EventTypeTransactionSuccess {
		return nil, fmt.Errorf("通知类型错误")
	}

	rc := notify.Resource
	data, err := DecodeCiphertext(rc.Algorithm, rc.Ciphertext, rc.Nonce, rc.AssociatedData, s.client.apiKey)
	if err != nil {
		return nil, err
	}

	order := new(CombineOrder)
	err = json.Unmarshal(data, order)
	return order, err
}

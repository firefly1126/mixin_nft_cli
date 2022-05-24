package trident

import (
	"encoding/json"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/shopspring/decimal"
)

type MetaData struct {
	Creator    Creator    `json:"creator,omitempty"`
	Collection Collection `json:"collection,omitempty"`
	Token      Token      `json:"token,omitempty"`
	Checksum   Checksum   `json:"checksum,omitempty"`
}

type Creator struct {
	ID      string          `json:"id,omitempty"` // uuid, required
	Name    string          `json:"name,omitempty"`
	Royalty decimal.Decimal `json:"royalty,omitempty"` //作者分红比例（每次售卖都有分红）
}

type Collection struct {
	ID          string          `json:"id,omitempty"`   // uuid,required
	Name        string          `json:"name,omitempty"` //required
	Description string          `json:"description,omitempty"`
	Icon        Icon            `json:"icon,omitempty"`
	Split       decimal.Decimal `json:"split,omitempty"` //早期拥有者收益比例（每次售卖，每个早期拥有者都建获得收益）
}

type Token struct {
	ID          string                 `json:"id,omitempty"`   //required, token identifier, number is the best
	Name        string                 `json:"name,omitempty"` //required
	Description string                 `json:"description,omitempty"`
	Icon        Icon                   `json:"icon,omitempty"`
	Media       Media                  `json:"media,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

type Icon struct {
	Url  string `json:"url,omitempty"`
	Hash string `json:"hash,omitempty"`
}

type Media struct {
	Url  string `json:"url,omitempty"`
	Hash string `json:"hash,omitempty"`
	MiMe string `json:"mime,omitempty"`
}

type Checksum struct {
	Fields    []string `json:"fields,omitempty"` //required
	Algorithm string   `json:"algorithm"`        //required, sha3-256
}

type CreateMetaDataPayload struct {
	MetaData MetaData `json:"metadata"`
	MetaHash string   `json:"metahash"`
}

func CreateMetaData(token string, payload *CreateMetaDataPayload) (interface{}, error) {
	req := mixin.GetRestyClient().NewRequest().SetAuthToken(token)
	req.SetBody(payload)
	rsp, err := req.Post("https://thetrident.one/api/collectibles")
	if err != nil {
		return nil, err
	}

	var b interface{}
	if err := json.Unmarshal(rsp.Body(), &b); err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateMetaData() error {
	return nil
}

func GetMetaData(metaHash string) (*MetaData, error) {
	rsp, err := mixin.GetRestyClient().NewRequest().Get("https://thetrident.one/api/collectibles/" + metaHash)
	if err != nil {
		return nil, err
	}

	var metaData MetaData
	if err := json.Unmarshal(rsp.Body(), &metaData); err != nil {
		return nil, err
	}

	return &metaData, nil
}

func GetOrders() error {
	return nil
}

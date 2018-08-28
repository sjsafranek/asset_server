package main

import (
	"encoding/json"
	"time"
)

var AssetSQL = `
	'{'||
		'"asset_id": "'  || asset_id  ||'",'||
		'"file_name": "' || file_name ||'",'||
		'"is_deleted": ' || (
			CASE
				WHEN is_deleted = 1 THEN 'true'
				ELSE 'false'
			END
		) ||','||
		'"create_at": "' || create_at ||'",'||
		'"update_at": "' || update_at ||'"'
	|| '}'
`

type Asset struct {
	AssetId   string    `json:"asset_id"`
	FileName  string    `json:"file_name"`
	IsDeleted string    `json:"is_deleted"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
}

func (self *Asset) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), self)
}

func (self Asset) Marshal() (string, error) {
	b, err := json.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), nil
}

package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ChNumber int

const (
	ChNum1 ChNumber = 1
	ChNum2 ChNumber = 2
	ChNum3 ChNumber = 3
	ChNum4 ChNumber = 4
)

type ChNum struct {
	Number ChNumber `json:"ChNum,omitempty"`
}

func (c ChNum) String() string {
	return strconv.Itoa(int(c.Number))
}

func (c *ChNum) UnmarshalJSON(data []byte) error {
	type Alias ChNum
	aux := struct {
		Number string `json:"ChNum"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Number == "" {
		return nil
	}

	var err error
	var num int
	if num, err = strconv.Atoi(aux.Number); err != nil {
		return fmt.Errorf("error parsing ChNum: %w", err)
	}

	c.Number = ChNumber(num)

	return nil
}

func (c ChNum) MarshalJSON() ([]byte, error) {
	type Alias ChNum
	return json.Marshal(&struct {
		Number string `json:"ChNum"`
		*Alias
	}{
		Number: c.String(),
		Alias:  (*Alias)(&c),
	})
}

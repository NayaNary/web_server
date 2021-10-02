package data

import (
	"encoding/json"
	"fmt"
)

type DataProcessing struct {
	Pages  map[uint64][]InputData
	LastId uint64
}

func NewDataProcessing() *DataProcessing {
	return &DataProcessing{
		Pages:  make(map[uint64][]InputData),
		LastId: 1,
	}
}

func (p DataProcessing) PreProcessing(input []byte) (answer []byte) {
	var iD []InputData
	//fmt.Println("Body: ", input)

	err := json.Unmarshal(input, &iD)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var answer_struct ResultInput
	answer_struct.Result = false
	if len(iD) != 0 {
		page := p.Pages[p.LastId]
		page = append(page, iD...)
		p.Pages[p.LastId] = page

		answer_struct.Result = true
		answer_struct.PageId = p.LastId

	}
	answer, err = json.Marshal(answer_struct)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return
}

func (p DataProcessing) Page(id uint64) (answer []byte){
	page := p.Pages[p.LastId]
	trans := make(map[string]OutputData)
	var data OutputData

     for i:=0;i<len(page); i++ {
		 data.LastTrade = page[i].LastTradePrice
		 data.Price = page[i].Price24h
		 data.Volume = page[i].Volume24h
		 trans[page[i].Symbol] = data
	 }
 		answer, _ = json.Marshal(trans)
	 return
}
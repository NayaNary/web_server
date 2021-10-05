package data

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"test.task/src/db"
)

type DataProcessing struct {
	mux    sync.Mutex
	Pages  map[uint64][]InputData
	LastId uint64
	objDB  db.Db
	timer  *time.Ticker
}

func NewDataProcessing() *DataProcessing {
	connParam := "user=postgres password=123456789 dbname=test_task sslmode=disable"
	database := db.NewConnect("postgres", connParam)
	dp := &DataProcessing{
		Pages:  make(map[uint64][]InputData),
		LastId: 1,
		objDB:  *database,
	}
	dp.timer = dp.StartWritePages()
	return dp
}

func (p *DataProcessing) PreProcessing(input []byte) (answer []byte) {
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
		p.mux.Lock()
		p.Pages[p.LastId] = page
		p.mux.Unlock()
		answer_struct.Result = true
		answer_struct.PageId = p.LastId

	}
	answer, err = json.Marshal(answer_struct)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return
}

func (p *DataProcessing) Page(id uint64) (answer []byte) {
	p.mux.Lock()
	page := p.Pages[p.LastId]
	p.mux.Unlock()
	trans := make(map[string]OutputData)
	var data OutputData

	for i := 0; i < len(page); i++ {
		data.LastTrade = page[i].LastTradePrice
		data.Price = page[i].Price24h
		data.Volume = page[i].Volume24h
		trans[page[i].Symbol] = data
	}
	answer, _ = json.Marshal(trans)
	return
}

func (p *DataProcessing) StartWritePages() (t *time.Ticker) {

	t = time.NewTicker(10 * time.Second)
	go func() {
		defer t.Stop()
		for now := range t.C {
			fmt.Println(now, time.Now().Format("2006-01-02 15:04:05"))
			rows := [...]string{"page", "create_at"}
			p.mux.Lock()
			fmt.Println("id:", p.LastId)
			fmt.Println("pages:", p.Pages)
			page := p.Pages[p.LastId]
			fmt.Println("befor delet", page)
			delete(p.Pages, p.LastId)
			fmt.Println("after delet", page)
			p.mux.Unlock()
			if len(page) != 0 {
				vi := make([]interface{},0)
				// var i []InputData
				// for k:=0; k<len(page);k++{
				// 	i = append(i, page[k])
				// 	kn:= k+1 
				// 	if kn < len(page){
				// 		i = append(i, ",")	
				// 	}
				
				// }
				dn, err:= json.Marshal(page)
				if err != nil {
					fmt.Println("marsh:", err)
				}
				
				
				vi = append(vi, string(dn))
				vi = append(vi, time.Now().Format("2006-01-02 15:04:05"))
				fmt.Println("vals:", vi)
				p.objDB.Insert("pages", rows[0:2], vi)
				if p.objDB.LastError() != nil {
					fmt.Println("insert:", p.objDB.LastError().Error())
				}
			}

		}
	}()
	return
}

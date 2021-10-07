package data

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"test.task/src/config"
	"test.task/src/db"
)

func New() *DataProcessing {
	env := config.New()

	var connParam string
	switch env.Db.NameDriver {
	case "postgres":
		connParam = fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v",
			env.Db.User, env.Db.Password, env.Db.DbName, env.Db.SslMode)
	case "mysql":
		connParam = fmt.Sprintf("%v:%v@/%v",
			env.Db.User, env.Db.Password, env.Db.DbName)
	default:
		connParam = fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v",
			env.Db.User, env.Db.Password, env.Db.DbName, env.Db.SslMode)
	}

	database := db.NewConnect(env.Db.NameDriver, connParam)

	dp := &DataProcessing{
		Pages:  make(map[uint64][]InputData),
		LastId: 1,
		objDB:  *database,
		Conf: env,
	}

	dp.Timer = dp.StartWritePages()

	return dp
}

func (p *DataProcessing) PreProcessing(input []byte) (answer []byte) {
	var inputData []InputData

	err := json.Unmarshal(input, &inputData)
	if err != nil {
		log.Println("Error:", err)
	}

	var bodyAnswer ResultInput
	bodyAnswer.Result = false
	if len(inputData) != 0 {
		page := p.Pages[p.LastId]
		page = append(page, inputData...)
		p.mux.Lock()
		p.Pages[p.LastId] = page
		p.mux.Unlock()
		bodyAnswer.Result = true
		bodyAnswer.PageId = p.LastId

	}
	answer, err = json.Marshal(bodyAnswer)
	if err != nil {
		log.Println("Error:", err)
	}
	return
}

func (p *DataProcessing) Page(id uint64) (answer []byte) {
	row := [...]string{"page"}
	where := fmt.Sprintf("where id = %v", id)
	page := p.objDB.Select("pages", row[0:1], where)
	defer page.Close()

	outputData := make(map[string]OutputData)
	var data OutputData
	if page.Next() {
		var dataDB []uint8
		err := page.Scan(&dataDB)
		if err != nil {
			log.Println("err scan:", err)
		}
		var tran []InputData
		err = json.Unmarshal([]byte(dataDB), &tran)
		if err != nil {
			log.Println("err unmarshal:", err)
		}
		for i := 0; i < len(tran); i++ {
			data.LastTrade = tran[i].LastTradePrice
			data.Price = tran[i].Price24h
			data.Volume = tran[i].Volume24h
			outputData[tran[i].Symbol] = data
		}
	}

	answer, _ = json.Marshal(outputData)
	return
}

func (p *DataProcessing) StartWritePages() (t *time.Ticker) {

	t = time.NewTicker(time.Duration(p.Conf.Db.TimeWrite) * time.Second)
	go func(t *time.Ticker) {
		defer t.Stop()
	    p.WritePages()
	}(t)
	return
}

func (p *DataProcessing) WritePages(){
	for now := range p.Timer.C {
		log.Println(now, time.Now().Format("2006-01-02 15:04:05"))
		rows := [...]string{"page", "create_at"}

		p.mux.Lock()
		page := p.Pages[p.LastId]
		p.mux.Unlock()

		if len(page) != 0 {
			values := make([]interface{}, 0)
			dataWrite, err := json.Marshal(page)
			if err != nil {
				fmt.Println("marsh:", err)
			}

			values = append(values, string(dataWrite))
			values = append(values, time.Now().Format("2006-01-02 15:04:05"))
			p.objDB.Insert("pages", rows[0:2], values)
			if p.objDB.LastError() == nil {
				p.mux.Lock()
				delete(p.Pages, p.LastId)
				p.mux.Unlock()
			} else {
				log.Println("insert:", p.objDB.LastError().Error())
			}
		}

	}
}

func (p *DataProcessing) CountPages() (answer []byte) {

	row := [...]string{"id"}
	where := "ORDER BY id DESC LIMIT 1"
	page := p.objDB.Select("pages", row[0:1], where)
	defer page.Close()

	if p.objDB.LastError() != nil {
		resultGet.Message = p.objDB.LastError().Error()
	}
	if page.Next() {
		var amount uint64
		err := page.Scan(&amount)
		if err != nil {
			log.Println("err scan:", err)
		}
		resultGet.AmountPages = amount
	}
	var err error
	answer, err = json.Marshal(resultGet)
	if err != nil {
		log.Println("err scan:", err)
	}
	return
}

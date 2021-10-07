package web

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"test.task/src/data"
)

type WebRouts struct {
	R     *mux.Router
	ProcData *data.DataProcessing
}

func New() *WebRouts {
	return &WebRouts{
		R:     mux.NewRouter(),
		ProcData: data.New(),
	}
}
func (w *WebRouts) StartServer() (srv *http.Server) {

	addrSrv := w.ProcData.Conf.Web.Host + ":" + w.ProcData.Conf.Web.Port

	srv = &http.Server{
		Handler:      w.RoutConnect(),
		Addr:         addrSrv,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	
	go func(srv *http.Server) {
		log.Println("Запущен сервер: ", srv.Addr)
		log.Fatal(srv.ListenAndServe())
	}(srv)

	return
}

func (w *WebRouts) RoutConnect() *mux.Router {

	w.R.HandleFunc("/", w.inputData).Methods("POST")
	w.R.HandleFunc("/amount", w.amountPages)
	w.R.HandleFunc("/page/{id}", w.outputData)
	http.Handle("/", w.R)
	return w.R
}

func (wr *WebRouts) inputData(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)

	answer := wr.ProcData.PreProcessing(body)

	w.Write(answer)
}
func (wr *WebRouts) amountPages(w http.ResponseWriter, r *http.Request) {
	w.Write(wr.ProcData.CountPages())
}
func (wr *WebRouts) outputData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlValue := vars["id"]
	id, _ := strconv.ParseUint(urlValue, 10, 0)
	answer := wr.ProcData.Page(id)
	w.Write(answer)
}

func readBody(r *http.Request) (bs []byte) {
	bs = make([]byte, 255)
	var bodyFul []byte
	for {
		n, errRead := r.Body.Read(bs)
		if errRead == io.EOF {
			for i := 0; i < n; i++ {
				bodyFul = append(bodyFul, bs[i])
			}
			break
		}
		for i := 0; i < n; i++ {
			bodyFul = append(bodyFul, bs[i])
		}
	}
	errClose := r.Body.Close()
	if errClose != nil {
		log.Println(errClose, "ReadDataRequest")
	}

	return bodyFul
}

package web

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"test.task/src/data"
)

type WebRouts struct {
	R     *mux.Router
	Pages data.DataProcessing
}

func NewWebRouts() *WebRouts {
	return &WebRouts{
		R:     mux.NewRouter(),
		Pages: *data.NewDataProcessing(),
	}
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
	answer_byte := wr.Pages.PreProcessing(body)

	w.Write(answer_byte)
  // fmt.Println("Pages: ", wr.Pages)
}
func (wr *WebRouts) amountPages(w http.ResponseWriter, r *http.Request) {
	w.Write(wr.Pages.CountPages())
	fmt.Println("amount pages")
}
func (wr *WebRouts) outputData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url:",r.URL)
	vars:= mux.Vars(r)
	urlValue := vars["id"]
	fmt.Println("url in:",urlValue)
	id, _ := strconv.ParseUint(urlValue, 10, 0)
	fmt.Println("id in:", id)
	answer := wr.Pages.Page(id)
	w.Write(answer)
	fmt.Println("output data")
}

func readBody(r *http.Request) (bs []byte) {
	bs = make([]byte, 255)
	var bodyFul []byte
	for {
		N, ErrRead := r.Body.Read(bs)
		if ErrRead == io.EOF {
			for i := 0; i < N; i++ {
				bodyFul = append(bodyFul, bs[i])
			}
			break
		}
		for i := 0; i < N; i++ {
			bodyFul = append(bodyFul, bs[i])
		}
	}
	errClose := r.Body.Close()
	if errClose != nil {
		fmt.Println(errClose, "ReadDataRequest")
	}

	return bodyFul
}

package multiroute

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const groupPattern = "/get/"

func Test_NewMultiRouter(t *testing.T) {
	mux := http.NewServeMux()

	n := NewMultiRouter(groupPattern, notFoundHandler)
	n.AddRoute(`^time/$`, timeHandler)
	n.AddRoute(`^hello\w{1,5}$`, echoHandler)
	mux.Handle(groupPattern, n)

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func(serv *http.Server) {
		time.Sleep(2 * time.Second)

		err := sendRequest()
		if err != nil {
			t.Fail()
			log.Println("sendRequest, err:", err)
		} else {
			log.Println("sendRequest success")
		}

		err = sendRequest2()
		if err != nil {
			t.Fail()
			log.Println("sendRequest2, err:", err)
		} else {
			log.Println("sendRequest2 success")
		}

		if err := serv.Shutdown(context.Background()); err != nil {
			t.Fail()
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}(&serv)

	err := serv.ListenAndServe()
	if err != nil {
		log.Println("ListenAndServe, error, ", err)
	}
}

func sendRequest() error {
	req, err := http.NewRequest("GET", "http://localhost:8080/get/time/", nil)

	if err != nil {
		log.Println(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("response:", string(f))
	if resp.StatusCode != 200 {
		return errors.New("err status, " + resp.Status)
	}
	return nil
}

func sendRequest2() error {
	req, err := http.NewRequest("GET", "http://localhost:8080/get/helloabc2", nil)

	if err != nil {
		log.Println(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("response:", string(f))
	if resp.StatusCode != 200 {
		return errors.New("err status, " + resp.Status)
	}
	return nil
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	m := time.Now().Format(time.RFC3339)

	length := len(m)
	x := 0
	y := length
	z := length
	sub := []byte(m)[x:y:z]

	w.Write(sub)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	lenURL := len(r.URL.Path)
	x := len(groupPattern)
	y := lenURL
	z := lenURL
	sub := []byte(r.URL.Path)[x:y:z]

	w.Write(sub)
}

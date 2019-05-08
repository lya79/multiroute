# Dynamic HTTP route handler
#### Features
- 基於 net/http標準庫套件
- 提供群組路由
- 提供動態路由
#### Sample
```golang

const groupPattern = "/get/"

func main(){
    mux := http.NewServeMux()

	n := NewMultiRouter(groupPattern, notFoundHandler)
	n.AddRoute(`^time/$`, timeHandler)          // ex: http://localhost:8080/get/time/
	n.AddRoute(`^hello\w{1,5}$`, echoHandler)   // ex: http://localhost:8080/get/helloabc1
	mux.Handle(groupPattern, n)

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
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
```

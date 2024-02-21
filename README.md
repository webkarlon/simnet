# SIMNET - обертка для работы с net/http
+ Имеет удобный роутер, с указанием http методов.
+ Есть возможность задавать динамические url
+ Мидалвары можно указывать списком, а не заворачивать один в другой
+ Поддержка мультипаттерна
+ Есть возможность выбора HTTP протокола
+ Сервер настраивается стандартными средствами net/http

### Пример простого сервера свыше перечисленными возможностями.

```go

package main

import (
	"fmt"
	simnet "github.com/webkarlon/simnet"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	middleware := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this middleware\n"))
	}

	myServer := simnet.NewServer(&simnet.Server{
		PortHTTP:        8888,
		ListenAddress:   "localhost",
		ShutdownTimeout: 2,
	})

	myServer.AddRouter(http.MethodGet, "/", middleware, Home)
	myServer.AddRouter(http.MethodPost, "/", middleware, Home)
	myServer.AddRouter(http.MethodGet, "/person/:name/:age", middleware, Person)
	myServer.AddRouter(http.MethodGet, "/*file", middleware, File)
	myServer.AddRouter(http.MethodGet, "/upload/:id/*file", middleware, Upload)

	go myServer.Start()

	osSignalsCh := make(chan os.Signal, 1)
	signal.Notify(osSignalsCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-osSignalsCh

	err := myServer.Stop()
	if err != nil {
		panic(err)
	}

}

func Upload(w http.ResponseWriter, r *http.Request) {
	pattern := simnet.GetParam(r, "pattern")
	file := simnet.GetParam(r, "file")
	id := simnet.GetParam(r, "id")
	w.Write([]byte(fmt.Sprintf("pattern:%s\nid:%s\nfile:%s", pattern, id, file)))
}

func Person(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pattern: " + simnet.GetParam(r, "pattern") + "\n"))
	w.Write([]byte("Name: " + simnet.GetParam(r, "name") + "\n"))
	w.Write([]byte("Age: " + simnet.GetParam(r, "age") + "\n"))
}

func File(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pattern: " + simnet.GetParam(r, "pattern") + "\n"))
	w.Write([]byte("File: " + simnet.GetParam(r, "file")))
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home " + r.Method))
}

```

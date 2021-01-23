package ms

import (
		"io"
		"os"
		"fmt"
		"time"
		"mime"
		"path"
		"net/http"
		"context"
		"strings"
		"os/signal"
		"math/rand"
		"github.com/gorilla/mux"
		"github.com/sirupsen/logrus"
		//"github.com/go-redis/redis/v8"
		"github.com/t-tomalak/logrus-easy-formatter"
)

type Mediaserver struct {
		router *mux.Router
		server *http.Server
}

var logger *logrus.Logger

func init_logger() *logrus.Logger {
		log_file, err := os.OpenFile("mediaserver.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
				fmt.Println(err)
				os.Exit(1)
		}
		return &logrus.Logger{
				Out: log_file,
				Level: logrus.DebugLevel,
				Formatter: &easy.Formatter{
						TimestampFormat: "2017-08-01 16:51:23",
						LogFormat: "[%lvl%]: %time% - %msg%\n",
				},
		}
}

func New(port, host string) *Mediaserver {
		addr := host+port
		r := mux.NewRouter()
		r.HandleFunc("/files/get/{file}", return_file)
		r.HandleFunc("/files/info/{file}", info_about_file)
		r.HandleFunc("/files/upload", upload_file)
		r.HandleFunc("/files/delete/{file}", delete_file)
		return &Mediaserver{
				router: r,
				server: &http.Server{
						Handler: r,
						Addr: addr,
						WriteTimeout: 15*time.Second,
						ReadTimeout: 15*time.Second,
				},
		}
}

func (ms *Mediaserver) Run() {
		logger = init_logger()
		go func() {
				if err := ms.server.ListenAndServe(); err != nil {
						time.Sleep(7*time.Second)
						logger.Fatal(err)
				}
		}()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ms.server.Shutdown(ctx)
		logger.Println("the server is shutted down")
		os.Exit(0)
}

func return_file(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file_name := "./media/"+vars["file"]
		logger.Debug(file_name)
		file, err := os.Open(file_name)
		if err != nil {
				w.Header().Set("Content-Type", "aplication/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":404, "msg":"file not found"}`))
		} else {
				defer file.Close()
				content_type := mime.TypeByExtension(path.Ext(file_name))
				stat, err := file.Stat()
				if err != nil {
						logger.Warn(err)
				}
				bytes := make([]byte, stat.Size())
				w.Header().Set("Content-Type", content_type)
				w.WriteHeader(http.StatusOK)
				file.Read(bytes)
				w.Write(bytes)
		}
}

func upload_file(w http.ResponseWriter, r *http.Request) {
		content_type := strings.Join(r.Header["Content-Type"], "")
		logger.Debug(content_type)
		if strings.Contains(content_type, "multipart/form-data") {
				src, _, err := r.FormFile("my-file")
				if err != nil {
						logger.Warn(err)
				}
				f, err := os.Create("./media/"+gen_file_name())
				if err != nil {
						logger.Warn(err)
				}
				defer f.Close()

				io.Copy(f, src)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":200, "msg": "file was successfuly downloaded"}`))
		} else {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
				<!DOCTYPE html>
						<html lang="en">
								<head>
										<meta charset="UTF-8" />
										<meta name="viewport" content="width=device-width, initial-scale=1.0" />
										<meta http-equiv="X-UA-Compatible" content="ie=edge" />
										<title>Document</title>
								</head>
								<body>
										<form
      enctype="multipart/form-data"
      action="http://localhost:8002/files/upload"
      method="post"
    >
												<input type="file" name="my-file" />
												<input type="submit" value="upload" />
										</form>
								</body>
						</html>`))
	}
}

func delete_file(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file_name := "./media/"+vars["file"]
		if _, err := os.Stat(file_name); os.IsNotExist(err) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":404, "msg":"file does not exist"}`))
		} else {
				os.Remove(file_name)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": 200, "msg": "file was successfuly deleted"}`))
		}
}

func info_about_file(w http.ResponseWriter, r *http.Request) {

}

func gen_file_name() string {
		symbols := "ABCDEFGHIJKLMNOPQRSTUVWXYZabsdefghijklmnopqrstuvwxyz1234567890"
		list_of_symbols := strings.Split(symbols, "")
		rand.Seed(time.Now().UnixNano())
		var generated_name []string
		for i := 0; i < 16; i++ {
				generated_name = append(generated_name, list_of_symbols[rand.Intn(len(list_of_symbols))])
		}
		return strings.Join(generated_name, "")
}

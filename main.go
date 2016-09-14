package main
import (
	"fmt"
	"net/http"
	"uploadbreak/upload"
)

func main(){
	http.HandleFunc("/upload", upload.UpLoad)
	http.HandleFunc("/hello", upload.HelloServer)
    err := http.ListenAndServe(":80", nil)
    if err != nil {
       fmt.Println("ListenAndServe: ", err)
    }
}
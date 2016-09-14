package upload

import (
	"io"
	"strconv"
	"os"
	"fmt"
	"net/http"
)
var UpTable map[string]string

func init(){
	UpTable = make(map[string]string)
}


//检查文件是否存在
func checkFileExist(file string)( bool,int64) {
	if f,err:=os.Stat(file) ;os.IsNotExist(err){
		return false,-1
	}else{
		return true,f.Size()
	}
}

//http://blog.sina.com.cn/s/blog_b37612030101drca.html
func UpLoad(w http.ResponseWriter,r *http.Request){
	fmt.Println("map:",UpTable)
	fmt.Println("header:", r.Header) //获取请求的方法
	path := r.Header.Get("file")
	hashcode := r.Header.Get("hashcode")
	fmt.Println("path:",path,"hashcode:",hashcode)
		
    if r.Method == "GET" {
	   	
		if code,ok:=UpTable[path];ok{//已经上传过
            fmt.Println("已经上传过")		
			if code==hashcode{	//上传的文件一致
				if 	exist,size:= checkFileExist(path);exist {
					fmt.Fprintf(w,"写了:%d",size)
					w.Header().Set("offset",strconv.FormatInt(size,10))//下次上传位置
					return 
				}
			}
		}
		w.Header().Set("offset","0")//下次上传位置，起始位置
    } else {
		
        r.ParseMultipartForm(32 << 20) //32M
        file, handler, err := r.FormFile("userfile")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
 		
		UpTable[path]=hashcode

        f, err := os.OpenFile(handler.Filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        writen,err:=io.Copy(f, file)
		if err!=nil{
			fmt.Println(err)
		}
		fmt.Fprintf(w,"写了:%d",writen)
    }
	
}



// hello world, the web server
func HelloServer(w http.ResponseWriter, r *http.Request) {
    // 上传页面
    w.Header().Add("Content-Type", "text/html")
    w.WriteHeader(200)
    html := `
<form enctype="multipart/form-data" action="/upload" method="Get">
    Send this file: <input name="userfile" type="file" />
    <input type="submit" value="Send File" />
</form>
`
	io.WriteString(w, html)

}



/* 测试
curl 192.168.1.71/hello -c cookie.txt -v
curl 192.168.1.71/hello?hasecode=123456 -c cookie.txt -v
curl 192.168.1.71/upload?hasecode=123456 -c cookie.txt -v
ll
curl 192.168.1.71/upload?hasecode=123456 -T s -c cookie.txt -v
curl 192.168.1.71/upload?hasecode=123456  -F "userfile=@./s;type=application/octet-stream" -c cookie.txt

curl 192.168.1.71/upload -H 'hashcode:123456' -c cookie.txt -v

*/
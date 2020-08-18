package conf

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Config struct {
	FilePath string
	MD5 string
	LastModifyTime time.Time
}


func (c Config)parse(){
	getMd5 := c.getMD5()
	if getMd5 != c.MD5{
		log.Println("重载配置文件")
	}
}

func (c Config)getMD5()(md5s string){
	file, err := os.Open(c.FilePath)
	if err != nil {
		log.Println("加载配置文件失败,",err)
		return
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println("", err)
		return
	}
	return string(hash.Sum(nil))
}
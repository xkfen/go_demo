package ftp_server

import (
	"fmt"
	"io"
	"os"
	"gcoresys/common"
	"gopkg.in/dutchcoders/goftp.v1"
)

var gftpObj *gftp

type gftpConfig struct {
	Username string `json:"username" fname:"用户名"`
	Password string `json:"password" fname:"密码"`
	Host     string `json:"host" fname:"访问地址"`
	Port     string `json:"port" fname:"端口"`
}

func getGftpConfig() (cnf *gftpConfig) {
	//cnf = &gftpConfig{
	//	Username: "xingye",
	//	Password: "xingye",
	//	Host:     "qiyuanfin.picp.io",
	//	Port:     "2221",
	//}
	//if common.GetUseDocker() == 2 {
	//	cnf = &gftpConfig{
	//		Username: "xingye",
	//		Password: "xingye",
	//		Host:     "qiyuanfin.picp.io",
	//		Port:     "2221",
	//	}
	//}
	cnf = &gftpConfig{
		Username: "910002",
		Password: "cib@910002",
		Host:     "117.144.184.152",
		Port:     "22",
	}
	if common.GetUseDocker() == 2 {
		cnf = &gftpConfig{
			Username: "910002",
			Password: "cib@910002",
			Host:     "117.144.184.152",
			Port:     "22",
		}
	}
	return
}

func NewGftp() (g *gftp, err error) {
	if gftpObj == nil {
		cnf := getGftpConfig()
		ftp, err := goftp.Connect(fmt.Sprintf("%s:%s", cnf.Host, cnf.Port))
		if err != nil {
			return nil, err
		}
		if err = ftp.Login(cnf.Username, cnf.Password); err != nil {
			return nil, err
		}
		gftpObj = &gftp{
			FtpObj:   ftp,
			Password: cnf.Password,
			Username: cnf.Username,
			Host:     cnf.Host,
			Port:     cnf.Port,
		}
	}
	return gftpObj, nil
}


type gftp struct {
	FtpObj   *goftp.FTP
	Password string
	Username string
	Host     string
	Port     string
}

type GftpDownloadHandlerFunc func(data []byte) error

// 文件上传到ftp服务器
func (g *gftp) UploadFile(filepath, outpath string) (err error) {
	var file *os.File
	if file, err = os.Open(filepath); err != nil {
		return
	}
	return g.FtpObj.Stor(outpath, file)
}

// 从ftp服务器下载文件
func (g *gftp) DownloadFile(filepath string, f GftpDownloadHandlerFunc) (err error) {
	_, err = g.FtpObj.Retr(filepath, func(r io.Reader) (readErr error) {
		bt := make([]byte, 20*1024*1024)
		if _, readErr = r.Read(bt); readErr != nil {
			return
		}
		return f(bt)
	})
	return
}

func (g *gftp) Close() {
	g.FtpObj.Close()
}
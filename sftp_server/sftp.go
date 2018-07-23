package sftp_server

import (
	"os"
	"golang.org/x/crypto/ssh"
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/pkg/sftp"
	"gopkg.in/dutchcoders/goftp.v1"
	"time"
	"path"
)

func ConnectSftp() (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	cnf := getGftpConfig()
	auth = append(auth, ssh.Password(cnf.Password))

	clientConfig = &ssh.ClientConfig{
		User:            cnf.Username,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", cnf.Host, cnf.Port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

// 上传文件到sftp服务器
func Upload2SftpServer(sftpClient *sftp.Client, localFilePath string, remotePath , busicode string) error{
	// 读取文件
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		//logger.Error("err", "UploadFile2SftpServer###open err", err.Error())
		return err
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)
	index := strings.Index(remoteFileName, ".")
	tmpName := remoteFileName[0 : index]

	remoteFileName = tmpName + "(" + busicode + ")"   + path.Ext(remoteFileName)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		//logger.Error("err", "UploadFile2SftpServer###sftpClient.Create error", err.Error())
		return err
	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		//logger.Error("err", "UploadFile2SftpServer###ReadAll error", err.Error())
		return err
	}
	dstFile.Write(ff)

	return nil
}

// 上传文件夹到sftp服务器
func UploadDir2SftpServer(sftpClient *sftp.Client, localPath string, remotePath, busicode string) error {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		//logger.Error("err", "UploadDir2SftpServer###read dir list fail", err.Error())
		return err
	}
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			sftpClient.Mkdir(remoteFilePath)
			UploadDir2SftpServer(sftpClient, localFilePath, remoteFilePath, busicode)
		} else {
			Upload2SftpServer(sftpClient, path.Join(localPath, backupDir.Name()), remotePath, busicode)
		}
	}

	//logger.Info(" copy directory to remote server finished!")
	return nil
}



// 从远程sftp服务器下载文件
func DownloadFromSftp(sftpClient *sftp.Client, remoteFilePath, localDir string) error {
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err := ConnectSftp()
	if err != nil {
		//logger.Error("err", "DownloadFromSftp###connect err", err.Error())
		return err
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		//logger.Error("err", "DownloadFromSftp###open file err", err.Error())
		return err
	}
	defer srcFile.Close()

	var localFileName = path.Base(remoteFilePath)
	dstFile, err := os.Create(path.Join(localDir, localFileName))
	if err != nil {
		//logger.Error("err", "DownloadFromSftp###create err", err.Error())
		return err
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		//logger.Error("err", "DownloadFromSftp####write to file err", err.Error())
		return err
	}

	//logger.Info("copy file from remote server finished!")

	return nil
}

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
	//if common.GetUseDocker() == 2 {
	//	cnf = &gftpConfig{
	//		Username: "910002",
	//		Password: "cib@910002",
	//		Host:     "117.144.184.152",
	//		Port:     "22",
	//	}
	//}
	return
}

type gftp struct {
	FtpObj   *goftp.FTP
	Password string
	Username string
	Host     string
	Port     string
}
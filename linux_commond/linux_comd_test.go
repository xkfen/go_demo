package linux_commond

import (
	"testing"
	"os/exec"
	"fmt"
	"io/ioutil"

)

func TestGzFile(t *testing.T){
	cmd := exec.Command("/bin/bash", "-c", `df -lh`)
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return
	}
	fmt.Printf("stdout:\n\n %s", bytes)
}


func TestFileGz(t *testing.T){
	// tar cvf a.gz protobuf-3.6.0
	file := "/home/qydev/test.txt"
	//cmd := exec.Command("tar", "cvf", "1.gz", file)
	cmd := exec.Command("gzip", file)
	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return
	}
}

//func TestGz1(t *testing.T){
//	filePath := "/home/qydev/test.txt"
//	filePath = path.Join(filePath)
//	fileZip, err := os.Create(filePath)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	w := zip.NewWriter(fileZip)
//	defer w.Close()
//
//}
package main

import (
	"flag"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/gomail.v2"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var (
	// sftp文件地址
	sftpFilePath string
	host         string
	port         int64
	username     string
	password     string
	mailTo       string
	mailToList   []string
)

func main() {
	var (
		message string
	)
	initParm()
	client := new(ClientConfig)
	err := client.CreateClient(host, port, username, password)
	if err != nil {
		panic(err)
	}

	if num := client.GetSFTPDirFile(sftpFilePath); num == 0 {
		message = "<h3>各位好:</h3><p>SFTP数据上报回执文件检查，无异常</p>"
	} else {
		message = fmt.Sprintf("<h3>各位好:</h3><p>SFTP数据上报回执文件检查，异常。 失败文件数: %d</p>", num)
	}

	for _, v := range strings.Split(mailTo, ",") {
		mailToList = append(mailToList, v)
	}

	if err = SendMail(mailToList, message); err != nil {
		panic(err)
	}

}

func initParm() {
	flag.StringVar(&sftpFilePath, "sftp", "/sftp/", "sftp文件地址")
	flag.StringVar(&host, "host", "192.168.232.128", "sftp地址")
	flag.Int64Var(&port, "port", 22, "sftp 端口")
	flag.StringVar(&username, "username", "root", "sftp 账号")
	flag.StringVar(&password, "password", "shinemo123", "sftp 密码")
	flag.StringVar(&mailTo, "mail", "changwenjie@shinemo.com,wangran@shinemo.com", "收件人地址")
	flag.Parse()
}

func SendMail(mailTo []string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress("changwenjie@shinemo.com", ""))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", "SFTP数据上报")
	m.SetBody("text/html", message)
	d := gomail.NewDialer("smtp.exmail.qq.com", 465, "changwenjie@shinemo.com", "wC4kCgGKBjuUjbk5")

	err := d.DialAndSend(m)
	return err
}

type ClientConfig struct {
	Host       string       // ip
	Port       int64        // 端口
	Username   string       // 用户名
	Password   string       // 密码
	sshClient  *ssh.Client  // ssh Client
	sftpClient *sftp.Client // sftp Client
	LastResult string       // 获取一次执行结果
}

// CreateClient 创建连接客户端
func (cliConf *ClientConfig) CreateClient(host string, port int64, username string, password string) error {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		err        error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username
	cliConf.Password = password
	cliConf.Port = port
	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)
	if sshClient, err = ssh.Dial("tcp", addr, &config); err != nil {
		return err
	}
	cliConf.sshClient = sshClient
	//此时获取了sshClient，下面使用sshClient构建sftpClient
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return err
	}
	cliConf.sftpClient = sftpClient
	return nil
}

// RunShell 执行Shell命令
func (cliConf *ClientConfig) RunShell(shell string) string {
	var (
		session *ssh.Session
		err     error
	)
	//获取session，这个session是用来远程执行操作的
	if session, err = cliConf.sshClient.NewSession(); err != nil {
		return err.Error()
	}
	//执行shell
	if output, err := session.CombinedOutput(shell); err != nil {
		return err.Error()
	} else {
		cliConf.LastResult = string(output)
	}
	return cliConf.LastResult
}

// Upload 上传文件
func (cliConf *ClientConfig) Upload(srcPath, dstPath string) error {
	srcFile, _ := os.Open(srcPath)                   //本地
	dstFile, _ := cliConf.sftpClient.Create(dstPath) //远程
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}
		_, err = dstFile.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

// Download 下载文件
func (cliConf *ClientConfig) Download(src, dst string) error {
	srcFile, _ := cliConf.sftpClient.Open(src) //远程
	dstFile, _ := os.Create(dst)               //本地
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	return nil
}

func (cliConf *ClientConfig) GetSFTPDirFile(dir string) int {

	var errFileNumber = 0
	readDir, err := cliConf.sftpClient.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, v := range readDir {
		if strings.Contains(v.Name(), getDateNow()) {
			if strings.Contains(v.Name(), "err") {
				errFileNumber++
			}
		}
	}
	return errFileNumber
}

func getDateNow() string {
	return time.Now().Format("200601")
}

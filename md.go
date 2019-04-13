package markdown
//MarkDown文件内容类
import (
	"github.com/NiuStar/utils"
	"fmt"
	"os/exec"
	"os"
	//"path/filepath"
	"io"
	"bufio"
	"strings"
	//"time"
	//"github.com/microcosm-cc/bluemonday"
	//"gopkg.in/russross/blackfriday.v2"
)

type MarkDown struct {
	name string
	path string

	content string
}

func NewMarkDown(name ,path string) *MarkDown {

	return &MarkDown{name:name,path:path}
}

func (md *MarkDown)WriteTitle(level int,content string) {

	for i := 0 ; i < level ; i++ {
		md.content += "#"
	}
	md.content += " " + content
}

func (md *MarkDown)WriteContent(content string) {
	md.content += content
}

func (md *MarkDown)WriteImportantContent(content string) {
	md.content += "**" + content + "**"
}


func (md *MarkDown)WriteCode(contents string,language string) {

	md.content += "```" + language +`
` + contents + `
` + "```" + `
`
}

func (md *MarkDown)WriteForm(contents [][]string) {

	for index,list := range contents {
		md.content += "\r\n| "
		c := ""
		for _,content := range list {
			c += content
			c += " |"
		}
		if len(c) == 0 {
			md.content += " |"
		} else {
			md.content += c
		}
		if index == 0 {
			c = ""
			for _ , _ = range list {
				c += " ---------- |"
			}
			if len(c) == 0 {
				md.content += " - |"
			} else {
				md.content += c
			}
		}
	}
}

func (md *MarkDown)Save() {

	fmt.Println("getCurrentPath:",getCurrentPath())
	if !utils.Mkdir(md.path) {
		fmt.Println("创建文件夹失败",md.path)
		return
	}
	data := "# " + md.name + `


` + md.content
	utils.WriteToFile(md.path + "/" + md.name + ".md",data)
	fmt.Println("写入：" + md.path + "/" + md.name + ".md")

	GOPATH := os.Getenv("GOPATH")
	//ok,out := execCommand("./","pandoc","-s","-f markdown","-t html5","--toc","--css=../css/markdownPad-github.css","./" + md.name + ".md","-o ./" + md.name + ".html","-B ../css/header.html")
	ok,out := execCommand("./","pandoc","-s","-f","markdown","-t","html","--toc","--css=" + GOPATH + "/src/github.com/NiuStar/server2/markdown/css/markdownPad-github.css",md.path + "/" + md.name + ".md","-o" +  md.name + ".html","-B",GOPATH + "/src/github.com/NiuStar/server2/markdown/css/header.html")
	//f, err := exec.Command("title", "123456").Output()
	if !ok {
		fmt.Println("出错了")
	}
	fmt.Println("output:",string(out))

	//time.Sleep(2 * time.Second)
	data = utils.ReadFileFullPath(md.path + "/" + md.name + ".html")
	data = strings.Replace(data,"code span.kw { color: #007020;","code span.kw { color: rgb(119, 0, 136);",-1)
	data = strings.Replace(data,"code span.dt { color: #902000;","code span.dt { color: rgb(119, 0, 136);",-1)
	data = strings.Replace(data,"code span.co { color: #60a0b0;","code span.co { color: rgb(170, 85, 0);",-1)
	data = strings.Replace(data,"code span.bu { }","code span.bu { color: rgb(34, 17, 153) }",-1)
	data = strings.Replace(data,"code span.st { color: #4070a0;","code span.st { color: rgb(170, 17, 17);",-1)

	fmt.Println("data:data:",data)
	WriteWithFileWrite(md.path + "/" + md.name + "_api.html",data)
	//fmt.Println("data:data:",data)

	//code span.kw { color: #007020;
	//pandoc -s -f markdown -t html --toc --css=./css/markdownPad-github.css test.md -o test.html  -B ./header.html
	//utils.WriteToFile(md.path + "/" + md.name + ".html",string(blackfriday.Run([]byte(data), blackfriday.WithNoExtensions())))
}

//使用os.OpenFile()相关函数打开文件对象，并使用文件对象的相关方法进行文件写入操作
//清空一次文件
func WriteWithFileWrite(name,content string){
	fileObj,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_TRUNC,0644)
	if err != nil {
		fmt.Println("Failed to open the file",err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()
	if _,err := fileObj.WriteString(content);err == nil {
		fmt.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.",content)
	}
	contents := []byte(content)
	if _,err := fileObj.Write(contents);err == nil {
		fmt.Println("Successful writing to thr file with os.OpenFile and *File.Write method.",content)
	}
}

func getCurrentPath() string {
	s, _ := exec.LookPath(os.Args[0])
	//checkErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

func execCommand(dir string,commandName string, params ...string) (bool, string) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("cmd.Env:",cmd.Env)
	cmd.Dir = dir
	//显示运行的命令
	fmt.Println(cmd.Args)

	/*stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println("err: ",err)
		return false,err.Error()
	}*/

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("cmd.StdoutPipe : ", err)
		return false, err.Error()
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("cmd.Start() err: ", err)
		return false, err.Error()
	}
	//fmt.Println("stdout:",getOutput(stdout))
	//fmt.Println("stderr:",getOutput(stderr))
	//go io.Copy(os.Stdout, stdout)
	//go io.Copy(os.Stderr, stderr)
	result_err := getOutput(stderr)
	cmd.Wait()
	return true, result_err
}

func getOutput(out io.ReadCloser) string {
	reader := bufio.NewReader(out)
	var result string
	//实时循环读取输出流中的一行内容
	for {
		line, _, err2 := reader.ReadLine()
		//fmt.Println("OK:", ok)
		if err2 != nil || io.EOF == err2 {
			//fmt.Println("reader.ReadLine:", err2)
			break
		}
		result += string(line) + "\r\n"
		//fmt.Println("line:",string(line))
	}
	return result
}

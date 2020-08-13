package main

import (
	"flag"
	"fmt"
	"github.com/Luxurioust/excelize"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	src string
	slowLog string
	excelName string
)

type SlowLogInfo struct {
	Count string
	Time string
	Lock string
	Rows string
	SQL string
}

func ExecCmd(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func main() {
	flag.StringVar(&slowLog,"f","", "the name of slow log")
	flag.StringVar(&src,"s","", "the path of slow log")
	flag.StringVar(&excelName,"e","", "the name of new general excel")

	flag.Parse()
	currentTimeBegin := time.Now().Format("2006-01-02T")+"00"
	currentTimeEnd := time.Now().Format("2006-01-02T15:04")
	cmdinfo := fmt.Sprintf("sed -n '/%s/,/%s/p' %s > %s", currentTimeBegin,currentTimeEnd, src+slowLog, currentTimeBegin+".log")
	log.Println("begin to execute:",cmdinfo)
	_, err := ExecCmd(cmdinfo)
	if err != nil{
		loginfo := fmt.Sprintf("can not run command:%s, err:%s", cmdinfo, err.Error())
		log.Println(loginfo)
		return
	}
	cmdinfo = "mysqldumpslow -s t " + "./" + currentTimeBegin+ ".log"
	log.Println("begin to execute:",cmdinfo)
	results, err := ExecCmd(cmdinfo)
	if err != nil {
		loginfo := fmt.Sprintf("can not run command:%s, err:%s", cmdinfo, err.Error())
		log.Println(loginfo)
		return
	}
	resultArrs := strings.Split(results,"Count:")
	var infos = []SlowLogInfo{}
	for index,result := range resultArrs{
		if index == 0 {
			continue
		}
		slowLoginfo := getSlowLogInfo(result)
		infos = append(infos,slowLoginfo)
	}

	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	header := map[string]string{"A1":"Count","B1":"Time","C1":"Lock","D1":"Rows","E1":"SQL"}
	for k,v := range header {
		f.SetCellValue("Sheet1",k,v)
	}


	for i, info := range infos {
		f.SetCellValue("Sheet1","A"+strconv.Itoa(i+2),info.Count)
		f.SetCellValue("Sheet1","B"+strconv.Itoa(i+2),info.Time)
		f.SetCellValue("Sheet1","C"+strconv.Itoa(i+2),info.Lock)
		f.SetCellValue("Sheet1","D"+strconv.Itoa(i+2),info.Rows)
		f.SetCellValue("Sheet1","E"+strconv.Itoa(i+2),info.SQL)

	}
	err = f.SaveAs("./"+excelName)
	if err != nil {
		log.Println("save excel file err:",err)
		return
	}
	log.Println("slowlog2excelold succeed")
	cmdinfo = "rm -rf "+ "./" + currentTimeBegin+ ".log"
	log.Println("begin to execute:",cmdinfo)
	_, err = ExecCmd(cmdinfo)
	if err != nil {
		log.Println("can not delete log file")
		return
	}
	log.Println("delete log file succeed, all done...")
}

//get the details from the arr
func getSlowLogInfo(line string) SlowLogInfo {
	line = strings.Replace(line,"\n","",-1)
	data := strings.Split(line," ")
	slowloginfo := SlowLogInfo{
		Count: data[1],
		Time:  strings.Split(data[3],"=")[1],
		Lock:  strings.Split(data[6],"=")[1],
		Rows:  strings.Split(data[9],"=")[1],
		SQL:   strings.Join(data[13:]," "),
	}
	return slowloginfo
}

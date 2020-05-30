//package Crawl
package main

import (
	"fmt"
	"time"
	"os"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"encoding/csv"
	//"reflect"
)

func CsvInit() (*os.File,*csv.Writer){
	file,err := os.OpenFile("test.csv",os.O_CREATE | os.O_RDWR,0644)
	if err != nil{
		panic(err)
	}
	file.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(file)
	writer.Write([]string{"小区名","板块","户型","面积","年份","总价（万元）","单价（元/平米）"})
	writer.Flush()
	return file,writer
}

func GetInfoInpage(doc *goquery.Document) []map[string]string{
	HouseInfoRet := make([]map[string]string,40)
	content := doc.Find("div[class=\"content \"]")
	content.Find("div[class=\"info clear\"]").Each(func(i int, s *goquery.Selection){
		HouseInfo := make(map[string]string)
		s.Find("div[class=\"positionInfo\"]").Find("a").Each(func(i int, s *goquery.Selection){
			if i == 0{
				HouseInfo["小区"] = s.Text()
			} else if i == 1{
				HouseInfo["板块"] = s.Text()
			}
		})

		detail := s.Find("div[class=\"houseInfo\"]").Text()
		str := strings.Replace(detail," ","",-1)
		rts := strings.Split(str,"|")
		for _,rt := range rts{
			r := regexp.MustCompile("([0-9]+)室([0-9]+)厅")
			matchs := r.FindStringSubmatch(rt)
			if len(matchs) > 0{
				HouseInfo["户型"] = rt
				continue
			}
	
			r = regexp.MustCompile("(\\d+(\\.\\d+)?)平米")
			matchs = r.FindStringSubmatch(rt)
			if len(matchs) >= 2{
				HouseInfo["面积"] = matchs[1]
				continue
			}
	
			r = regexp.MustCompile("([0-9]+)年建")
			matchs = r.FindStringSubmatch(rt)
			if len(matchs) > 0{
				HouseInfo["年份"] = matchs[1]
				continue
			}
		}
		totalPrice := s.Find("div[class=\"totalPrice\"]").Text()
		r := regexp.MustCompile("(\\d+(\\.\\d+)?)万")
		matchs := r.FindStringSubmatch(totalPrice)
		if len(matchs) >= 2 {
			HouseInfo["总价（万元）"] = matchs[1]
		}
		unitPrice := s.Find("div[class=\"unitPrice\"]").Text()
		r = regexp.MustCompile("单价(\\d+(\\.\\d+)?)元/平米")
		matchs = r.FindStringSubmatch(unitPrice)
		if len(matchs) == 2 {
			HouseInfo["单价（元/平米）"] = matchs[1]
		}

		HouseInfoRet = append(HouseInfoRet,HouseInfo)
		
	})

	return HouseInfoRet
}

func GetInfo(start int,end int,writer *csv.Writer){
	for page := start ; page <= end; page ++{
		doc, err := goquery.NewDocument(fmt.Sprintf("https://sh.lianjia.com/ershoufang/pg%d",page))
		if err != nil{
			panic(err)
		}
		house_info := GetInfoInpage(doc)
		time.Sleep(1 * time.Second)
		cnt := 0
		for _,info := range house_info{
			if len(info) != 0 {
				writer.Write([]string{info["小区"],info["板块"],info["户型"],info["面积"],
					info["年份"],info["总价（万元）"],info["单价（元/平米）"]})
				writer.Flush()
				cnt++
			}
		}
		if cnt == 0{
			fmt.Println("not got any info in page ",page)
			break
		}
	}
}

func main(){
	t := time.Now()
	file,writer := CsvInit()
	defer file.Close()
	GetInfo(1,100,writer)
	fmt.Println("TimeCost:",time.Since(t))
}


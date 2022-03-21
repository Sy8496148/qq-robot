package spider

import (
	"encoding/json"
	"github.com/wxnacy/wgo/arrays"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func GetChinaTotal() (*ChinaTotal, error) {
	r, err := download(chinaTotalUri)
	if err != nil {
		log.Println("download error", err)
		return nil, err
	}

	data := r["data"]
	detail := make(map[string]interface{})
	err = json.Unmarshal([]byte(data.(string)), &detail)
	if err != nil {
		log.Println("load json error", err)
		return nil, err
	}

	var ch ChinaTotal

	typeInfo := reflect.TypeOf(ch)
	valInfo := reflect.ValueOf(&ch)

	num := typeInfo.NumField()
	for i := 0; i < num; i++ {
		key := typeInfo.Field(i).Tag.Get("json")
		if i == 0 {
			valInfo.Elem().Field(i).SetString(detail[key].(string))
		} else {
			valInfo.Elem().Field(i).SetInt(int64(detail["chinaTotal"].(map[string]interface{})[key].(float64)))
		}
	}
	return &ch, nil
}

func GetProvinceTotal(province string) (*ProvinceCityTotal, error) {
	r, err := download(provinceCityTotalUri)
	if err != nil {
		log.Println("download error", err)
		return nil, err
	}

	data := r["data"].(map[string]interface{})["diseaseh5Shelf"].(map[string]interface{})["areaTree"].([]interface{})[0].(map[string]interface{})["children"].([]interface{})

	var pc ProvinceCityTotal

	for _, v := range data {
		detail := v.(map[string]interface{})
		if detail["name"] == province {
			return pares(pc, detail), nil
		}
	}

	pc.Exists = false
	return &pc, nil
}

func GetCityTotal(city string) (*ProvinceCityTotal, error) {
	r, err := download(provinceCityTotalUri)
	if err != nil {
		log.Println("download error", err)
		return nil, err
	}

	data := r["data"].(map[string]interface{})["diseaseh5Shelf"].(map[string]interface{})["areaTree"].([]interface{})[0].(map[string]interface{})["children"].([]interface{})

	var pc ProvinceCityTotal
	//判断是否是直辖市和特别行政区
	citys := []string{"香港", "台湾", "上海", "北京", "重庆", "天津", "澳门"}
	index := arrays.ContainsString(citys, city)
	if index == -1 {
		for _, v := range data {
			detail := v.(map[string]interface{})
			for _, c := range detail["children"].([]interface{}) {
				cc := c.(map[string]interface{})
				if cc["name"] == city {
					return pares(pc, cc), nil
				}
			}
		}
	} else {
		for _, v := range data {
			detail := v.(map[string]interface{})
			if detail["name"] == city {
				return pares(pc, detail), nil
			}
		}
	}

	pc.Exists = false
	return &pc, nil
}

func GetProvinceCityTotal(province string, city string) (*ProvinceCityTotal, error) {
	r, err := download(provinceCityTotalUri)
	if err != nil {
		log.Println("download error", err)
		return nil, err
	}

	data := r["data"].(map[string]interface{})["diseaseh5Shelf"].(map[string]interface{})["areaTree"].([]interface{})[0].(map[string]interface{})["children"].([]interface{})

	var pc ProvinceCityTotal

	for _, v := range data {
		detail := v.(map[string]interface{})
		if detail["name"] == province {
			for _, c := range detail["children"].([]interface{}) {
				cc := c.(map[string]interface{})
				if cc["name"] == city {
					return pares(pc, cc), nil
				}
			}
		}
	}

	pc.Exists = false
	return &pc, nil
}

func pares(pc ProvinceCityTotal, cc map[string]interface{}) *ProvinceCityTotal {
	typeInfo := reflect.TypeOf(pc)
	valInfo := reflect.ValueOf(&pc)
	num := typeInfo.NumField()
	for i := 0; i < num-1; i++ {
		key := typeInfo.Field(i).Tag.Get("json")
		valInfo.Elem().Field(i).SetInt(int64(cc["total"].(map[string]interface{})[key].(float64)))
	}
	pc.Exists = true
	return &pc
}

func download(url string) (map[string]interface{}, error) {
	var result map[string]interface{}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	response, err := client.Do(req)
	if err != nil {
		log.Println("http get error", err)
		return nil, err
	}
	//函数结束后关闭相关链接
	defer func() {
		if e := response.Body.Close(); e != nil {
			log.Println("close response error", e, response.Request.URL.String())
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("read error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("load json error", err)
		return nil, err
	}

	return result, nil
}

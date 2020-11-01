package service

import (
	"claps-test/model"
	"claps-test/util"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

const CONTENTTYPE = "application/json"

/**
 * @Description: options选项
 * @return nc
 */
const (
	DEVVAL       = "dev_value"  //开发价值
	COMMIT_NUM   = "commit_num" //commit number
	CHANGE_LINES = "loc"        //change lines
)

type OneMember struct {
	PrimaryEmail string          `json:"primary_email"`
	Value        decimal.Decimal `json:"value"`
}

type mericoReciveData struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    []OneMember `json:"data"`
}

/**
 * @Description:签名类
 */
type SignUtils struct {
	params map[string]interface{}
	key    string
}

/**
 * @Description: 选项类
 */
type Options struct {
	SelectColumn       string `json:"selectColumn"`
	TargetTimezoneName string `json:"targetTimezoneName"`
	SelectProjectId    string `json:"selectProjectId,omitempty"`
	SelectGroupId      string `json:"selectGroupId,omitempty"`
}

/**
 * @Description: 实例化签名
 * @return sign
 */
func newSign() (sign *SignUtils) {
	sign = new(SignUtils)
	sign.key = util.MericoSecret
	sign.params = make(map[string]interface{})
	sign.params["appid"] = util.MericoAppid
	return sign
}

/**
 * @Description: 设置nonStr
 * @receiver s
 * @param nonceStr nonStr为Merico签名必填字段
 */
func (s *SignUtils) setNonceStr(nonceStr string) {
	s.params["nonce_str"] = nonceStr
}

/**
 * @Description: 对map的key排序，返回有序的slice
 * @receiver s
 * @return keys
 */
func (s *SignUtils) sortMapbyKey() (keys []string) {
	for k := range s.params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

/**
 * @Description: 设置k-v键值对
 * @receiver s
 * @param key
 * @param value
 */
func (s *SignUtils) set(key, value string) {
	s.params[key] = value
}

/**
 * @Description: 设置数组,对象类型的k-v
 * @receiver s
 * @param key
 * @param value
 * @return err
 */
func (s *SignUtils) setObjectOrArray(key string, value interface{}) (err error) {
	s.params[key] = value
	return
}

/**
 * @Description: 生成sign
 * @receiver s
 * @return result
 */
func (s *SignUtils) sign() (result string) {
	//对key排序
	keys := s.sortMapbyKey()
	//拼接
	for _, val := range keys {
		//对于array和object需要先转义再拼接
		v, err := s.params[val].(string)
		//断言失败ok为false
		if !err {
			//序列化
			if b, err2 := json.Marshal(s.params[val]); err2 == nil {
				v = string(b)
			} else {
				fmt.Println(err2)
				return
			}
		}
		result += val + "=" + v + "&"
	}
	result += "key=" + util.MericoSecret
	//md5加密
	result = fmt.Sprintf("%x", md5.Sum([]byte(result)))
	//转化大写
	result = strings.ToUpper(result)
	return
}

/**
 * @Description: 获取需要post的数据
 * @receiver s
 * @return result
 * @return err
 */
func (s *SignUtils) getPostData() (result *strings.Reader, err error) {
	//获取sign值
	s.params["sign"] = s.sign()
	//序列化 两次序列化导致转义
	b, err := json.Marshal(s.params)
	if err != nil {
		return
	}
	result = strings.NewReader(string(b))
	return
}

func GetMetricByGroupIdAndUserEmails(groupId string, metric string, members []model.User) (primaryEmailStrs []OneMember, err error) {
	recv, err := getMetric(groupId, metric, members)
	if err != nil {
		return
	}

	//may exist sum of value is not 1 or value is integer,handle this situation
	var sum decimal.Decimal
	for _, v := range recv.Data {
		sum = sum.Add(v.Value)
	}

	for _, v := range recv.Data {
		v.Value = v.Value.Div(sum)
		log.Info("success get devValue:%v", v)
	}
	primaryEmailStrs = recv.Data

	return
}

func getMetric(groupId string, selectColumn string, members []model.User) (recv *mericoReciveData, err error) {
	signTool := newSign()
	signTool.setNonceStr(util.RandUp(16))

	//Set options
	options := Options{
		SelectColumn:       selectColumn,
		TargetTimezoneName: "UTC",
		SelectGroupId:      groupId,
	}

	err = signTool.setObjectOrArray("options", options)
	if err != nil {
		log.Error("Sign set options error: %v", err)
		return
	}

	//get All the email
	var primaryEmails []string
	//traverse emails and append it to primaryEmails
	for _, v := range members {
		primaryEmails = append(primaryEmails, v.Email)
	}

	err = signTool.setObjectOrArray("primaryEmailStrs", primaryEmails)
	if err != nil {
		log.Error("Sign set objectOr Array error. ", err)
		return
	}

	res, err := signTool.getPostData()
	if err != nil {
		log.Error("Sign get post data error. ", err)
		return
	}

	url := util.Merico + "/openapi/openapi/developer/query-efficiency-metric"
	resp, err := http.Post(url, CONTENTTYPE, res)
	if err != nil {
		log.Error("Sign post data error. ", err)
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	recv = &mericoReciveData{}
	err = json.Unmarshal(b, recv)
	if err != nil {
		log.Error("Unmarshal merico data error:%v.", err)
		return
	}

	if recv.Code != 200 {
		log.Error("Merico return code:%v", recv.Code)
		return
	}

	return
}

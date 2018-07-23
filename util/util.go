package util

import (
	"context"
	"fmt"
	//qyhttp "gcoresys/common/http"
	demoHttp "go_demo/http"
	"math"
	"math/rand"
	"net/http"
	"time"

	"bytes"
	"database/sql"
	"errors"
	"flag"
	"gcoresys/common"
	"gcoresys/common/logger"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/json-iterator/go"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/resty.v0"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

// 应用的配置
type AppConfig struct {
	RedisUrl   string
	MysqlHost  string
	MysqlPort  string
	MysqlUname string
	// http服务的端口（在docker中都是3000）
	HttpServerPort string
	MongoHost      string
	MongoUname     string
}

var appConfig *AppConfig

func GetAppConfig() *AppConfig {
	if appConfig != nil {
		return appConfig
	}
	useDocker := GetUseDocker()
	if useDocker == 1 {
		fmt.Println("采用的是devops容器的配置")
		appConfig = GetDevDockerConf()
	} else if useDocker == 2 {
		fmt.Println("采用的是prod生产的配置")
		appConfig = GetProdDockerConf()
	} else {
		fmt.Println("采用的是非docker环境的配置")
		appConfig = GetAppDefaultConf()
	}
	return appConfig
}

// 获取use docker的配置，避免每次打包修改use docker配置
// 表明是否使用容器配置 0是本地开发   1 是 dev  2是 prod
func GetUseDocker() int {
	f := flag.Lookup("docker_env")
	if f == nil || f.Value.String() == "0" {
		// 非docker
		return 0
	} else if f.Value.String() == "2" {
		// 生产
		return 2
	} else {
		// 默认返回开发和测试的配置
		return 1
	}
}

// 获取正常环境下的配置
func GetAppDefaultConf() *AppConfig {
	return &AppConfig{
		RedisUrl:   "127.0.0.1:6379",
		MysqlHost:  "localhost",
		MysqlPort:  "3306",
		MysqlUname: "root",
		MongoHost:  "localhost:27017",
		MongoUname: "admin",
	}
}

// 获取容器中的配置
func GetDevDockerConf() *AppConfig {
	return &AppConfig{
		RedisUrl:  "redis-master:6379",
		MysqlHost: "qy-mysql",
		//MysqlHost: "172.16.0.101",
		MysqlPort:      "3306",
		MysqlUname:     "qy",
		HttpServerPort: "3000",
		MongoHost:      "",
		MongoUname:     "",
	}
}

// 获取容器中的配置
func GetProdDockerConf() *AppConfig {
	uname := "qy"
	if data, err := ioutil.ReadFile("/usr/local/.db/mysql.uname"); err != nil {
		fmt.Println("读取mysql用户名文件出错:" + err.Error() + "。使用默认用户名。")
	} else {
		uname = strings.TrimSpace(string(data))
	}
	return &AppConfig{
		RedisUrl:  "redis-master:6379",
		MysqlHost: "qy-mysql",
		//MysqlHost: "172.16.1.90",
		MysqlPort:      "3306",
		MysqlUname:     uname,
		HttpServerPort: "3000",
		MongoHost:      "",
		MongoUname:     "",
	}
}




// 计算两个时间相差天数
func CountDaysOfTimes(start *time.Time, end *time.Time) int {
	if start == nil || end == nil {
		logger.Warn("start或者end为空")
		return 0
	}
	return SiSheWuRuToInt(GetDate(*end).Sub(GetDate(*start)).Hours() / 24)
}

// 获取一年的天数
func GetYearDays(year int) int {
	if year == 0 {
		return 0
	}
	if isLeapYear(year) {
		return 366
	}
	return 365
}

//判断是否为闰年
func isLeapYear(year int) bool { //y == 2000, 2004
	//判断是否为闰年
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}

// repaychan转中文
func TransRepayChan(rc string) string {
	switch rc {
	case "yb":
		return "易宝"
	case "yfb":
		return "苏宁"
	default:
		return "苏宁"
	}
}

// 获取当前纳秒
func GetCurNanoStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 获取msg id，当前纳秒后加6位随机数
func GetMsgId() string {
	curNano := time.Now().UnixNano()
	r := rand.New(rand.NewSource(curNano))
	return fmt.Sprintf("%d%06v", curNano, r.Int31n(1000000))
}

// 字符串转json对象
func ParseJson(str string, result interface{}) error {
	//return json.Unmarshal([]byte(str), &result)
	return jsoniter.Unmarshal([]byte(str), result)
}

func ParseJsonFromBytes(b []byte, result interface{}) error {
	//return json.Unmarshal(b, &result)
	return jsoniter.Unmarshal(b, result)
}

// json对象转字符串
func StringifyJson(obj interface{}) string {
	//b, err := json.Marshal(obj)
	b, err := jsoniter.Marshal(obj)
	if err != nil {
		fmt.Println("转换json字符串出错")
		return ""
	}
	return string(b)
}

func StringifyJsonToBytes(obj interface{}) []byte {
	//b, err := json.Marshal(obj)
	b, err := jsoniter.Marshal(obj)
	if err != nil {
		fmt.Println("转换json字符串出错")
		return nil
	}
	return b
}

// 获取offset
func GetOffset(page int, perPage int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * perPage
}

//API Gateway Decode
func DecodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r.Body, nil
}

// 包装返回
func EncodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//return json.NewEncoder(w).Encode(response)
	return jsoniter.NewEncoder(w).Encode(response)
}

// 获取不含时分秒的时间
func GetDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

// 在现有时间上加一个月。
func AddOneMonth(t time.Time) time.Time {
	return AddMonthes(t, 1)
}

// 在现有时间上加上n个月
func AddMonthes(t time.Time, times time.Month) time.Time {
	// 在一个月的第一天对month+1是一定正确的
	tBeginM := now.New(t).BeginningOfMonth()
	m := tBeginM.Month()
	// 这里月数大于12小于1会自动换成年
	nextBeginM := time.Date(tBeginM.Year(), m+times, tBeginM.Day(), 0, 0, 0, 0, time.Local)
	nextEndM := now.New(nextBeginM).EndOfMonth()
	if t.Day() > nextEndM.Day() {
		return time.Date(nextEndM.Year(), nextEndM.Month(), nextEndM.Day(), 0, 0, 0, 0, time.Local)
	} else {
		return time.Date(nextEndM.Year(), nextEndM.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	}
}

func CeilFloat64Toint(f float64) int {
	return int(math.Ceil(f))
}

func CeilFloat64ToUint64(f float64) uint64 {
	return uint64(math.Ceil(f))
}

func FloorFloat64ToUint64(f float64) uint64 {
	return uint64(math.Floor(f))
}

// 对传进来的float64做四舍五入，第二个参数是保留小数点后几位
func SiSheWuRu(f float64, remain int) float64 {
	expand := math.Pow10(remain)
	// +0.5是为了执行floor时实现四舍五入
	tmp := f*expand + 0.5
	tmp = math.Floor(tmp)
	return tmp / expand
}

// 四舍五入去掉小数位
func SiSheWuRuCutDecimal(f float64) float64 {
	return SiSheWuRu(f, 0)
}

// 算本金或利息时需要根据配置保留小数
func SiSheWuRuForMoney(f float64) float64 {
	return SiSheWuRu(f, common.KeepFigures)
}

func GetParseJsonErrResp() *demoHttp.BaseResp {
	return &demoHttp.BaseResp{
		Success: false,
		Info:    "请求解析出错",
	}
}

// 从gin中读出请求体的string
func GetStringBodyFromGin(c *gin.Context, keepBody bool) []byte {
	if c.ContentType() == "multipart/form-data" {
		return nil
	}
	params, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	if keepBody {
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(params))
	}
	return params
}

// 获取基本结果
func GetBaseResp(err error, successInfo string) *demoHttp.BaseResp {
	// 如果错误不为空则报错
	if err != nil {
		return &demoHttp.BaseResp{
			Success:  false,
			Info:     err.Error(),
			ErrorMsg: err.Error(),
		}
	}
	return &demoHttp.BaseResp{
		Success: true,
		Info:    successInfo,
	}
}

// 获取通用的成功结果
func GetSuccessBaseResp(info string) *demoHttp.BaseResp {
	return &demoHttp.BaseResp{
		Success: true,
		Info:    info,
	}
}

// 获取通用的错误结果
func GetErrorBaseResp(err string) *demoHttp.BaseResp {
	return &demoHttp.BaseResp{
		Success:  false,
		Info:     err,
		ErrorMsg: err,
	}
}

// 四舍五入保留0位小数并返回int
func SiSheWuRuToInt(f float64) int {
	return int(SiSheWuRu(f, 0))
}

// 获取总页数和总数
func GetTotalPagesAndCount(db *gorm.DB, m interface{}, perPage int) (totalPages int, totalCount int) {
	if perPage <= 0 {
		logger.Warn("获取总页数和总数时perPage不对", "perPage", perPage)
		return
	}
	var count int
	db = db.Offset(-1).Limit(-1)
	db.Model(m).Count(&count)
	totalCount = count
	totalPages = count / perPage
	if count%perPage != 0 {
		totalPages += 1
	}
	return
}

// 获取分页页数总数及数据列表。m是查询的表的model，result是列表结果传指针进来！
func GetDataByPageAndPerPage(db *gorm.DB, page int, perPage int, m interface{}, result interface{}) (totalPages int, totalCount int) {
	offset := GetOffset(page, perPage)
	if err := db.Offset(offset).Limit(perPage).Find(result).Error; err != nil {
		return 0, 0
	}
	totalPages, totalCount = GetTotalPagesAndCount(db, m, perPage)
	return
}

// 返回gin结果，成功失败都可以用这个（但传进来的对象一定要有Success字段）
func RenderGinResult(c *gin.Context, result interface{}) {
	rv := reflect.ValueOf(result)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	sField := rv.FieldByName("Success")
	// 不为空且为true，则说明是成功, bool类型的field不能调IsNil,会炸
	if sField.Kind() == reflect.Bool && sField.Bool() {
		c.JSON(200, result)
	} else {
		c.JSON(400, result)
	}
}

// -------------------- render gin err ---------------------
func RenderGinError(info string, c *gin.Context) {
	c.JSON(400, gin.H{"success": false, "info": info, "errmsg": info})
}

func RenderGinErrorWithCode(statusCode int, info string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"success": false, "info": info, "errmsg": info})
}

func RenderGinErrorJson(info string, data *gin.H, c *gin.Context) {
	result := gin.H{"success": false, "info": info, "errmsg": info}
	if data != nil {
		for k, v := range *data {
			result[k] = v
		}
	}
	c.JSON(400, result)
}

func RenderGinErrorJsonWithCode(statusCode int, info string, data *gin.H, c *gin.Context) {
	result := gin.H{"success": false, "info": info, "errmsg": info}
	if data != nil {
		for k, v := range *data {
			result[k] = v
		}
	}
	c.JSON(statusCode, result)
}

func GinRenderError(c *gin.Context, info string) {
	c.JSON(400, gin.H{"success": false, "info": info, "errmsg": info})
}

func GinRenderAuthError(c *gin.Context, info string) {
	c.JSON(403, gin.H{"success": false, "info": info, "errmsg": info})
}

// -------------------- render gin err ---------------------

func RenderGinSuccess(info string, c *gin.Context) {
	c.JSON(200, gin.H{"success": true, "info": info})
}

// 如果用baseresp则可以不用以下三个方法
func RenderGinSuccessJson(info string, data *gin.H, c *gin.Context) {
	result := gin.H{"success": true, "info": info}
	if data != nil {
		for k, v := range *data {
			result[k] = v
		}
	}
	c.JSON(200, result)
}

// 可以传入任何obj，而不是字符串
func GinRenderJsonObjSuccess(c *gin.Context, resultJson interface{}) {
	c.JSON(200, resultJson)
}

// 返回判断过success的结果
func GinRenderJudgedSuccess(c *gin.Context, resultJson interface{}) {
	var result map[string]interface{}
	if err := ParseJson(resultJson.(string), &result); err != nil {
		logger.Error("解析json报错", "err", err.Error())
	}
	// 判断如果请求返回中成功为false，则返回400
	success := reflect.ValueOf(result["success"])
	if success.Kind() == reflect.Bool && !success.Bool() {
		c.JSON(400, result)
		return
	}
	c.JSON(200, result)
}

// api gw use
func GinRenderSuccess(c *gin.Context, resultJson interface{}) {
	var result map[string]interface{}
	ParseJson(resultJson.(string), &result)
	c.JSON(200, result)
}

// 判断是否是测试环境
func IsTestEnv() bool {
	return flag.Lookup("test.v") != nil
}

// 代理时需要把get的参数割到目标url上
func CutParamsInUrlToTargetUrl(fromUrl string, toUrl string) string {
	splitResult := strings.Split(fromUrl, "?")
	if len(splitResult) > 1 {
		toUrl += "?" + splitResult[1]
	}
	return toUrl
}

// proxy to other server
func ProxyReq(originReq *http.Request, targetUrl string) ([]byte, error) {
	targetUrl = CutParamsInUrlToTargetUrl(originReq.RequestURI, targetUrl)
	logger.Info("代理请求到：" + targetUrl)
	req, _ := http.NewRequest(originReq.Method, targetUrl, originReq.Body)
	return HttpReq(req)
}

// http请求
func HttpReq(r *http.Request) ([]byte, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}
	return b, nil
}

// 开启事务处必须defer调用该函数，否则可能在程序报错后事务既没有commit也没有rollback
func ClearTransaction(tx *gorm.DB) {
	err := tx.Rollback().Error
	if err != sql.ErrTxDone && err != nil {
		logger.Error("关闭事务时出错", "err", err.Error())
	}
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

// 发起JSON请求
func JsonPost(reqUrl string, params interface{}) (respStatus int, respB []byte, err error) {
	logger.Info("发送Json请求", "url", reqUrl, "params", StringifyJson(params))
	if resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(StringifyJsonToBytes(params)).
		Post(reqUrl); err != nil {
		return 0, nil, err
	} else {
		return resp.StatusCode(), resp.Body(), nil
	}
}

// 通用的json请求，把请求当方法用
func GeneralJsonPost(reqUrl string, params interface{}, result interface{}) (err error) {
	respStatus, respB, rErr := JsonPost(reqUrl, params)
	if respStatus != 200 {
		if rErr != nil {
			logger.Error("请求报错", "err", rErr.Error())
			return rErr
		}
		errStr := string(respB)
		logger.Warn("请求返回不是200", "resp", errStr)
		// 远端返回报错也要解析出去
		if pErr := ParseJsonFromBytes(respB, result); pErr != nil {
			logger.Warn("解析返回请求报错", "err", pErr.Error())
			return errors.New("解析返回请求报错:" + pErr.Error())
		}
		return nil
	}
	if err = ParseJsonFromBytes(respB, result); err != nil {
		logger.Error("无法解析请求返回", "err", err.Error())
		return errors.New("无法解析请求返回:" + err.Error())
	}
	return
}

// get请求
func HttpGet(reqUrl string) (respStatus int, respB []byte, err error) {
	if resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Get(reqUrl); err != nil {
		return 0, nil, err
	} else {
		//fmt.Println(string(resp.Body()))
		return resp.StatusCode(), resp.Body(), nil
	}
}

// get请求
func GeneralHttpJsonGet(reqUrl string, result interface{}) (err error) {
	respS, respB, rErr := HttpGet(reqUrl)
	if rErr != nil {
		logger.Warn("http请求报错", "err", rErr.Error())
		return rErr
	} else if respS != http.StatusOK {
		logger.Warn("服务器返回状态不是200", "status", respS, "resp", string(respB))
		return errors.New("服务器返回状态不是200")
	}
	if true {
		logger.Info("请求返回:" + string(respB))
	} else {
		logger.Info("请求返回过长，不打印")
	}
	if err = ParseJsonFromBytes(respB, result); err != nil {
		logger.Warn("无法解析服务器的返回", "err", err.Error(), "resp", string(respB))
		return
	}
	return
}

// 该方法会默认返回id倒序排的字符串
// reqDesc 0则是不倒序
func GetOrderBySql(reqOrderBy string, reqDesc string) string {
	orderBy := "id"
	isDesc := true
	if reqOrderBy != "" {
		orderBy = reqOrderBy
	}
	if reqDesc == "0" {
		isDesc = false
	}
	if isDesc {
		orderBy += " desc"
	}
	return orderBy
}

// 保存上传的文件（图片压缩也可以在这里做）
func SaveUploadFile(fHeader *multipart.FileHeader, savePath string) (err error) {
	// 先创建文件夹
	fileDir := path.Dir(savePath)
	fmt.Println("文件保存位置:" + fileDir)
	os.MkdirAll(fileDir, 0755)
	// 上传的文件
	tmpF, fErr := fHeader.Open()
	if fErr != nil {
		return fErr
	}
	defer tmpF.Close()
	// 目标文件
	outF, cErr := os.Create(savePath)
	if cErr != nil {
		return cErr
	}
	defer outF.Close()
	// 拷过去
	_, err = io.Copy(outF, tmpF)
	return
}

// 检查两个时间是不是同一天
func IsSameDate(firstDate time.Time, secondDate time.Time) bool {
	return firstDate.Year() == secondDate.Year() && firstDate.Month() == secondDate.Month() && firstDate.Day() == secondDate.Day()
}

// 根据第一个还款日计算当前期是哪期
func CountCurTermAndCurTermRepayDate(firstRepayAt time.Time, loanTerm int, haveAdjustment bool) (curTerm int, repayDate time.Time) {
	// 没调整期则第一次还款日在下个月
	if !haveAdjustment {
		firstRepayAt = AddOneMonth(firstRepayAt)
	}
	nowT := time.Now()
	// 如果当前时间比第一期还款时间小，那么就算第一期
	if nowT.Sub(firstRepayAt) < 0 {
		curTerm = 1
		repayDate = firstRepayAt
	} else {
		if nowT.Day() > firstRepayAt.Day() {
			curTerm = 1
		} else {
			curTerm = 0
		}
		// 求差即可
		curTerm += (nowT.Year()-firstRepayAt.Year())*12 + int(nowT.Month()-firstRepayAt.Month()) + 1
		// 如果没有调整期，但计算出来的期数大于贷款期数，那当前期就是最后一期
		if curTerm > loanTerm && !haveAdjustment {
			curTerm = loanTerm
		} else if curTerm > (loanTerm+1) && haveAdjustment {
			curTerm = loanTerm + 1
		}
		repayDate = AddMonthes(firstRepayAt, time.Month(curTerm-1))
	}
	return
}

// 其它服务器返回的数据可能是unicode，因此需要做一个转换
func UnicodeBytesToUTF8Str(ub []byte) (result string, err error) {
	return UnicodeStrToUTF8Str(string(ub))
}

// 解析含unicode的字符串
func UnicodeStrToUTF8Str(originStr string) (result string, err error) {
	tmpStr := originStr
	for len(tmpStr) > 0 {
		v, _, t, qErr := strconv.UnquoteChar(tmpStr, ' ')
		// 如果报错则不解析，只向前移动一个位置
		if qErr != nil {
			result += tmpStr[0:1]
			tmpStr = tmpStr[1:]
		} else {
			result += string(v)
			tmpStr = t
		}
	}
	return
}

// 下划线转驼峰，首字母大写
func StrToTF1(str string) (result string) {
	tmp := strings.Split(str, "_")
	for _, tmpS := range tmp {
		if len(tmpS) > 0 {
			result += strings.ToUpper(tmpS[0:1]) + tmpS[1:]
		}
	}
	return
}

// 迭代一个对象的所有字段名
func EnumAnObjFieldNames(rv reflect.Type, cb func(f reflect.StructField)) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		tmpType := tmpF.Type
		// 如果是时间就不能迭代了
		if tmpType.Kind() == reflect.Struct && !strings.Contains(tmpType.Name(), "Time") && tmpF.Tag.Get("skip") != "true" {
			EnumAnObjFieldNames(tmpType, cb)
		} else {
			cb(tmpF)
		}

	}
}

// 迭代一个对象的所有字段名(可以返回深度), 深度规则与json转换相同, 深度默认最外层为0, currentDepth 不需要传值
func EnumAnObjFieldNamesWithDepth(rv reflect.Type, cb func(f reflect.StructField, depth int), currentDepth ...int) {
	var depth = 0
	if len(currentDepth) > 0 {
		depth = currentDepth[0]
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		tmpType := tmpF.Type
		// 如果是时间就不能迭代了
		if tmpType.Kind() == reflect.Struct && !strings.Contains(tmpType.Name(), "Time") && tmpF.Tag.Get("skip") != "true" {
			if tmpF.Anonymous {
				EnumAnObjFieldNamesWithDepth(tmpType, cb, depth)
			} else {
				cb(tmpF, depth + 1)
				EnumAnObjFieldNamesWithDepth(tmpType, cb, depth + 1)
			}
		} else {
			cb(tmpF, depth)
		}
	}
}

// 迭代一个对象的所有字段的值    -------- alternated by Zebreay
func EnumAnObjFieldValues(rv reflect.Value, cb func(f reflect.Value)) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		// 如果是时间就不能迭代了
		if tmpF.Kind() == reflect.Struct && !strings.Contains(tmpF.Type().Name(), "Time") && rv.Type().Field(i).Tag.Get("skip") != "true" {
			EnumAnObjFieldValues(tmpF, cb)
		} else {
			cb(tmpF)
		}

	}
}

// 终极迭代大法    -------- alternated by Zebreay
func EnumAnStruct(rv reflect.Value, cb func(f reflect.StructField, v reflect.Value)) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpT := rv.Type().Field(i)
		tmpV := rv.Field(i)
		// 如果是时间就不能迭代了
		if tmpV.Kind() == reflect.Struct && !strings.Contains(tmpV.Type().Name(), "Time") && tmpT.Tag.Get("skip") != "true" {
			EnumAnStruct(tmpV, cb)
		} else if tmpV.Kind() == reflect.Slice {
			var nonStructSlice bool
			for i:= 0; i<tmpV.Len() ; i++ {
				if item := tmpV.Index(i); item.Kind() == reflect.Struct && !strings.Contains(tmpV.Type().Name(), "Time") && tmpT.Tag.Get("skip") != "true" {
					EnumAnStruct(item, cb)
				} else {
					nonStructSlice = true
				}
			}
			if nonStructSlice || tmpV.Len() == 0 {
				cb(tmpT, tmpV)
			}
		} else {
			cb(tmpT, tmpV)
		}
	}
}

// 看一个数组中是否含有某个元素
func StrSliceContains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

// 检查文件是否存在
func FileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// 获取当前秒数的字符串
func CurTimeSecStr() string {
	return time.Now().Format("20060102_150405")
}

func IsNull(str string) bool {
	str = strings.Replace(str, " ", "", -1)
	if str != "" {
		str = strings.ToLower(str)
		if str == "null" {
			str = ""
		}
	}
	return str == ""
}

// 将json转map
func JsonToMap(data interface{}) map[string]interface{} {
	if result, ok := data.(map[string]interface{}); ok {
		return result
	} else {
		return map[string]interface{}{}
	}
}

// 根据key将json map转map[string]interface{}
func GetJsonFromJson(data map[string]interface{}, key string) map[string]interface{} {
	if result, ok := data[key].(map[string]interface{}); ok {
		return result
	} else {
		return map[string]interface{}{}
	}
}

// 根据给定的json map key 将得到的值string
func GetStrFromJson(data map[string]interface{}, key string) string {
	if result, ok := data[key].(string); ok {
		return result
	} else {
		return ""
	}
}

// 根据给定的json map key 将得到的值转interface数组:[]interface{}
func GetArrFromJson(data map[string]interface{}, key string) []interface{} {
	if result, ok := data[key].([]interface{}); ok {
		return result
	} else {
		return []interface{}{}
	}
}

// 根据给定的json map key 将得到的值转为float64
func GetFloatFromJson(data map[string]interface{}, key string) float64 {
	if result, ok := data[key].(float64); ok {
		return result
	}
	return 0
}

// 根据给定的json map key 将得到的值转为int
func GetIntFromJson(data map[string]interface{}, key string) int {
	if result, ok := data[key].(int); ok {
		return result
	}
	return 0
}

// 根据给定的json map key 将得到的值转为bool
func GetBoolFromJson(data map[string]interface{}, key string) bool {
	if result, ok := data[key].(bool); ok {
		return result
	}
	return false
}

// 根据给定的json map key 将得到的值转interface{}
func GetMapToInterface(data map[string]interface{}, key string) interface{} {
	if result, ok := data[key]; ok {
		return result
	}
	return []interface{}{}
}

// string to int
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		logger.Warn("string 转 int err", "string", str, "int", i)
		return 0
	}
	return i
}

// 检查文件路径
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 信仰
func WishNoBug() {
	logger.Info("")
	logger.Info("")
	logger.Info(" ━━━━━━ 神兽出没 ━━━━━━			 ")
	logger.Info(" 　　　┏┓　 ┏┓							   ")
	logger.Info(" 　　┏┛┻━━━┛┻━━┓						 ")
	logger.Info(" 　　┃　　　　   ┃						 ")
	logger.Info(" 　　┃　  ━　    ┃						 ")
	logger.Info(" 　　┃　>　  <   ┃						 ")
	logger.Info(" 　　┃ ...　...  ┃					 ")
	logger.Info(" 　　┃　　⌒　    ┃						 ")
	logger.Info(" 　　┃　　　　　 ┃						 ")
	logger.Info(" 　　┗━┓　　　┏━┛						 ")
	logger.Info(" 　　　　┃　　　┃    神兽保佑, 永无BUG!							                        ")
	logger.Info(" 　　　　 ┃　　  ┃ Code is far away from bug with the animal protecting		")
	logger.Info(" 　　　　┃　　　┗━━━┓					  ")
	logger.Info(" 　　　　┃　　　　　　┣┓				")
	logger.Info(" 　　　　┃　　　　　　┏┛				")
	logger.Info(" 　　　　┗┓┓┏━┳┓┏┛ 						")
	logger.Info(" 　　　　　┃┫┫┃┫┫							")
	logger.Info(" 　　　　　┗┻┛┗┻┛							")
	logger.Info("")
	logger.Info("")
}

// gin解析multipart/form-data，支持文件上传
func MultipartForm(c *gin.Context) (*multipart.Form, error) {
	// MaxMultipartMemory = 8 << 20  // 8 MiB
	err := c.Request.ParseMultipartForm(8 << 20)
	return c.Request.MultipartForm, err
}

// 打印版本号与检测数据库连接
func PrintVerAndCheckDb(version string, db *gorm.DB) {
	if version == "" {
		panic("没有传入版本号")
	}
	if db == nil {
		panic("没有传入数据库连接对象")
	}
	logger.Info("当前系统版本:", "version", version)
	rs, err := db.Raw("show tables;").Rows()
	if err != nil {
		panic("列出数据库所有表表名失败:" + err.Error())
	}
	ts := []string{}
	var tName string
	for rs.Next() {
		if err := rs.Scan(&tName); err != nil {
			panic("遍历表名结果报错:" + err.Error())
		}
		ts = append(ts, tName)
	}
	logger.Info("当前系统表清单:", "tables", ts)
}

// 计算逾期情况（m1，m2，m3...），返回逾期天数和状态
func CountOverdueStatus(repayDate time.Time) (uint, string) {
	curDate := GetDate(time.Now())
	repayDate = GetDate(repayDate)
	// 还款日就等于今天就不算逾期
	dur := curDate.Sub(repayDate)
	if dur <= 0 {
		return 0, ""
	}
	days := uint(dur / (24 * time.Hour))
	switch {
	case days < 30:
		return days, "M1"
	case days < 60:
		return days, "M2"
	case days < 90:
		return days, "M3"
	case days < 120:
		return days, "M4"
	case days < 150:
		return days, "M5"
	case days < 180:
		return days, "M6"
	default:
		return days, "M6+"
	}
}

// 三目表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfFunc(condition bool, trueFunc, falseFunc interface{}, params ...interface{}) []reflect.Value {
	var paramsValue []reflect.Value
	for _, x := range params {
		paramsValue = append(paramsValue, reflect.ValueOf(x))
	}
	if condition {
		return reflect.ValueOf(trueFunc).Call(paramsValue)
	}
	return reflect.ValueOf(falseFunc).Call(paramsValue)
}

// Case When Then
func CaseWhen(whenThen ...interface{}) interface{} {
	for i := 0; i < len(whenThen)-1; i += 2 {
		if whenThen[i].(bool) {
			return whenThen[i+1]
		}
	}
	return GenEmptyValue(whenThen[len(whenThen)-1])
}

// switch case
func SwitchCase(switchValue interface{}, def interface{}, caseValue ...interface{}) interface{} {
	for i := 0; i < len(caseValue)-1; i += 2 {
		if switchValue == caseValue[i] {
			return caseValue[i+1]
		}
	}
	return def
}

// interface 转换为 []interface
func ToSlice(array interface{}) []interface{} {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice {
		panic("ToSlice array not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

// 判断 item 是否在 数组 array 里
func InArray(item interface{}, array interface{}) (bool) {
	arr := ToSlice(array)
	for _, x := range arr {
		if item == x {
			return true
		}
	}
	return false
}

// 判断 item 是否在 数组 array 里, 并且返回首个匹配值的index  ---by Zebreay
func InArrayWithIndex(item interface{}, array interface{}) (bool, int) {
	arr := ToSlice(array)
	for i, x := range arr {
		if item == x {
			return true, i
		}
	}
	return false, -1
}

// 根据给定参数的类型, 生成该类型的空值
func GenEmptyValue(v interface{}) interface{} {
	//return reflect.New(reflect.TypeOf(v)).Elem().Interface()
	return reflect.Zero(reflect.TypeOf(v)).Interface()
}

// 判断一个变量是否为空
func IsEmpty(object interface{}) bool {

	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	for _, v := range numericZeros {
		if object == v {
			return true
		}
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Slice, reflect.Chan:
		{
			return objValue.Len() == 0
		}
	case reflect.Struct:
		switch object.(type) {
		case time.Time:
			return object.(time.Time).IsZero()
		}
	case reflect.Ptr:
		{
			if objValue.IsNil() {
				return true
			}
			switch object.(type) {
			case *time.Time:
				return object.(*time.Time).IsZero()
			default:
				return false
			}
		}
	}
	return false
}
var numericZeros = []interface{}{
	int(0),
	int8(0),
	int16(0),
	int32(0),
	int64(0),
	uint(0),
	uint8(0),
	uint16(0),
	uint32(0),
	uint64(0),
	float32(0),
	float64(0),
}
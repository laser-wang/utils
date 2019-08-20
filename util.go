// 其他项目也可使用的工具类函数
package utils

import (
	"bufio"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	//	"math"
	mRand "math/rand"
	"net"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dchest/pbkdf2"

	clog "github.com/cihub/seelog"
	"github.com/denisbrodbeck/machineid"
	"github.com/jaypipes/ghw"
	"github.com/satori/go.uuid"
	"github.com/ssimunic/gosensors"
	"gopkg.in/mgo.v2/bson"
)

// 普通常量
const (
	// checkErr 函数用到的常量 1表示记录日志并退出 2表示仅记录日志，程序不退出
	CHECK_FLAG_EXIT    = 1
	CHECK_FLAG_LOGONLY = 2

	// 分隔符
	SPLIT_CHAR = "_"
	// 注释常亮
	COMMENT_STR = "===================="

	TIME_SECOND_TEN  = 10
	TIME_MINUTE_ONE  = 60
	TIME_MINUTE_FIVE = 5 * 60
	TIME_MINUTE_TEN  = 10 * 60
	TIME_HOUR_ONE    = 60 * 60
	TIME_HOUR_TWO    = 2 * 60 * 60
	TIME_DAY_ONE     = 24 * 60 * 60
	TIME_MAX         = TIME_DAY_ONE * 365 * 100

	TIME_CACHE = TIME_HOUR_TWO

	PathProcCpuinfo = "/proc/cpuinfo"
)

type SizeSymbol struct {
	Size uint64
	Name string
}

var (
	SizeSymbols = []SizeSymbol{
		{uint64(1024 * 1024 * 1024 * 1024), "T"},
		{uint64(1024 * 1024 * 1024), "G"},
		{uint64(1024 * 1024), "M"},
		{uint64(1024), "K"},
	}
)

func GetHumanSize(v uint64) string {
	var s *SizeSymbol
	for _, ss := range SizeSymbols {
		if v >= ss.Size {
			s = &ss
			break
		}
	}
	if s != nil {
		ret := float64(v) / float64(s.Size)
		return fmt.Sprintf("%.1f%s", ret, s.Name)
	}
	return fmt.Sprintf("%dB", v)
}

// CheckErr 错误处理函数，程序中错误分2种，一种需要终止程序，一种仅仅是记录错误日志 flag:  1:exit  2:log only
func CheckErr(err error, flag int) {

	var path string

	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		path = " -- " + file + ":" + strconv.Itoa(line)

		switch flag {
		case CHECK_FLAG_EXIT:
			clog.Critical(err.Error() + path)
			clog.Critical(StackTrace(false))
			panic(err)
		case CHECK_FLAG_LOGONLY:
			clog.Critical(err.Error() + path)
			clog.Critical(StackTrace(false))
		default:
			clog.Info(err.Error() + path)
		}
	}

}

func CheckErrEcho(err error, flag int) {

	var path string

	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		path = " -- " + file + ":" + strconv.Itoa(line)

		switch flag {
		case CHECK_FLAG_EXIT:
			fmt.Println(err.Error() + path)
			fmt.Println(StackTrace(false))
			panic(err)
		case CHECK_FLAG_LOGONLY:
			fmt.Println(err.Error() + path)
			fmt.Println(StackTrace(false))
		default:
			fmt.Println(err.Error() + path)
		}
	}

}

func StackTrace(all bool) string {
	// Reserve 10K buffer at first
	buf := make([]byte, 10240)

	for {
		size := runtime.Stack(buf, all)
		// The size of the buffer may be not enough to hold the stacktrace,
		// so double the buffer size
		if size == len(buf) {
			buf = make([]byte, len(buf)<<1)
			continue
		}
		break
	}
	return string(buf)
}

// FileIsExist 检查文件或目录是否存在,如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func FileIsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func FileDel(filename string) bool {
	err := os.Remove(filename)
	if err == nil {
		return true
	} else {
		clog.Error(err)
		return false
	}
}

func ObjToStr(obj interface{}) string {
	if obj == nil {
		return ""
	}
	b, err := json.Marshal(obj)
	if err != nil {
		clog.Error(err)
		return ""
	}
	return (string(b))
}

func Atoi(str string) int {
	if str == "" {
		return 0
	}
	i, err := strconv.Atoi(str)
	CheckErr(err, CHECK_FLAG_LOGONLY)
	return i
}
func Atoi64(str string) int64 {
	if str == "" {
		return 0
	}
	i, err := strconv.ParseInt(str, 10, 64)
	CheckErr(err, CHECK_FLAG_LOGONLY)
	return i
}
func Atof64(str string) float64 {
	if str == "" {
		return 0.0
	}
	i, err := strconv.ParseFloat(str, 64)
	CheckErr(err, CHECK_FLAG_LOGONLY)
	return i
}

func Itof64(i int) float64 {
	a := Itoa(i)
	f := Atof64(a)
	return f
}

// Ftoa64 转字符串 prec：保留几位小数
func F64toa(f float64, prec int) string {
	str := strconv.FormatFloat(f, 'f', prec, 64)
	return str
}

// Ftoa64 转int，如有小数，则抛弃
func F64toi(f float64) int {
	y := int(f)
	return y
}

func Itoa(i int) string {
	str := strconv.Itoa(i)
	return str
}
func Itoa64(i int64) string {
	str := strconv.FormatInt(i, 10)
	return str
}

// ItoaZero 数字转字符串，前面补零
func ItoaZero(i int, lenth int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(lenth)+"d", i)
}

// ItoaZero64 数字转字符串，前面补零
func ItoaZero64(i int64, lenth int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(lenth)+"d", i)
}
func Btoa(b bool) string {
	str := strconv.FormatBool(b)
	return str
}

// Itob int -> bool
func Itob(i int) bool {
	if i == 1 {
		return true
	} else {
		return false
	}
}

// Itob int -> bool
func Atob(i string) bool {
	if i == "1" {
		return true
	} else {
		return false
	}
}

// GetToken 生成token
func GetToken() string {
	// Creating UUID Version 4
	u1, _ := uuid.NewV4()

	return Sha1(u1.String())
}

func Sha256(str string) string {
	sum := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sum)
}
func Sha1(str string) string {
	sum := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", sum)
}
func Md5(str string) string {
	sum := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", sum)
}
func Pbkdf2(str string) string {
	password := str
	// Get random salt
	salt := make([]byte, 32)
	if _, err := rand.Reader.Read(salt); err != nil {
		panic("random reader failed")
	}
	//12345678901234567890123456789012
	salt = []byte("ae260d66ed648a7ffb6d9286b080ee91")
	// Derive key
	//	key := pbkdf2.WithHMAC(sha256.New, []byte(password), salt, 9999, 64)
	key := pbkdf2.WithHMAC(sha256.New, []byte(password), salt, 9999, 16)
	return fmt.Sprintf("%x", key)
}
func HmacSha1(str string, secretKey string) string {
	key := []byte(secretKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// wrap 给string外面加上""
func Wrap(str string) string {
	return `"` + str + `"`
}

// UnWrap 去掉string外面的双引号（从redis中取出数据时会有） eg. "dddd" -> dddd
func UnWrap(str string) string {
	rs := []rune(str)
	quote := []rune(`"`)
	lenth := len(rs)
	if lenth >= 2 {
		if rs[0] == quote[0] && rs[lenth-1] == quote[0] {
			return string(rs[1 : lenth-1])
		} else if rs[0] == quote[0] {
			return string(rs[1:lenth])
		} else if rs[lenth-1] == quote[0] {
			return string(rs[0 : lenth-1])
		} else {
			return str
		}
	} else if lenth == 1 {
		if rs[0] == quote[0] {
			return ""
		} else {
			return str
		}
	} else {
		return str
	}
}

// 返回两个时间戳之间的间隔，单位：秒  ret = time2 - time1
func TimeDiff(time1 int64, time2 int64) int64 {
	return time2 - time1
}

func GenKeyPair() []string {
	accessKey := Md5(GetToken())
	secretKey := Md5(GetToken())

	ret := []string{accessKey, secretKey}
	return ret
}

// setEnv 设置环境变量
func SetEnv(key string, value string) bool {
	err := os.Setenv(key, value)
	if err != nil {
		return false
	} else {
		return true
	}
}

// GetEnv 取得环境变量
func GetEnv(key string) string {
	ret := os.Getenv(key)
	return ret
}

// ProcWritePID 把进程的pid写入文件,参数是pid路径和文件名。
func ProcWritePID(PIDFileNameWithPath string) {
	iManPid := fmt.Sprint(os.Getpid())
	tmpFileName := PIDFileNameWithPath

	if FileExist(tmpFileName) == true {
		if ret := ProcExist(tmpFileName); ret == 1 {
			return
		}
	}
	pidFile, _ := os.Create(tmpFileName)
	defer pidFile.Close()
	pidFile.WriteString(iManPid)

	return

}

// 判断进程是否启动 1:已启动 0：未启动
func ProcExist(PIDFileNameWithPath string) (ret int) {
	tmpFileName := PIDFileNameWithPath
	if FileExist(tmpFileName) == false {
		return 0
	}
	iManPidFile, err := os.Open(tmpFileName)
	defer iManPidFile.Close()
	if err == nil {
		filePid, err := ioutil.ReadAll(iManPidFile)
		if err == nil {
			pidStr := fmt.Sprintf("%s", filePid)
			pid, _ := strconv.Atoi(pidStr)
			_, err := os.FindProcess(pid)
			if err == nil {
				return 1
			}
		}
	}
	return 0
}

// FileExist 判断文件是否存在 true:存在 false:不存在
func FileExist(fileNameWithPath string) bool {
	_, err := os.Stat(fileNameWithPath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// AddOneDay 日期加一天 yyyy-mm-dd
func AddOneDay(dt string) (ret string) {

	tm, _ := time.Parse(DATE_FORMAT_SHORT, dt)
	ta := tm.AddDate(0, 0, 1)
	ret = ta.Format(DATE_FORMAT_SHORT)
	return ret
}

// AddDays 日期加n天 yyyy-mm-dd
func AddDays(dt string, n int) (ret string) {

	tm, _ := time.Parse(DATE_FORMAT_SHORT, dt)
	ta := tm.AddDate(0, 0, n)
	ret = ta.Format(DATE_FORMAT_SHORT)
	return ret
}

// AddOneMonth 日期加一个月 yyyy-mm-dd
func AddOneMonth(dt string) (ret string) {

	tm, _ := time.Parse(DATE_FORMAT_SHORT, dt)
	ta := tm.AddDate(0, 1, 0)
	ret = ta.Format(DATE_FORMAT_SHORT)
	return ret
}

// AddMonths 日期加n个月 yyyy-mm-dd
func AddMonths(dt string, n int) (ret string) {

	tm, _ := time.Parse(DATE_FORMAT_SHORT, dt)
	ta := tm.AddDate(0, n, 0)
	ret = ta.Format(DATE_FORMAT_SHORT)
	return ret
}

// AddOneHour 时间加一小时 yyyy-mm-dd hh:mm:ss
func AddOneHour(dt string) (ret string) {
	tm, _ := time.Parse(DATE_FORMAT_LONG, dt)
	ta := tm.Add(60 * time.Minute)
	ret = ta.Format(DATE_FORMAT_LONG)
	return ret
}

// AddMinute 时间加若干分钟 yyyy-mm-dd hh:mm:ss
func AddMinute(dt string, cnt int) (ret string) {
	tm, _ := time.Parse(DATE_FORMAT_LONG, dt)
	ta := tm.Add(time.Duration(cnt) * time.Minute)
	ret = ta.Format(DATE_FORMAT_LONG)
	return ret
}

// AddSecond 时间加若干秒钟 yyyy-mm-dd hh:mm:ss
func AddSecond(dt string, cnt int) (ret string) {
	tm, _ := time.Parse(DATE_FORMAT_LONG, dt)
	ta := tm.Add(time.Duration(cnt) * time.Second)
	ret = ta.Format(DATE_FORMAT_LONG)
	return ret
}

// GetTimeStr 本函数返回某一个时间开始，第n个线程的起始时间
// 比如 workDtStart='2017-01-01 00:00:00' ,period =24*60,threadCnt=6 threadId = 2,
// 那么 返回值等于'2017-01-01 08:00:00'，如果threadId=1,那么返回值等于'2017-01-01 04:00:00'
// workDtStart：开始时间(yyyy-mm-dd hh:mm:ss)
// period：时间长度，单位分钟
// threadCnt：线程数
// threadId：线程序号
func GetTimeStr(workDtStart string, period int, threadCnt int, threadId int) (ret string) {
	totalSec := period * 60
	secPerThead := totalSec / threadCnt
	tm, _ := time.Parse(DATE_FORMAT_LONG, workDtStart)
	dur, _ := time.ParseDuration(Itoa(secPerThead*threadId) + "s")
	ta := tm.Add(dur)
	ret = ta.Format(DATE_FORMAT_LONG)
	return ret
}

// Time2StrL long time -> yyyy-mm-dd hh:mm:ss
func Time2StrL(t time.Time) string {
	return t.Format(DATE_FORMAT_LONG)
}

// Time2StrS short time -> yyyy-mm-dd
func Time2StrS(t time.Time) string {
	return t.Format(DATE_FORMAT_SHORT)
}

// TimeStamp2StrL long time -> yyyy-mm-dd hh:mm:ss
func TimeStamp2StrL(ts int64) string {
	t := time.Unix(ts, 0).UTC()
	return t.Format(DATE_FORMAT_LONG)
}

// TimeStamp2StrS short time -> yyyy-mm-dd
func TimeStamp2StrS(ts int64) string {
	t := time.Unix(ts, 0).UTC()
	return t.Format(DATE_FORMAT_SHORT)
}

// Str2TimeL yyyy-mm-dd hh:mm:ss -> time
func Str2TimeL(s string) time.Time {
	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation(DATE_FORMAT_LONG, s, loc)
	CheckErr(err, CHECK_FLAG_LOGONLY)
	return t
}

// Str2TimeS yyyy-mm-dd -> time
func Str2TimeS(s string) time.Time {
	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation(DATE_FORMAT_SHORT, s, loc)
	CheckErr(err, CHECK_FLAG_LOGONLY)
	return t
}

// Str2TimeStampL yyyy-mm-dd hh:mm:ss -> timeStamp
func Str2TimeStampL(s string) int64 {
	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation(DATE_FORMAT_LONG, s, loc)
	CheckErr(err, CHECK_FLAG_LOGONLY)
	return t.Unix()
}

// Str2TimeStampS yyyy-mm-dd -> timeStamp
func Str2TimeStampS(s string) int64 {
	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation(DATE_FORMAT_SHORT, s, loc)
	CheckErr(err, CHECK_FLAG_LOGONLY)
	return t.Unix()
}

func IsDate(s string) bool {
	var err error
	loc, _ := time.LoadLocation("UTC")
	if len(s) == 10 {
		_, err = time.ParseInLocation(DATE_FORMAT_SHORT, s, loc)
	} else {
		_, err = time.ParseInLocation(DATE_FORMAT_LONG, s, loc)
	}

	if err != nil {
		return false
	}
	return true
}

// GetNow 取得当前日期时间
func GetNow(format string) string {
	return time.Now().Format(format)
}
func GetNowUTC(format string) string {
	return time.Now().UTC().Format(format)
}
func GetNowUTC2() string {
	return time.Now().UTC().Format(DATE_FORMAT_LONG)
}
func GetNowUTC2Num() int64 {
	return Str2TimeStampL(GetNowUTC2())
}

// 取得本月最后一天
func GetLastDayOfMonth(date string) string {
	now := Str2TimeL(date)
	loc, _ := time.LoadLocation("UTC")
	currentYear, currentMonth, _ := now.Date()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, loc)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return Time2StrL(lastOfMonth)
}

// 取得本月第一天 yyyy-mm-dd hh:mm:ss
func GetFirstDayOfMonth(date string) string {
	now := Str2TimeL(date)
	loc, _ := time.LoadLocation("UTC")
	currentYear, currentMonth, _ := now.Date()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, loc)
	return Time2StrL(firstOfMonth)
}

// 取得下月第一天 yyyy-mm-dd hh:mm:ss
func GetFirstDayOfNextMonth(date string) string {
	dayNext := AddOneMonth(Substring(date, 0, 10)) + " 00:00:00"
	day := GetFirstDayOfMonth(dayNext)
	return day
}

// Random 取得随机值, 0 <= 返回值 < iMax
func RandomInt(iMax int) int {
	return mRand.Intn(iMax)
}

// RandomInt64 取得随机值, 0 <= 返回值 < iMax
func RandomInt64(iMax int64) int64 {
	return mRand.Int63n(iMax)
}

// GetMongoObjectId 取得mongo objectid，24位
func GetMongoObjectId() string {
	ret := bson.NewObjectId()
	return ret.Hex()
}

// Substring 取得子字符串  substring(yyyy-mm-dd hh:mm:ss,0,13) -> yyyy-mm-dd hh
func Substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

// Args2Str 把不定长参数转成字符串，主要用以打印日志
func Args2Str(args ...interface{}) (ret string) {
	split := ","
	for _, v := range args {
		switch v.(type) {
		case string:
			ret = ret + v.(string) + split
		case int:
			ret = ret + Itoa(v.(int)) + split
		case int64:
			ret = ret + Itoa64(v.(int64)) + split
		default:
			ret = ret + fmt.Sprintf("%T", v) + split
		}
	}
	return
}

// DtDiff 计算两个时间差 ret = timeB - timeA ,参数：yyyy-mm-dd hh:mm:ss 返回值：time.Duration
func DtDiff(timeA string, timeB string) (ret time.Duration) {
	dtA, _ := time.Parse(DATE_FORMAT_LONG, timeA)
	dtB, _ := time.Parse(DATE_FORMAT_LONG, timeB)
	d := dtB.Sub(dtA)
	return d
}

// DtCompare 比较两个时间。
// timeB > timeA 返回1
// timeB == timeA 返回0
// timeB < timeA 返回-1
func DtCompare(timeA string, timeB string) int {
	diff := DtDiff(timeA, timeB)
	if diff.Seconds() > 0 {
		return 1
	} else if diff.Seconds() < 0 {
		return -1
	} else {
		return 0
	}
}

func Len(str string) int {
	return len([]rune(str))
}

func WaitEndSignal() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done
}

// AddSlash   eg1."xxxx" -> "xxxx/"   eg2."xxxx/" -> "xxxx/"
func AddSlash(str string) string {

	if str == "" {
		return ""
	}
	rs := []rune(str)
	quote := []rune(`/`)
	lenth := len(rs)
	if rs[lenth-1] == quote[0] {
		return str
	} else {
		return str + `/`
	}
}

// DelSlash   eg1."xxxx/" -> "xxxx"   eg2."xxxx" -> "xxxx"
func DelSlash(str string) string {

	if str == "" {
		return ""
	}
	rs := []rune(str)
	quote := []rune(`/`)
	lenth := len(rs)
	if rs[lenth-1] == quote[0] {
		return string(rs[0 : lenth-1])
	} else {
		return str
	}
}

// DtPatch 检查日期有效性，如果是 2017-9-1这种形式的日期，则转换为 2017-09-01
func DtPatch(dt string) (ret bool, dtNew string) {
	dtNew = ""
	if dt == "" {
		return false, ""
	}
	if len(dt) > 10 {
		return false, ""
	}
	splitStr := strings.Split(dt, "-")
	if len(splitStr) != 3 {
		return false, ""
	}
	year := Atoi(splitStr[0])
	month := Atoi(splitStr[1])
	day := Atoi(splitStr[2])

	if month > 12 || month < 1 {
		return false, ""
	}

	if day > 31 || day < 1 {
		return false, ""
	}

	dtNew = ItoaZero(year, 4) + "-" + ItoaZero(month, 2) + "-" + ItoaZero(day, 2)
	retCheck := DtCheck(dtNew)
	if retCheck == false {
		return false, ""
	}

	return true, dtNew
}

// DtCheck 判断是否是合法日期 yyyy-mm-dd
func DtCheck(dt string) bool {
	loc, _ := time.LoadLocation("UTC")
	_, err := time.ParseInLocation(DATE_FORMAT_SHORT, dt, loc)
	if err != nil {
		return false
	}
	return true
}

// GetLMSensorData 检测cpu温度等，
// 使用前需要运行 1）yum install lm_sensors  2） sensors-detect
func GetLMSensorData() (ret bool, lmData *LMData) {
	lmData = new(LMData)

	sensors, err := gosensors.NewFromSystem()

	if err != nil {
		return false, nil
	}

	lmData.Data = sensors.JSON()
	lmData.Cpu = make([]*LMCpu, 0)
	lmData.Fan = make([]*LMFan, 0)

	// Iterate over chips
	for chip := range sensors.Chips {
		// Iterate over entries
		for key, value := range sensors.Chips[chip] {
			if strings.Contains(key, "Core") ||
				strings.Contains(key, "CORE") ||
				strings.Contains(key, "core") ||
				strings.Contains(key, "cpu") ||
				strings.Contains(key, "Cpu") ||
				strings.Contains(key, "CPU") {
				idxStart := strings.Index(value, `+`)
				idxEnd := strings.Index(value, `°C`)
				temperature := Substring(value, idxStart, idxEnd)
				lmCpu := new(LMCpu)
				lmCpu.Temperature = int(Atof64(temperature))
				lmData.Cpu = append(lmData.Cpu, lmCpu)
			}
			if strings.Contains(key, "Fan") ||
				strings.Contains(key, "fan") ||
				strings.Contains(key, "FAN") {

				idxStart := strings.Index(value, `:`)
				idxEnd := strings.Index(value, `RPM`)
				rpm := Substring(value, idxStart, idxEnd)

				lmFan := new(LMFan)
				lmFan.Rpm = int(Atof64(rpm))
				lmData.Fan = append(lmData.Fan, lmFan)
			}
		}
	}

	return true, lmData
}

// IsLetterNumber 是否只包含英文字母和数字
func IsLetterNumber(str string) bool {
	match, _ := regexp.MatchString(`^[~!@#%&_<>,/\|\{\}\\\.\+\$\*\(\)\^\?\[\]A-Za-z0-9]+$`, str)
	//`^[~!@#%&_<>,/\|\{\}\\\.\+\$\*\(\)\^\?\[\]A-Za-z0-9]+$`
	//`^(?=.*[0-9].*)(?=.*[A-Z].*)(?=.*[a-z].*).{6,20}`

	return match
}

// HasCapNumLow 至少一个大写字母一个小写字母一个数字
func HasCapNumLow(str string) bool {
	match1, _ := regexp.MatchString(`[0-9]+`, str)
	match2, _ := regexp.MatchString(`[A-Z]+`, str)
	match3, _ := regexp.MatchString(`[a-z]+`, str)

	if match1 == true && match2 == true && match3 == true {
		return true
	} else {
		return false
	}
}

func IsEmail(str string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9\._-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, str)
	//	^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$
	//	^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$
	//	^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$
	return match
}

// IsValidPassword 用于检查是否是合法密码。要求：大于6位，只能包含大小写英文和数字以及特殊字符，至少一个大写字母一个小写字母一个数字
// 返回值：-1:密码小于6位 -2:至少需要包含一个大写字符一个小写字符一个数字 -3：只能包含大小写英文和数字以及特殊字符
func IsValidPassword(str string) int {
	if Len(str) < 6 {
		return -1
	}
	if IsLetterNumber(str) == false {
		return -3
	}
	if HasCapNumLow(str) == false {
		return -2
	}
	return 0
}

// IsIp 检查是否是合法ip地址
func IsIp(ipStr string) bool {
	a := net.ParseIP(ipStr)
	if a == nil {
		return false
	}
	return true
}

// GetIpLocal 取得本机局域网ip
func GetIpLocal() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Panic2Str panic产生的错误，转化为error
func Panic2Error(r interface{}) error {
	var err error

	switch x := r.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = errors.New("Unknow panic")
	}
	return err
}

func NvlStr(value *string) string {
	if value == nil {
		return ""
	} else {
		return *value
	}
}
func NvlInt(value *int) int {
	if value == nil {
		return 0
	} else {
		return *value
	}
}
func NvlInt64(value *int64) int64 {
	if value == nil {
		return 0
	} else {
		return *value
	}
}

// GetMachineId 取得机器码
func GetMachineId() string {
	id, err := machineid.ID()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return id
}

func GetFinger() string {
	hardwareStr := getHardwareStr()
	//	fmt.Println("hardwareStr:\n" + hardwareStr)
	machineId := GetMachineId()
	//	fmt.Println("machineId:\n" + machineId)

	// 主板id
	//sys/class/dmi/id/product_uuid

	fingerOrigin := hardwareStr + machineId
	finger := Md5(fingerOrigin)
	return finger
}

// getHardwareStr 取得硬件字符串
func getHardwareStr() string {

	retHardwareStr := ""

	// --------------------------------------
	//	fmt.Println("--------------mem-------------")
	memHeadStr := fmt.Sprintf("--------------mem-------------\n")
	retHardwareStr = retHardwareStr + memHeadStr
	//	fmt.Println("")

	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}

	memoryStr := fmt.Sprintf("%v\n", memory)
	retHardwareStr = retHardwareStr + memoryStr

	//	fmt.Println(memoryStr)

	//	// cpu	// --------------------------------------
	//	fmt.Println("--------------cpu-------------")
	//	fmt.Println("")

	cpuHeadStr := fmt.Sprintf("--------------cpu-------------\n")
	retHardwareStr = retHardwareStr + cpuHeadStr

	cpuInfo := getCpuInfo()
	procStr := fmt.Sprintf("%s\n", cpuInfo)
	retHardwareStr = retHardwareStr + procStr

	//	fmt.Println(procStr)

	//	cpu, err := ghw.CPU()
	//	if err != nil {
	//		fmt.Printf("Error getting CPU info: %v", err)
	//	}

	//	cpuStr := fmt.Sprintf("%v\n", cpu)
	//	retHardwareStr = retHardwareStr + cpuStr

	//	for _, proc := range cpu.Processors {
	//		//		procStr := fmt.Sprintf("%v\n", proc)
	//		//		retHardwareStr = retHardwareStr + procStr
	//		for _, core := range proc.Cores {
	//			//			coreStr := fmt.Sprintf("%v\n", core)
	//			//			retHardwareStr = retHardwareStr + coreStr
	//		}
	//		if len(proc.Capabilities) > 0 {
	//			// pretty-print the (large) block of capability strings into rows
	//			// of 6 capability strings
	//			rows := int(math.Ceil(float64(len(proc.Capabilities)) / float64(6)))
	//			for row := 1; row < rows; row = row + 1 {
	//				rowStart := (row * 6) - 1
	//				rowEnd := int(math.Min(float64(rowStart+6), float64(len(proc.Capabilities))))
	//				rowElems := proc.Capabilities[rowStart:rowEnd]
	//				capStr := strings.Join(rowElems, " ")
	//				if row == 1 {
	//					fmt.Printf("  capabilities: [%s\n", capStr)
	//					//					capStrA := fmt.Sprintf("  capabilities: [%s\n", capStr)
	//					//					retHardwareStr = retHardwareStr + capStrA
	//				} else if rowEnd < len(proc.Capabilities) {
	//					fmt.Printf("                 %s\n", capStr)
	//					//					capStrB := fmt.Sprintf("                 %s\n", capStr)
	//					//					retHardwareStr = retHardwareStr + capStrB
	//				} else {
	//					fmt.Printf("                 %s]\n", capStr)
	//					//					capStrC := fmt.Sprintf("                 %s\n", capStr)
	//					//					retHardwareStr = retHardwareStr + capStrC
	//				}
	//			}
	//		}
	//	}

	// --------------------------------------
	//	fmt.Println("--------------Block storage-------------")
	//	fmt.Println("")

	storageHeadStr := fmt.Sprintf("--------------Block storage-------------\n")
	retHardwareStr = retHardwareStr + storageHeadStr

	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}

	//	fmt.Printf("%v\n", block)
	blockStr := fmt.Sprintf("%v\n", block)
	retHardwareStr = retHardwareStr + blockStr

	for _, disk := range block.Disks {
		//		a := fmt.Sprintf("%d\n", disk.SizeBytes)
		//		b := fmt.Sprintf("%s\n", disk.SerialNumber)
		//		retHardwareStr = retHardwareStr + a
		//		retHardwareStr = retHardwareStr + b
		if disk.SizeBytes > 0 && disk.SerialNumber != "" {

			//		fmt.Printf(" %v\n", disk)
			diskStr := fmt.Sprintf("%v\n", disk)
			// 忽略光驱
			isODD := strings.Index(diskStr, "ODD")
			if isODD >= 0 {
				continue
			}
			retHardwareStr = retHardwareStr + diskStr

			//			for _, part := range disk.Partitions {
			//			fmt.Printf("  %v\n", part)
			//				partStr := fmt.Sprintf("%v\n", part)
			//				retHardwareStr = retHardwareStr + partStr
			//			}
		}

	}

	// --------------------------------------
	//	fmt.Println("--------------Topology-------------")
	//	fmt.Println("")

	//	topologyHeadStr := fmt.Sprintf("--------------Topology-------------\n")
	//	retHardwareStr = retHardwareStr + topologyHeadStr

	//	topology, err := ghw.Topology()
	//	if err != nil {
	//		fmt.Printf("Error getting topology info: %v", err)
	//	}

	//	//	fmt.Printf("%v\n", topology)
	//	topologyStr := fmt.Sprintf("%v\n", topology)
	//	retHardwareStr = retHardwareStr + topologyStr

	//	for _, node := range topology.Nodes {
	//		//		fmt.Printf(" %v\n", node)
	//		nodeStr := fmt.Sprintf("%v\n", node)
	//		retHardwareStr = retHardwareStr + nodeStr
	//		for _, cache := range node.Caches {
	//			//			fmt.Printf("  %v\n", cache)
	//			cacheStr := fmt.Sprintf("%v\n", cache)
	//			retHardwareStr = retHardwareStr + cacheStr
	//		}
	//	}

	// --------------------------------------
	//	fmt.Println("--------------Network-------------")
	//	fmt.Println("")
	networkHeadStr := fmt.Sprintf("--------------Network-------------\n")
	retHardwareStr = retHardwareStr + networkHeadStr

	net, err := ghw.Network()
	if err != nil {
		fmt.Printf("Error getting network info: %v", err)
	}

	//	//	fmt.Printf("%v\n", net)
	//	netStr := fmt.Sprintf("%v\n", net)
	//	retHardwareStr = retHardwareStr + netStr
	//	for _, nic := range net.NICs {
	//		//		fmt.Printf(" %v\n", nic)
	//		nicStr := fmt.Sprintf("%v\n", nic)
	//		retHardwareStr = retHardwareStr + nicStr
	//	}

	//	for _, nic := range net.NICs {
	//		//		fmt.Printf("%s, %s, %s, %s \n", nic.Name, nic.MacAddress, utils.Btoa(nic.IsVirtual), nic.Vendor)
	//		addressStr := fmt.Sprintf("mac:%s\n", nic.MacAddress)
	//		retHardwareStr = retHardwareStr + addressStr
	//	}

	if net != nil {
		hmMac := make(map[string]string)
		macList := make([]string, 0)
		for _, nic := range net.NICs {
			mac := nic.MacAddress
			if mac != "" {
				macTmp := hmMac[mac]
				if macTmp == "1" {
					// 已经存在，则continue，主要发生于vlan场合
					continue
				}
				//				nicStr := fmt.Sprintf("%v\n", nic)
				//				retHardwareStr = retHardwareStr + nicStr
				//				addressStr := fmt.Sprintf("mac:%s\n", mac)
				//				retHardwareStr = retHardwareStr + addressStr

				hmMac[mac] = "1"
				macList = append(macList, mac)
			}
		}

		// 对mac排序，这样网卡顺序发生变化时不会影响机器码
		if len(macList) > 0 {
			sort.Strings(macList)
			for _, v := range macList {
				mac := fmt.Sprintf("mac:%s\n", v)
				retHardwareStr = retHardwareStr + mac
			}
		}
	}

	// --------------------------------------
	//	fmt.Println("--------------PCI-------------")
	//	fmt.Println("")

	//	pci, err := ghw.PCI()
	//	if err != nil {
	//		fmt.Printf("Error getting PCI info: %v", err)
	//	}

	//	for _, devClass := range pci.Classes {
	//		//		fmt.Printf(" Device class: %v ('%v')\n", devClass.Name, devClass.Id)
	//		retHardwareStr = retHardwareStr + devClass.Name
	//		retHardwareStr = retHardwareStr + devClass.Id
	//		for _, devSubclass := range devClass.Subclasses {
	//			//			fmt.Printf("    Device subclass: %v ('%v')\n", devSubclass.Name, devSubclass.Id)
	//			retHardwareStr = retHardwareStr + devSubclass.Name
	//			retHardwareStr = retHardwareStr + devSubclass.Id
	//			for _, progIface := range devSubclass.ProgrammingInterfaces {
	//				//				fmt.Printf("        Programming interface: %v ('%v')\n", progIface.Name, progIface.Id)
	//				retHardwareStr = retHardwareStr + progIface.Name
	//				retHardwareStr = retHardwareStr + progIface.Id
	//			}
	//		}
	//	}

	// --------------------------------------
	//	fmt.Println("--------------PCIDevice-------------")
	//	fmt.Println("")
	pciHeadStr := fmt.Sprintf("--------------PCIDevice-------------\n")
	retHardwareStr = retHardwareStr + pciHeadStr

	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}

	addr := "0000:00:00.0"
	if len(os.Args) == 2 {
		addr = os.Args[1]
	}
	//	fmt.Printf("PCI device information for %s\n", addr)

	deviceInfo := pci.GetDevice(addr)
	if deviceInfo == nil {
		fmt.Printf("could not retrieve PCI device information for %s\n", addr)
	}

	vendor := deviceInfo.Vendor
	//	fmt.Printf("Vendor: %s [%s]\n", vendor.Name, vendor.Id)
	vendorStr := fmt.Sprintf("Vendor: %s [%s]\n", vendor.Name, vendor.ID)
	retHardwareStr = retHardwareStr + vendorStr

	product := deviceInfo.Product
	//	fmt.Printf("Product: %s [%s]\n", product.Name, product.Id)
	productStr := fmt.Sprintf("Product: %s [%s]\n", product.Name, product.ID)
	retHardwareStr = retHardwareStr + productStr

	subsystem := deviceInfo.Subsystem
	subvendor := pci.Vendors[subsystem.VendorID]
	subvendorName := "UNKNOWN"
	if subvendor != nil {
		subvendorName = subvendor.Name
	}
	//	fmt.Printf("Subsystem: %s [%s] (Subvendor: %s)\n", subsystem.Name, subsystem.Id, subvendorName)
	subsystemStr := fmt.Sprintf("Subsystem: %s [%s] (Subvendor: %s)\n", subsystem.Name, subsystem.ID, subvendorName)
	retHardwareStr = retHardwareStr + subsystemStr

	class := deviceInfo.Class
	//	fmt.Printf("Class: %s [%s]\n", class.Name, class.Id)
	classStr := fmt.Sprintf("Class: %s [%s]\n", class.Name, class.ID)
	retHardwareStr = retHardwareStr + classStr

	subclass := deviceInfo.Subclass
	//	fmt.Printf("Subclass: %s [%s]\n", subclass.Name, subclass.Id)
	subclassStr := fmt.Sprintf("Subclass: %s [%s]\n", subclass.Name, subclass.ID)
	retHardwareStr = retHardwareStr + subclassStr

	progIface := deviceInfo.ProgrammingInterface
	//	fmt.Printf("Programming Interface: %s [%s]\n", progIface.Name, progIface.Id)
	progIfaceStr := fmt.Sprintf("Programming Interface: %s [%s]\n", progIface.Name, progIface.ID)
	retHardwareStr = retHardwareStr + progIfaceStr

	// --------------------------------------
	//	fmt.Println("--------------GPU-------------")
	//	fmt.Println("")

	gpuHeadStr := fmt.Sprintf("--------------GPU-------------\n")
	retHardwareStr = retHardwareStr + gpuHeadStr

	gpu, err := ghw.GPU()
	if err != nil {
		fmt.Printf("Error getting GPU info: %v", err)
	}

	//	fmt.Printf("%v\n", gpu)

	for _, card := range gpu.GraphicsCards {
		//		fmt.Printf(" %v\n", card)
		cardStr := fmt.Sprintf("%v\n", card)
		retHardwareStr = retHardwareStr + cardStr
	}

	//	fmt.Println(retHardwareStr)

	return retHardwareStr

}

func getCpuInfo() string {

	retCpuInfo := ""

	r, err := os.Open(PathProcCpuinfo)
	if err != nil {
		return ""
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())
		//		if strings.Contains(line, "cpu MHz") == true ||
		//			strings.Contains(line, "bogomips") == true {
		//			// 这些指标经常会变，所以需要排除
		//			continue
		//		}
		if strings.Contains(line, "model name") == true {
			retCpuInfo = retCpuInfo + line
		}

	}

	return retCpuInfo
}

package cache

import (
	"encoding/json"
	//	"fmt"
	"time"

	"caton/xw.wang/utils"

	clog "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

func newPool(url string, max_idle, idle_timeout int, db int, password string) *redis.Pool {
	timeout := time.Duration(idle_timeout) * time.Second
	return &redis.Pool{
		MaxIdle:     max_idle,
		IdleTimeout: timeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(url)
			if err != nil {
				return nil, err

			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err

				}
				//				clog.Info("[redis auth ok]")
			}

			c.Do("SELECT", db)

			//			clog.Info("[sredis pool open]")
			return c, err

		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil

			}
			_, err := c.Do("PING")
			return err

		},
	}
}

func Init(url string, max_idle, idle_timeout int, db int) {
	pool = newPool(url, max_idle, idle_timeout, db, "")
}

func InitWithPassword(url string, max_idle, idle_timeout int, db int, password string) {
	pool = newPool(url, max_idle, idle_timeout, db, password)
}

func Get() redis.Conn {
	if pool == nil {
		clog.Critical("Please set cache pool first!")
		return nil
	}
	return pool.Get()
}

func Ping() error {
	conn := Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}

func Close() {
	clog.Info("[redis pool close]")
	pool.Close()
}

func DoMapSet(key string, obj map[int]map[int]string, expire int) {
	redisConn := Get()
	defer redisConn.Close()

	value, _ := json.Marshal(obj)
	// 存入redis
	_, err := redisConn.Do("SETEX", key, expire, value)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)

}
func DoMapGet(key string) (obj map[int]map[int]string) {
	redisConn := Get()
	defer redisConn.Close()

	value, err := redis.Bytes(redisConn.Do("GET", key))
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)

	// 将json解析成map类型
	errShal := json.Unmarshal(value, &obj)
	utils.CheckErr(errShal, utils.CHECK_FLAG_LOGONLY)

	return obj
}

func DoSet(key string, obj interface{}, expire int) bool {
	redisConn := Get()
	defer redisConn.Close()

	value, _ := json.Marshal(obj)
	// 存入redis
	_, err := redisConn.Do("SETEX", key, expire, value)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}

	return true
}

func DoSetNx(key string, expire int) bool {
	redisConn := Get()
	defer redisConn.Close()

	retSetNx, errSetNx := redis.Int(redisConn.Do("SETNX", key, "1"))
	utils.CheckErr(errSetNx, utils.CHECK_FLAG_LOGONLY)
	if errSetNx != nil {
		return false
	}

	if retSetNx == 0 {
		//		fmt.Println("DoSetNx false")
		return false
	}

	_, err2 := redisConn.Do("EXPIRE", key, expire)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false
	}

	//	fmt.Println("DoSetNx true")

	return true
}

func DoHSet(key string, field string, obj interface{}, expire int) bool {
	redisConn := Get()
	defer redisConn.Close()

	value, _ := json.Marshal(obj)
	// 存入redis
	_, err := redisConn.Do("HSET", key, field, value)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}
	_, err2 := redisConn.Do("EXPIRE", key, expire)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false
	}

	return true
}

func DoHDel(key string, field string) bool {
	redisConn := Get()
	defer redisConn.Close()

	_, err := redisConn.Do("HDEL", key, field)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}

	return true
}

func DoHGet(key string, field string, obj interface{}) bool {
	redisConn := Get()
	defer redisConn.Close()

	ret, err1 := redisConn.Do("HGET", key, field)
	if ret == nil {
		return false
	}

	value, err2 := redis.Bytes(ret, err1)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false
	}

	if value == nil {
		return false
	}

	err3 := json.Unmarshal(value, obj)
	utils.CheckErr(err3, utils.CHECK_FLAG_LOGONLY)
	if err3 != nil {
		return false
	}

	return true
}

// 设置key的过期时间
func DoExpire(key string, expire int) bool {
	redisConn := Get()
	defer redisConn.Close()

	ret, err := redisConn.Do("EXPIRE", key, expire)
	//	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}
	value, _ := redis.Int(ret, err)
	if value == 1 {
		return true
	} else {
		return false
	}
}

func DoDel(key string) bool {
	redisConn := Get()
	defer redisConn.Close()

	_, err := redisConn.Do("DEL", key)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}

	return true
}

// DoGet obj:结构体指针 返回值 true：取到值 false：未取到值
func DoGet(key string, obj interface{}) bool {
	redisConn := Get()
	defer redisConn.Close()

	ret, err1 := redisConn.Do("GET", key)
	if ret == nil {
		return false
	}

	value, err2 := redis.Bytes(ret, err1)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false
	}

	err3 := json.Unmarshal(value, obj)
	utils.CheckErr(err3, utils.CHECK_FLAG_LOGONLY)
	if err3 != nil {
		return false
	}

	return true
}

func DoFlushDb() bool {
	redisConn := Get()
	defer redisConn.Close()

	_, _ = redisConn.Do("FLUSHDB")
	return true
}

func DoStrSet(key string, obj string, expire int) bool {
	redisConn := Get()
	defer redisConn.Close()

	if expire == 0 {
		// 存入redis
		ret, err := redisConn.Do("SET", key, obj)
		utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
		if ret == false {
			return false
		} else {
			return true
		}
	} else {
		// 存入redis
		ret, err := redisConn.Do("SETEX", key, expire, obj)
		utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
		if ret == false {
			return false
		} else {
			return true
		}
	}

}

func DoStrGet(key string) (ret bool, obj string) {
	redisConn := Get()
	defer redisConn.Close()

	retGet, err1 := redisConn.Do("get", key)
	if retGet == nil {
		return false, ""
	}

	if err1 == nil {
		value, err2 := redis.String(retGet, err1)
		utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
		if err2 == nil {
			return true, value
		} else {
			return false, ""
		}
	} else {
		return false, ""
	}

}

func DoKeys(key string) (ret bool, keys []string) {

	redisConn := Get()
	defer redisConn.Close()

	retGet, err1 := redisConn.Do("keys", key)
	if retGet == nil {
		return false, nil
	}

	if err1 == nil {
		value, err2 := redis.Strings(retGet, err1)
		utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
		if err2 == nil {
			return true, value
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}

}

func DoStrHSetConn(key string, field string, value string, conn redis.Conn) {

	// 存入redis
	_, err := conn.Do("hset", key, field, value)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)

}

func DoStrHGetConn(key string, field string, conn redis.Conn) (obj string) {
	ret, err1 := conn.Do("hget", key, field)
	utils.CheckErr(err1, utils.CHECK_FLAG_LOGONLY)
	if ret != nil {
		value, err2 := redis.String(ret, err1)
		utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
		return value
	}
	return ""
}

func DoDelConn(key string, conn redis.Conn) bool {
	_, err := conn.Do("DEL", key)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}
	return true
}
func DoHDelConn(key string, field string, conn redis.Conn) bool {
	_, err := conn.Do("HDEL", key)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}
	return true
}

func DoHkeys(key string) []string {
	redisConn := Get()
	defer redisConn.Close()

	ret, err := redis.Strings(redisConn.Do("hkeys", key))
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)

	return ret
}

func DoHVals(key string) (bool, []interface{}) {
	redisConn := Get()
	defer redisConn.Close()

	ret, err1 := redisConn.Do("HVals", key)
	if ret == nil {
		return false, nil
	}

	value, err2 := redis.ByteSlices(ret, err1)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false, nil
	}

	obj := make([]interface{}, 0)
	for _, v := range value {
		var f interface{}
		err3 := json.Unmarshal(v, &f)
		if err3 != nil {
			return false, nil
		}
		obj = append(obj, f)
	}

	return true, obj
}

func DoHLen(key string) int64 {
	redisConn := Get()
	defer redisConn.Close()
	// 存入redis
	ret, err := redis.Int64(redisConn.Do("hlen", key))
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)

	return ret
}

func DoSetConn(key string, obj interface{}, expire int, conn redis.Conn) bool {

	value, _ := json.Marshal(obj)
	// 存入redis
	_, err := conn.Do("SETEX", key, expire, value)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}

	return true
}

func DoGetConn(key string, obj interface{}, conn redis.Conn) bool {

	ret, err1 := conn.Do("GET", key)
	if ret == nil {
		return false
	}

	value, err2 := redis.Bytes(ret, err1)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false
	}

	err3 := json.Unmarshal(value, obj)
	utils.CheckErr(err3, utils.CHECK_FLAG_LOGONLY)
	if err3 != nil {
		return false
	}

	return true
}

func DoRPush(key string, obj interface{}, expire int) bool {
	redisConn := Get()
	defer redisConn.Close()

	value, _ := json.Marshal(obj)
	// 存入redis
	_, err := redisConn.Do("RPUSH", key, value)
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}
	_, err2 := redisConn.Do("EXPIRE", key, expire)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false
	}

	return true
}

func DoLPop(key string, obj interface{}) bool {
	redisConn := Get()
	defer redisConn.Close()

	ret, err1 := redisConn.Do("LPop", key)
	if ret == nil {
		return false
	}

	value, err2 := redis.Bytes(ret, err1)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false
	}

	if value == nil {
		return false
	}

	err3 := json.Unmarshal(value, obj)
	utils.CheckErr(err3, utils.CHECK_FLAG_LOGONLY)
	if err3 != nil {
		return false
	}

	return true
}

// 共享锁
// -------------------------------------------------------------------
// lockStart 开始一个分布式锁,retLock:是否锁成功 ，尝试n次，每次间隔100毫秒
// cntTry:尝试次数 redisKeyEx：锁的超时时间，单位：秒，该值必须大于1秒
func LockStart(redisKey string, redisKeyEx int, cntTry int) (retLock bool) {
	retLock = false
	for i := 0; i < cntTry; i++ {
		retLock = DoSetNx(redisKey, redisKeyEx)
		if retLock == true {
			break
		} else {
			// 休眠500毫秒
			time.Sleep(500 * time.Millisecond)
		}
	}

	return retLock
}

// lockEnd 结束一个分布式锁
func LockEnd(redisKey string) {
	DoDel(redisKey)
}

// lockHeart 锁的心跳,锁的超时时间很短，一旦没有心跳，锁就自动解锁
// redisKeyEx：每次心跳时会重置key的超时时间，用来保持锁定状态，该值必须大于1秒
// expire:心跳超时时间，单位：秒，如果忘记关闭心跳，超时后心跳结束
func LockHeart(redisKey string, redisKeyEx int, expire float64) {
	nowT := utils.GetNowUTC2()
	for {
		ret := DoExpire(redisKey, redisKeyEx)
		if ret == false {
			break
		}
		time.Sleep(1 * time.Second)
		// 如果超过心跳超时时间，则心跳退出
		diff := utils.DtDiff(nowT, utils.GetNowUTC2())
		if diff.Seconds() > expire {
			break
		}
	}
}

func DoZAdd(key string, score float64, obj interface{}) (int, bool) {
	redisConn := Get()
	defer redisConn.Close()

	value, _ := json.Marshal(obj)
	// 存入redis
	ret, err1 := redisConn.Do("ZAdd", key, score, value)
	utils.CheckErr(err1, utils.CHECK_FLAG_LOGONLY)
	if err1 != nil {
		return -1, false
	}

	ret2, err2 := redis.Int(ret, err1)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return -1, false
	}

	return ret2, true
}

func DoZRange(key string, start int, stop int) (bool, []interface{}) {
	redisConn := Get()
	defer redisConn.Close()

	ret, err1 := redisConn.Do("ZRange", key, start, stop)
	if ret == nil {
		return false, nil
	}

	value, err2 := redis.ByteSlices(ret, err1)
	utils.CheckErr(err2, utils.CHECK_FLAG_LOGONLY)
	if err2 != nil {
		return false, nil
	}

	obj := make([]interface{}, 0)
	for _, v := range value {
		var f interface{}
		err3 := json.Unmarshal(v, &f)
		if err3 != nil {
			return false, nil
		}
		obj = append(obj, f)
	}

	//	err3 := json.Unmarshal(value, obj)
	//	utils.CheckErr(err3, utils.CHECK_FLAG_LOGONLY)
	//	if err3 != nil {
	//		return false
	//	}

	return true, obj
}

func DoExists(key string) bool {
	redisConn := Get()
	defer redisConn.Close()

	exists, err := redis.Bool(redisConn.Do("EXISTS", key))
	utils.CheckErr(err, utils.CHECK_FLAG_LOGONLY)
	if err != nil {
		return false
	}

	return exists
}

// -------------------------------------------------------------------

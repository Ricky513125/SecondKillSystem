package data

import (
	"SecKill/conf"
	"github.com/go-redis/redis/v7"
)

var client *redis.Client

// open redis connection pool
func initRedisConnection(config conf.AppConfig) {
	// use the default connections to the database
	client = redis.NewClient(&redis.Options{
		Addr:     config.App.Redis.Address,
		Password: config.App.Redis.Password,
		DB:       0,
	})
	// clear all the keys when just connect to Redis
	if _, err := FlushAll(); err != nil {
		println("Error when flushAll. " + err.Error())
	}

}

// use for tests
func FlushAll() (string, error) {
	return client.FlushAll().Result()
}

// ensure redis to load the lua script, if not then load
func PrepareScript(script string) string {
	// sha := sha1.Sum([]byte(script))
	// check if scripts are available  // []bool
	scriptsExists, err := client.ScriptExists(script).Result()
	if err != nil {
		panic("Failed to check if script exists: " + err.Error())
	}
	if !scriptsExists[0] {
		scriptSHA, err := client.ScriptLoad(script).Result()
		if err != nil {
			panic("Failed to load script " + script + " err: " + err.Error())
		}
		return scriptSHA
	}
	print("Script Exists.")
	return ""
}

// execute lua script by SHA
func EvalSHA(sha string, args []string) (interface{}, error) {
	// Result() *cmd() -> (interface{}, error)
	val, err := client.EvalSha(sha, args).Result()
	if err != nil {
		print("Error executing evalSHA... " + err.Error())
		return nil, err
	}
	return val, nil
}

// redis operation SET
func SetForever(key string, value interface{}) (string, error) {
	// set a key-value pair in Redis with no expiration time(0).
	// It turns the result of the SET command
	val, err := client.Set(key, value, 0).Result() // expiration: no expire time
	return val, err
}

// redis operation hmset
// This function sets multiple fields of a Redis hash.
// The hash is stored under the given key and contains the key-value pairs specified
// in the field map.
func SetMapForever(key string, field map[string]interface{}) (string, error) {
	return client.HMSet(key, field).Result()
}

// redis operation hmget
func GetMap(key string, fields ...string) ([]interface{}, error) {
	return client.HMGet(key, fields...).Result()
}

// redis SADD
// This function adds a value (field) to a Redis set (key).
// If the field is already a member of the set, Redis will ignore it.
func SetAdd(key string, field string) (int64, error) {
	return client.SAdd(key, field).Result()
}

// redis SISMEMBER
func SetIsMember(key string, field string) (bool, error) {
	return client.SIsMember(key, field).Result()
}

// redis SMEMBERS
func GetSetMembers(key string) ([]string, error) {
	return client.SMembers(key).Result()
}

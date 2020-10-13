package public

import (
	"github.com/garyburd/redigo/redis"
	"src/gatewayProject/golang_common/lib"
)

func RedisConfPipeline(pip ...func(c redis.Conn)) error {
	conn, err := lib.RedisConnFactory("default")
	if err != nil {
		return err
	}
	defer conn.Close()
	for _, f := range pip {
		f(conn)
	}
	_ = conn.Flush()
	return nil
}

func RedisConfDo(commandName string, args ...interface{}) (interface{}, error) {
	conn, err := lib.RedisConnFactory("default")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Do(commandName, args...)
}

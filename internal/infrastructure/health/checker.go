package health

import (
	"context"
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
)

type Checker interface {
	Liveness() map[string]string
	Readiness(ctx context.Context) (map[string]string, bool)
}

type checker struct {
	mysql *sql.DB
	redis *redis.Client
}

func NewChecker(mysql *sql.DB, redis *redis.Client) Checker {
	return &checker{mysql: mysql, redis: redis}
}

func (c *checker) Liveness() map[string]string {
	return map[string]string{"status": "ok"}
}

func (c *checker) Readiness(ctx context.Context) (map[string]string, bool) {
	status := map[string]string{
		"status": "ok",
		"mysql":  "ok",
		"redis":  "ok",
	}

	mysqlCtx, mysqlCancel := context.WithTimeout(ctx, 2*time.Second)
	defer mysqlCancel()
	if err := c.mysql.PingContext(mysqlCtx); err != nil {
		status["status"] = "degraded"
		status["mysql"] = "down"
	}

	redisCtx, redisCancel := context.WithTimeout(ctx, 2*time.Second)
	defer redisCancel()
	if err := c.redis.Ping(redisCtx).Err(); err != nil {
		status["status"] = "degraded"
		status["redis"] = "down"
	}

	return status, status["status"] == "ok"
}

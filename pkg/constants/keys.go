package constants

// MySQL configuration keys in Consul
const (
	MySQLHostKey               = "config/mysql/host"
	MySQLPortKey               = "config/mysql/port"
	MySQLUserKey               = "config/mysql/user"
	MySQLPasswordKey           = "config/mysql/password"
	MySQLDBKey                 = "config/mysql/db"
	MySQLMaxLifetimeSecondsKey = "config/mysql/connection/maxLifeTime"
	MySQLIdleConnectionsKey    = "config/mysql/connection/idleConnections"
	MySQLMaxOpenConnectionsKey = "config/mysql/connection/maxOpenConnections"
)

// Redis configuration keys in Consul
const (
	RedisHostKey     = "config/redis/host"
	RedisPortKey     = "config/redis/port"
	RedisPasswordKey = "config/redis/password"
)

// Environment variable keys
const (
	EnvKeyConsulAddress = "CONSUL_ADDRESS"
	EnvKeyConsulScheme  = "CONSUL_SCHEME"
	EnvKeyVaultAddr     = "VAULT_ADDR"
	EnvKeyVaultJWTPath  = "VAULT_JWT_PATH"
	EnvKeyVaultToken    = "VAULT_TOKEN"
)

const (
	VaultJWTSecretKey        = "jwtSecret"
	VaultJWTRefreshSecretKey = "jwtRefreshSecret"
)

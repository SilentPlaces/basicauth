package constants

// MySQL configuration keys in Consul
const (
	MySQLHostKey               = "config/mysql/connection/host"
	MySQLPortKey               = "config/mysql/connection/port"
	MySQLUserKey               = "config/mysql/connection/user"
	MySQLPasswordKey           = "config/mysql/connection/password"
	MySQLDBKey                 = "config/mysql/connection/db"
	MySQLMaxLifetimeSecondsKey = "config/mysql/connection/maxLifeTime"
	MySQLIdleConnectionsKey    = "config/mysql/connection/idleConnections"
	MySQLMaxOpenConnectionsKey = "config/mysql/connection/maxOpenConnections"
)

// Redis configuration keys in Consul
const (
	RedisHostKey     = "config/redis/connection/host"
	RedisPortKey     = "config/redis/connection/port"
	RedisPasswordKey = "config/redis/connection/password"
)

// config for smtp server
const (
	SMTPHostKey     = "config/smtp/connection/host"
	SMTPPortKey     = "config/smtp/connection/port"
	SMTPUsernameKey = "config/smtp/connection/username"
	SMTPPasswordKey = "config/smtp/connection/password"
)

// Environment variable keys
const (
	EnvKeyConsulAddress = "CONSUL_ADDRESS"
	EnvKeyConsulScheme  = "CONSUL_SCHEME"
	EnvKeyVaultAddr     = "VAULT_ADDR"
	EnvKeyMountPath     = "VAULT_MOUNT_PATH"
	EnvKeySecretPath    = "VAULT_SECRET_PATH"
	EnvKeyVaultToken    = "VAULT_TOKEN"
)

const (
	VaultJWTSecretKey        = "jwtSecret"
	VaultJWTRefreshSecretKey = "jwtRefreshSecret"
)

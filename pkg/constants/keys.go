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

// config for mail server
const (
	SMTPHostKey     = "config/mail/connection/host"
	SMTPPortKey     = "config/mail/connection/port"
	SMTPUsernameKey = "config/mail/connection/username"
	SMTPPasswordKey = "config/mail/connection/password"
)

// General config
const (
	GeneralDomainKey           = "config/general/domain"
	GeneralHTTPListenerPortKey = "config/general/httpListenerPort"
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

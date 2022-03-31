package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigurationFromEnv(t *testing.T) {

	// t.Run("predefined values", func(t *testing.T) {
	// 	cfg := NewConfig()
	// 	require.Equal(t, "127.0.0.1:8095", cfg.HTTPHostAddr)
	// })

	t.Run("from env variables", func(t *testing.T) {
		_ = os.Setenv("MYSQL_URL", "goLogTest")

		cfg := NewConfig()
		require.Equal(t, "goLogTest", cfg.Databases.PostgresURL)
	})
}

package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	l, err := New("scouts3000", "test", "demoTest", "today", "0", "000000")
	require.NoError(t, err)
	l.Info("infoTest")

	t.Run("prod version", func(t *testing.T) {
		l, err = New("scouts3000", "prod", "prodTest", "today", "0", "000000")
		require.NoError(t, err)
		l.Info("infoProd")
	})

	t.Run("dev version", func(t *testing.T) {
		l, err = New("scouts3000", "dev", "prodTest", "today", "0", "000000")
		require.NoError(t, err)
		l.Info("infoDev")
	})

	t.Run("error version", func(t *testing.T) {
		_, err = New("scouts3000", "prod", "prodTest", "today", "0", "000000")
		require.NoError(t, err)
	})

	t.Run("debug version", func(t *testing.T) {
		l, err = New("scouts3000", "debug", "debugTest", "today", "0", "000000")
		require.NoError(t, err)
		l.Info("infoDebug")
	})

}

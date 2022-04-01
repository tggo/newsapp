package helpers

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func TestRegisterErrorChannel(t *testing.T) {

	ch := RegisterErrorChannel(func() error {
		return errors.New("some error")
	})

	httpServerErrCh := <-ch

	require.Equal(t, "some error", httpServerErrCh.Error())
}

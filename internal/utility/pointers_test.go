package utility

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ptr(t *testing.T) {
	i := 42
	s := "hello"
	f := 3.14
	b := true

	iptr := Ptr(i)
	sptr := Ptr(s)
	fptr := Ptr(f)
	bptr := Ptr(b)

	require.NotNil(t, iptr)
	require.NotNil(t, sptr)
	require.NotNil(t, fptr)
	require.NotNil(t, bptr)

	require.Equal(t, i, *iptr)
	require.Equal(t, s, *sptr)
	require.Equal(t, f, *fptr)
	require.Equal(t, b, *bptr)
}

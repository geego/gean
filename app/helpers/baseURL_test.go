package helpers

import (
	"testing"

	"github.com/gostores/require"
)

func TestBaseURL(t *testing.T) {
	b, err := newBaseURLFromString("http://example.com")
	require.NoError(t, err)
	require.Equal(t, "http://example.com", b.String())

	p, err := b.WithProtocol("webcal://")
	require.NoError(t, err)
	require.Equal(t, "webcal://example.com", p)

	p, err = b.WithProtocol("webcal")
	require.NoError(t, err)
	require.Equal(t, "webcal://example.com", p)

	_, err = b.WithProtocol("mailto:")
	require.Error(t, err)

	b, err = newBaseURLFromString("mailto:hugo@rules.com")
	require.NoError(t, err)
	require.Equal(t, "mailto:hugo@rules.com", b.String())

	// These are pretty constructed
	p, err = b.WithProtocol("webcal")
	require.NoError(t, err)
	require.Equal(t, "webcal:hugo@rules.com", p)

	p, err = b.WithProtocol("webcal://")
	require.NoError(t, err)
	require.Equal(t, "webcal://hugo@rules.com", p)

	// Test with "non-URLs". Some people will try to use these as a way to get
	// relative URLs working etc.
	b, err = newBaseURLFromString("/")
	require.NoError(t, err)
	require.Equal(t, "/", b.String())

	b, err = newBaseURLFromString("")
	require.NoError(t, err)
	require.Equal(t, "", b.String())

}

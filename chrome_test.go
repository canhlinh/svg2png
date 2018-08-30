package svg2png

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChrome(t *testing.T) {
	chrome := NewChrome()
	assert.NotNil(t, chrome)
}

func TestScreenshoot(t *testing.T) {
	chrome := NewChrome().SetHeight(600).SetWith(600)
	require.NotNil(t, chrome)

	require.Error(t, chrome.Screenshoot("https://www.chromestatus.com/", ""))
	require.Error(t, chrome.Screenshoot("https://www.chromestatus.com/", "/"))
	require.Error(t, chrome.Screenshoot("https://www.chromestatus.com/", "."))
	require.Error(t, chrome.Screenshoot("https://www.chromestatus.com/", "./abc"))

	filepath := "Soccerball_mask_transparent_background.png"
	err := chrome.Screenshoot("https://upload.wikimedia.org/wikipedia/commons/8/84/Example.svg", filepath)
	require.NoError(t, err)

	file, err := os.Open(filepath)
	require.NoError(t, err)

	hash := md5.New()
	io.Copy(hash, file)
	assert.Equal(t, "29528ee2bac0175e1bb62f3b2cec992b", fmt.Sprintf("%x", (hash.Sum(nil))))
	os.RemoveAll(filepath)
}

func BenchmarkScreenshootRunParallel(b *testing.B) {
	chrome := NewChrome().SetHeight(600).SetWith(600).SetTimeout(time.Minute)
	require.NotNil(b, chrome)

	b.ResetTimer()
	b.SetParallelism(3)
	rand.Seed(time.Now().UnixNano())

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			filepath := fmt.Sprintf("test_file_%d.png", rand.Intn(10000000000))
			b.StartTimer()
			assert.NoError(b, chrome.Screenshoot("https://upload.wikimedia.org/wikipedia/commons/8/84/Example.svg", filepath))
			b.StopTimer()

			file, err := os.Open(filepath)
			require.NoError(b, err)
			hash := md5.New()
			io.Copy(hash, file)
			assert.Equal(b, "29528ee2bac0175e1bb62f3b2cec992b", fmt.Sprintf("%x", (hash.Sum(nil))))
			os.RemoveAll(filepath)
		}
	})

}

func BenchmarkScreenshoot(b *testing.B) {
	chrome := NewChrome().SetHeight(600).SetWith(600).SetTimeout(time.Minute)
	require.NotNil(b, chrome)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		filepath := fmt.Sprintf("test_file_%d.png", n)
		b.StartTimer()
		assert.NoError(b, chrome.Screenshoot("https://upload.wikimedia.org/wikipedia/commons/8/84/Example.svg", filepath))
		b.StopTimer()

		file, err := os.Open(filepath)
		require.NoError(b, err)
		hash := md5.New()
		io.Copy(hash, file)
		assert.Equal(b, "29528ee2bac0175e1bb62f3b2cec992b", fmt.Sprintf("%x", (hash.Sum(nil))))
		os.RemoveAll(filepath)
	}

}

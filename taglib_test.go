package taglib_test

import (
	_ "embed"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"go.senan.xyz/taglib-wasm"
	_ "go.senan.xyz/taglib-wasm/embed"
)

//go:embed testdata/eg.flac
var egFlac []byte

func TestFile(t *testing.T) {
	path := tmpf(t, egFlac, "eg.flac")
	f, err := taglib.New(path)
	nilErr(t, err)

	tags := f.ReadTags()
	eq(t, len(tags["ARTIST"]), 2)
	eq(t, tags["ARTIST"][0], "Artist One")
	eq(t, tags["ARTIST"][1], "Artist Two")

	eq(t, 1*time.Second, f.Length())
	eq(t, 1460, f.Bitrate())
	eq(t, 48_000, f.SampleRate())
	eq(t, 2, f.Channels())

	f.Close()
}

func TestWrite(t *testing.T) {
	path := tmpf(t, egFlac, "eg.flac")

	tags := map[string][]string{
		"ONE":   {"one", "two", "three", "four"},
		"FIVE":  {"six", "seven"},
		"EIGHT": {"nine"},
	}

	{
		f, err := taglib.New(path)
		nilErr(t, err)

		f.WriteTags(tags)

		err = f.Save()
		nilErr(t, err)

		f.Close()
	}
	{
		f, err := taglib.New(path)
		nilErr(t, err)

		got := f.ReadTags()
		if !maps.EqualFunc(tags, got, slices.Equal) {
			t.Errorf("%v != %v", got, tags)
		}

		f.Close()
	}
}

func TestMultiOpen(t *testing.T) {
	{
		path := tmpf(t, egFlac, "eg.flac")
		f, err := taglib.New(path)
		nilErr(t, err)
		f.Close()
	}
	{
		path := tmpf(t, egFlac, "eg.flac")
		f, err := taglib.New(path)
		nilErr(t, err)
		f.Close()
	}
}

func BenchmarkOpen(b *testing.B) {
	path := tmpf(b, egFlac, "eg.flac")
	b.ResetTimer()

	for range b.N {
		f, err := taglib.New(path)
		nilErr(b, err)
		_ = f.ReadTags()
		f.Close()
	}
}

func tmpf(t testing.TB, b []byte, name string) string {
	p := filepath.Join(t.TempDir(), name)
	err := os.WriteFile(p, b, os.ModePerm)
	nilErr(t, err)
	return p
}

func nilErr(t testing.TB, err error) {
	if err != nil {
		t.Helper()
		t.Fatalf("err: %v", err)
	}
}
func eq[T comparable](t testing.TB, a, b T) {
	if a != b {
		t.Helper()
		t.Fatalf("%v != %v", a, b)
	}
}

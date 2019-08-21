package uuid

import (
	"testing"

	goog "github.com/google/uuid"
)

func TestNilParse(t *testing.T) {
	var nilUUID UUID
	u, err := Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatal(err)
	}
	if u != nilUUID {
		t.Fatal("parse nil uuid failed")
	}
}

func TestParseToString(t *testing.T) {
	str := "15588635-a45e-4867-aadb-dbf0385ade95"
	u, err := Parse(str)
	if err != nil {
		t.Fatal(err)
	}
	res := u.String()
	if res != str {
		t.Errorf("parse/string roundtrip failed, expected %s, got %s", str, res)
	}
}

func TestParseToAppend(t *testing.T) {
	str := "15588635-a45e-4867-aadb-dbf0385ade95"
	u, err := Parse(str)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 36)
	res := u.AppendFormatted(buf[0:0])
	if string(res) != str {
		t.Errorf("parse/string roundtrip failed, expected %s, got %s", str, res)
	}
}

var (
	examples = []struct {
		namespace string
		input     string
		expected  string
	}{
		{
			"15588635-a45e-4867-aadb-dbf0385ade95",
			"input 1",
			"4c816dc1-9418-502e-9b91-f17b83891bf8",
		},
		{
			"15588635-a45e-4867-aadb-dbf0385ade95",
			"input 123456123456123456123456123456123456123456123456123456123456123456123456123456123456",
			"0a550dec-ba3e-52c2-a54d-60795ea13c90",
		},
	}
)

func doTest(t *testing.T, f func(ns, input string) string) {
	t.Helper()
	for _, example := range examples {
		result := f(example.namespace, example.input)
		if result != example.expected {
			t.Errorf("expected %s with ns %s to output %s but got %s", example.input, example.namespace, example.expected, result)
		}
	}
}

func TestExamples(t *testing.T) {
	doTest(t, func(namespace, input string) string {
		ns, err := Parse(namespace)
		if err != nil {
			t.Fatal(err)
		}
		var u UUID
		NewSHA1Gen(ns).Generate(&u, []byte(input))
		return u.String()
	})
}

// This is simply to test that the widely used google implementation gets the
// same test results as our implementation.
func TestExamplesGoogle(t *testing.T) {
	doTest(t, func(namespace, input string) string {
		return goog.NewSHA1(goog.MustParse(namespace), []byte(input)).String()
	})
}

func BenchmarkGenerate(b *testing.B) {
	for _, input := range []struct {
		name  string
		input []byte
	}{
		{"short", []byte("example.com")},
		{"longer", []byte("benchmark input string 15588635-a45e-4867-aadb-dbf0385ade95 15588635-a45e-4867-aadb-dbf0385ade95 15588635-a45e-4867-aadb-dbf0385ade95 15588635-a45e-4867-aadb-dbf0385ade95")},
	} {

		b.Run(input.name, func(b *testing.B) {
			b.ReportAllocs()

			ns, err := Parse("15588635-a45e-4867-aadb-dbf0385ade95")
			if err != nil {
				b.Fatal(err)
			}

			gen := NewSHA1Gen(ns)

			b.ResetTimer()

			var u UUID
			for i := 0; i < b.N; i++ {
				gen.Generate(&u, input.input)
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	b.ReportAllocs()

	example, err := Parse("15588635-a45e-4867-aadb-dbf0385ade95")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = example.String()
	}

}
func BenchmarkAppendBytes(b *testing.B) {
	b.ReportAllocs()

	example, err := Parse("15588635-a45e-4867-aadb-dbf0385ade95")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	buf := make([]byte, 36)
	for i := 0; i < b.N; i++ {
		buf = example.AppendFormatted(buf[0:0])
	}

}

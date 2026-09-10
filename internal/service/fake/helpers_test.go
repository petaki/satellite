package fake

import "testing"

func TestArg(t *testing.T) {
	args := []any{"first", []byte("second"), 3}

	cases := map[int]string{0: "first", 1: "second", 2: "3", 9: ""}

	for index, want := range cases {
		if got := arg(args, index); got != want {
			t.Errorf("arg(args, %d) = %q, want %q", index, got, want)
		}
	}
}

func TestGlobMatch(t *testing.T) {
	cases := []struct {
		pattern string
		value   string
		want    bool
	}{
		{"*", "anything", true},
		{"*", "", true},
		{"", "", true},
		{"", "x", false},
		{"exact", "exact", true},
		{"exact", "exacts", false},
		{"exac", "exact", false},

		// findPaths: prefix scans
		{"web-01:disk:200:*", "web-01:disk:200:Lw==", true},
		{"web-01:disk:200:*", "web-01:disk:100:Lw==", false},
		{"web-01:disk:200:*", "web-02:disk:200:Lw==", false},

		// FindAll: a wildcard on both sides
		{"*cpu:*", "web-01:cpu:1788991200", true},
		{"*cpu:*", "web-01:process:cpu:1788991200", true},
		{"*cpu:*", "web-01:memory:1788991200", false},

		// Delete: probe prefix
		{"web-01:*", "web-01:cpu:1", true},
		{"web-01:*", "web-011:cpu:1", false},

		// escape(): a probe name containing glob metacharacters must match literally
		{`web\*01:*`, "web*01:cpu:1", true},
		{`web\*01:*`, "webXX01:cpu:1", false},
		{`web\?01:*`, "web?01:cpu:1", true},
		{`web\?01:*`, "webX01:cpu:1", false},

		// unescaped metacharacters still glob
		{"web?01:*", "webX01:cpu:1", true},
		{"web*01:*", "webANYTHING01:cpu:1", true},
	}

	for _, c := range cases {
		if got := globMatch(c.pattern, c.value); got != c.want {
			t.Errorf("globMatch(%q, %q) = %v, want %v", c.pattern, c.value, got, c.want)
		}
	}
}

func TestToString(t *testing.T) {
	cases := []struct {
		value any
		want  string
	}{
		{"text", "text"},
		{[]byte("bytes"), "bytes"},
		{7, "7"},
		{int64(8), "8"},
		{1.5, ""},
		{nil, ""},
	}

	for _, c := range cases {
		if got := toString(c.value); got != c.want {
			t.Errorf("toString(%#v) = %q, want %q", c.value, got, c.want)
		}
	}
}

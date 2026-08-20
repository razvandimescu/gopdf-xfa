package main

import "testing"

// The abbreviation rules are ordered, and bytes.ReplaceAll consumes what it
// matches, so a shorter pattern placed first makes a longer one unreachable.
// That is not visible by reading either rule alone — only by reading the whole
// list in order — which is what these cases pin.
func TestPostProcessFixesAbbreviations(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{
			// The regression: " ct." used to match the space inside "din ct."
			// first, stranding the Romanian "din" in front of "acc.".
			name: "din ct. at the start of an element body",
			in:   "<span>din ct. 191</span>",
			want: "<span>from acc. 191</span>",
		},
		{
			name: "din ct. after a space",
			in:   "<span>total din ct. 191</span>",
			want: "<span>total from acc. 191</span>",
		},
		{
			name: "din ct. in parentheses",
			in:   "<span>(din ct. 191)</span>",
			want: "<span>(from acc. 191)</span>",
		},
		{
			name: "bare ct. is still translated",
			in:   "<span>ct. 191 si ct. 192</span>",
			want: "<span>acc. 191 si acc. 192</span>",
		},
		{
			name: "ct. inside a formula keeps its operator",
			in:   "<span>ct. 191+ct. 192-ct. 193,ct. 194</span>",
			want: "<span>acc. 191+acc. 192-acc. 193,acc. 194</span>",
		},
		{
			name: "rd. becomes row without doubling the space",
			in:   "<span>rd. 191 la 192</span>",
			want: "<span>row 191 to 192</span>",
		},
		{
			// The period carries the separation when no space follows.
			name: "rd. with no following space",
			in:   "<span>(rd.191)</span>",
			want: "<span>(row 191)</span>",
		},
		{
			// Scripts are partitioned out and must come back untouched,
			// or a rule would corrupt executable form logic.
			name: "script bodies are left alone",
			in:   "<span>ct. 1</span><script>var x = \" ct. 1\";</script>",
			want: "<span>acc. 1</span><script>var x = \" ct. 1\";</script>",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := string(postProcessFixes([]byte(c.in))); got != c.want {
				t.Errorf("postProcessFixes(%q)\n got %q\nwant %q", c.in, got, c.want)
			}
		})
	}
}

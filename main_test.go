package main

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func strip(s string) string { return ansi.ReplaceAllString(s, "") }

var today = time.Date(2026, time.November, 17, 0, 0, 0, 0, time.UTC)

func TestNov2026StartsSunday(t *testing.T) {
	lines := monthLines(2026, time.November, today, false, true, false)
	if !strings.HasPrefix(lines[2], " 1  2") {
		t.Fatalf("first week = %q", lines[2])
	}
	if !strings.Contains(strings.Join(lines, "\n"), "30") || strings.Contains(strings.Join(lines[2:], "\n"), "31") {
		t.Fatal("November should have 30 days")
	}
}

func TestMondayFirst(t *testing.T) {
	lines := monthLines(2026, time.November, today, true, true, false)
	if !strings.HasPrefix(lines[1], "Mo") {
		t.Fatalf("header = %q", lines[1])
	}
	if !strings.HasSuffix(lines[2], " 1") {
		t.Fatalf("first week = %q", lines[2])
	}
}

func TestLeapFebruary(t *testing.T) {
	body := strings.Join(monthLines(2028, time.February, today, false, true, false)[2:], "\n")
	if !strings.Contains(body, "29") || strings.Contains(body, "30") {
		t.Fatalf("Feb 2028 should have 29 days:\n%s", body)
	}
}

func TestWidthWithColor(t *testing.T) {
	for m := time.January; m <= time.December; m++ {
		for _, mon := range []bool{false, true} {
			lines := monthLines(2026, m, today, mon, true, true)
			if len(lines) != 8 {
				t.Fatalf("%v: %d lines", m, len(lines))
			}
			for i, l := range lines {
				if n := len(strip(l)); n != width {
					t.Errorf("%v line %d: width %d: %q", m, i, n, strip(l))
				}
			}
		}
	}
}

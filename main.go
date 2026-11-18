// gocal is a colorful take on the classic Unix cal command.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	width = 20 // visible width of one month block: 7 cells * 2 + 6 spaces
	gap   = 3  // spaces between side-by-side months
)

// One rainbow color per month (ANSI 256-color palette).
var monthColors = [12]int{196, 208, 220, 154, 46, 49, 51, 39, 27, 93, 165, 201}

// Fixed-date holidays, keyed by {month, day}.
var holidays = map[[2]int]bool{
	{1, 1}: true, {2, 14}: true, {7, 4}: true, {10, 31}: true, {12, 25}: true,
}

const (
	styleHeader   = "2;36"
	styleSunday   = "31"
	styleSaturday = "35"
	styleHoliday  = "4;32"
	styleToday    = "1;30;103"
)

type painter bool

// paint wraps already-padded text in an ANSI style so alignment is unaffected.
func (p painter) paint(code, text string) string {
	if !p || code == "" {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}

func center(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	left := (w - len(s)) / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", w-len(s)-left)
}

func weekendStyle(wd time.Weekday) string {
	switch wd {
	case time.Sunday:
		return styleSunday
	case time.Saturday:
		return styleSaturday
	}
	return ""
}

// monthLines renders one month as 8 lines, each exactly `width` visible chars.
func monthLines(year int, month time.Month, today time.Time, mondayFirst, showYear bool, p painter) []string {
	lines := make([]string, 0, 8)

	title := month.String()
	if showYear {
		title += " " + strconv.Itoa(year)
	}
	lines = append(lines, p.paint(fmt.Sprintf("1;38;5;%d", monthColors[month-1]), center(title, width)))

	colDay := func(col int) time.Weekday {
		if mondayFirst {
			return time.Weekday((col + 1) % 7)
		}
		return time.Weekday(col)
	}

	header := make([]string, 7)
	for col := range 7 {
		wd := colDay(col)
		style := weekendStyle(wd)
		if style == "" {
			style = styleHeader
		}
		header[col] = p.paint(style, wd.String()[:2])
	}
	lines = append(lines, strings.Join(header, " "))

	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	offset := int(first.Weekday())
	if mondayFirst {
		offset = (offset + 6) % 7
	}
	days := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	isThisMonth := today.Year() == year && today.Month() == month

	for row := range 6 {
		cells := make([]string, 7)
		for col := range 7 {
			d := row*7 + col - offset + 1
			if d < 1 || d > days {
				cells[col] = "  "
				continue
			}
			style := weekendStyle(colDay(col))
			if holidays[[2]int{int(month), d}] {
				style = styleHoliday
			}
			if isThisMonth && today.Day() == d {
				style = styleToday
			}
			cells[col] = p.paint(style, fmt.Sprintf("%2d", d))
		}
		lines = append(lines, strings.Join(cells, " "))
	}
	return lines
}

// joinColumns lays month blocks out side by side.
func joinColumns(blocks [][]string) []string {
	out := make([]string, len(blocks[0]))
	sep := strings.Repeat(" ", gap)
	for i := range out {
		parts := make([]string, len(blocks))
		for j, b := range blocks {
			parts[j] = b[i]
		}
		out[i] = strings.TrimRight(strings.Join(parts, sep), " ")
	}
	return out
}

func printLines(lines []string) {
	for _, l := range lines {
		fmt.Println(l)
	}
}

func printYear(year int, today time.Time, mondayFirst bool, p painter) {
	total := 3*width + 2*gap
	banner := center(strconv.Itoa(year), total)
	fmt.Println(p.paint("1;38;5;"+strconv.Itoa(monthColors[(year%12+12)%12]), strings.TrimRight(banner, " ")))
	fmt.Println()
	for q := range 4 {
		blocks := make([][]string, 3)
		for i := range 3 {
			m := time.Month(q*3 + i + 1)
			blocks[i] = monthLines(year, m, today, mondayFirst, false, p)
		}
		printLines(joinColumns(blocks))
		if q < 3 {
			fmt.Println()
		}
	}
}

func useColor(disabled bool) bool {
	if disabled || os.Getenv("NO_COLOR") != "" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	enableVT()
	return true
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "gocal: "+format+"\n", args...)
	os.Exit(1)
}

func parseYear(s string) int {
	y, err := strconv.Atoi(s)
	if err != nil || y < 1 || y > 9999 {
		fail("year %q not in range 1..9999", s)
	}
	return y
}

func main() {
	three := flag.Bool("3", false, "show previous, current and next month")
	wholeYear := flag.Bool("y", false, "show the whole year")
	monday := flag.Bool("m", false, "start weeks on Monday")
	noColor := flag.Bool("no-color", false, "disable colors")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: gocal [-3] [-y] [-m] [-no-color] [[month] year]")
		flag.PrintDefaults()
	}
	flag.Parse()

	today := time.Now()
	year, month := today.Year(), today.Month()
	args := flag.Args()
	switch len(args) {
	case 0:
	case 1:
		year = parseYear(args[0])
		*wholeYear = true
	case 2:
		m, err := strconv.Atoi(args[0])
		if err != nil || m < 1 || m > 12 {
			fail("month %q not in range 1..12", args[0])
		}
		month = time.Month(m)
		year = parseYear(args[1])
	default:
		flag.Usage()
		os.Exit(2)
	}

	p := painter(useColor(*noColor))

	switch {
	case *wholeYear:
		printYear(year, today, *monday, p)
	case *three:
		blocks := make([][]string, 3)
		for i := range 3 {
			t := time.Date(year, month+time.Month(i-1), 1, 0, 0, 0, 0, time.UTC)
			blocks[i] = monthLines(t.Year(), t.Month(), today, *monday, true, p)
		}
		printLines(joinColumns(blocks))
	default:
		printLines(monthLines(year, month, today, *monday, true, p))
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gdamore/tcell"
	"os"
	"time"
	"unicode"
	"math/rand"
)

var (
	opt_version = flag.Bool("version", false, "print version information")
	opt_delay   = flag.Int("delay", 10, "set delay in ms")
	opt_random  = flag.Bool("random", true, "randomize reveal")
	opt_auto    = flag.Bool("auto", false, "no user interaction")

	charset string // String containing the printable ASCII characters for scrambling
)

type NmsChar struct {
	ch    rune
	scram rune
}

func init() {
	for i := 1; i <= 127; i++ {
		if unicode.IsPrint(rune(i)) && i != ' ' {
			charset += string(rune(i))
		}
	}
	rand.Seed(time.Now().UnixNano())
}

func nms_scramble(c *NmsChar) {
	if c.ch != ' ' {
		c.scram = rune(charset[rand.Intn(len(charset))])
	} else {
		c.scram = ' '
	}
}

func nms_reveal(c *NmsChar) {
	c.scram = c.ch
}

func delay(ms time.Duration) {
	time.Sleep(ms * time.Millisecond)
}

func nms_read_stdin() []string {
	lines := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func nms_read_file(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func nms_process_input(input string) []NmsChar {
	nms_chars := make([]NmsChar, len(input))
	for i, ch := range input {
		nms_chars[i] = NmsChar{ch: ch}
		nms_scramble(&nms_chars[i])
	}
	return nms_chars
}

func main() {
	flag.Parse()

	if *opt_version {
		fmt.Println("nms version 0.3.0")
		os.Exit(0)
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	// Read the entire input into a single string
	var inputLines []string
	if len(flag.Args()) > 0 {
		// Read from the file if a filename is provided as a command-line argument
		inputLines = nms_read_file(flag.Arg(0))
	} else {
		// Otherwise, read from stdin
		inputLines = nms_read_stdin()
	}

	// Process each line
	nms_lines := make([][]NmsChar, len(inputLines))
	for i, line := range inputLines {
		nms_lines[i] = nms_process_input(line)
	}

	// Scramble and display all lines
	for y, nms_chars := range nms_lines {
		for x, ch := range nms_chars {
			screen.SetContent(x, y, ch.scram, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
	}
	screen.Show()
	delay(time.Duration(*opt_delay))

	// Prepare the list of all characters to reveal
	type Pos struct{ x, y int }
	var allChars []Pos
	for y, nms_chars := range nms_lines {
		for x := range nms_chars {
			allChars = append(allChars, Pos{x, y})
		}
	}

	// If random flag is set, shuffle the characters
	if *opt_random {
		rand.Shuffle(len(allChars), func(i, j int) { allChars[i], allChars[j] = allChars[j], allChars[i] })
	}

	// Reveal all characters
	for _, pos := range allChars {
		nms_reveal(&nms_lines[pos.y][pos.x])
		screen.SetContent(pos.x, pos.y, nms_lines[pos.y][pos.x].scram, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.Show()
		delay(time.Duration(*opt_delay))
	}

	// Wait for a key press before exiting
	if !*opt_auto {
		screen.PollEvent()
	}
}

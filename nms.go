package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gdamore/tcell"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var (
	opt_version = flag.Bool("version", false, "print version information")
	opt_delay   = flag.Int("delay", 1000, "set delay in ms")
	opt_random  = flag.Bool("random", false, "randomize reveal")
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

func nms_read_stdin() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input
}

func nms_read_file(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
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
	var input string
	if len(flag.Args()) > 0 {
		// Read from the file if a filename is provided as a command-line argument
		input = nms_read_file(flag.Arg(0))
	} else {
		// Otherwise, read from stdin
		input = nms_read_stdin()
	}

	// Split the input into lines and process each line
	lines := strings.Split(input, "\n")
	for y, line := range lines {
		// Process the line
		nms_chars := nms_process_input(line)

		// If random flag is set, shuffle the characters
		if *opt_random {
			rand.Shuffle(len(nms_chars), func(i, j int) { nms_chars[i], nms_chars[j] = nms_chars[j], nms_chars[i] })
		}

		// Output the processed line
		for x, ch := range nms_chars {
			// Scramble the character
			screen.SetContent(x, y, ch.scram, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			screen.Show()
			if !*opt_auto {
				delay(time.Duration(*opt_delay))
			}
		}

		// Reveal the original text
		for x, ch := range nms_chars {
			// Reveal the character
			nms_reveal(&ch)
			screen.SetContent(x, y, ch.scram, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			screen.Show()
			if !*opt_auto {
				delay(time.Duration(*opt_delay))
			}
		}
	}

	// Wait for a key press before exiting
	if !*opt_auto {
		screen.PollEvent()
	}
}

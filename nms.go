package main

import (
    "bufio"
    "flag"
    "fmt"
	"os"
    "io/ioutil"
    "math/rand"
    "time"
    "github.com/nsf/termbox-go"
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
    for i := ' '; i <= '~'; i++ {
        charset += string(i)
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

    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    var input string
    if len(flag.Args()) > 0 {
        // Read from the file if a filename is provided as a command-line argument
        input = nms_read_file(flag.Arg(0))
    } else {
        // Otherwise, read from stdin
        input = nms_read_stdin()
    }

    // Process the input
    nms_chars := nms_process_input(input)

    // Output the processed input
    for i, ch := range nms_chars {
        // Scramble the character
        termbox.SetCell(i, 0, ch.scram, termbox.ColorWhite, termbox.ColorBlack)
        termbox.Flush()
        delay(time.Duration(*opt_delay) * time.Millisecond)

        // Reveal the character
        nms_reveal(&ch)
        termbox.SetCell(i, 0, ch.scram, termbox.ColorWhite, termbox.ColorBlack)
        termbox.Flush()
        delay(time.Duration(*opt_delay) * time.Millisecond)
    }
}

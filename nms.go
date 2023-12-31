package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "time"
    "math/rand"

    "github.com/gdamore/tcell"
    "github.com/acarl005/stripansi"
)

var (
    opt_version = flag.Bool("version", false, "print version information")
    opt_delay   = flag.Int("delay", 5, "set delay in ms")
    opt_sneakers = flag.Bool("sneakers", false, "display the 'Sneakers' screen")

    charset string // String containing the printable ASCII characters for scrambling
)

type NmsChar struct {
    ch    rune
    scram rune
}

func init() {
    // Standard ASCII.
    for i := 33; i <= 126; i++ {
        charset += string(rune(i))
    }
    // Line-drawing characters from the extended ASCII.
    for _, i := range []int{179, 180, 191, 192, 193, 194, 195, 196, 197, 217, 218, 219} {
        charset += string(rune(i))
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

func nms_sneakers_screen() []string {
    return []string{
        "",
        "    DATANET PROC RECORD:  45-3456-W-3452                                                           Transnet on/xc-3",
        "",
        "                               FEDERAL RESERVE TRANSFER NODE",
        "",
        "                               National Headquarters",
        "",
        "       ************  Remote Systems Network Input Station  ************",
        "       ================================================================================",
        "",
        "       [1] Interbank Funds Transfer  (Code Prog: 485-GWU)",
        "       [2] International Telelink Access  (Code Lim: XRP-262)",
        "       [3] Remote Facsimile Send/Receive  (Code Tran:  2LZP-517)",
        "       [4] Regional Bank Interconnect  (Security Code:  47-B34)",
        "       [5] Update System Parameters  (Entry Auth. Req.)",
        "       [6] Remote Operator Logon/Logoff",
        "",
        "       ================================================================================",
        "",
        "       [ ] Select Option or ESC to Abort",
    }
}

func nms_process_input(input string) []NmsChar {
    input = stripansi.Strip(input)
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
        fmt.Println("nms v0.5.0")
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

    // Get the size of the terminal.
    width, height := screen.Size()

    var inputLines []string
    if *opt_sneakers {
        inputLines = nms_sneakers_screen()
    } else if len(flag.Args()) > 0 {
        inputLines = nms_read_file(flag.Arg(0))
    } else {
        inputLines = nms_read_stdin()
    }

    nms_lines := make([][]NmsChar, len(inputLines))
    for i, line := range inputLines {
        nms_lines[i] = nms_process_input(line)
    }

    type Pos struct{ x, y int }
    var allChars []Pos
    for y, nms_chars := range nms_lines {
        for x := range nms_chars {
            allChars = append(allChars, Pos{x, y})
        }
    }

    rand.Shuffle(len(allChars), func(i, j int) { allChars[i], allChars[j] = allChars[j], allChars[i] })

    eventCh := make(chan tcell.Event)
    go func() {
        for {
            eventCh <- screen.PollEvent()
        }
    }()

    // Calculate vertical padding based on terminal height and number of lines.
    vPad := (height - len(nms_lines)) / 2

    // Calculate horizontal padding based on terminal width and the maximum line length.
    maxLength := 0
    for _, nms_chars := range nms_lines {
        if len(nms_chars) > maxLength {
            maxLength = len(nms_chars)
        }
    }
    hPad := (width - maxLength) / 2

    // Adjust x to include padding when setting content.
    for y, nms_chars := range nms_lines {
        for x, ch := range nms_chars {
            screen.SetContent(x+hPad, y+vPad, ch.scram, nil, tcell.StyleDefault.Foreground(tcell.ColorLightBlue))
        }
        screen.Show()
        time.Sleep(time.Duration(*opt_delay) * time.Millisecond)
    }

    time.Sleep(2 * time.Second) // Pause before starting to reveal

    const scrambleIterations = 5 // number of times to scramble each character before revealing
    for iteration := 0; iteration < scrambleIterations; iteration++ {
        for _, pos := range allChars {
            nms_scramble(&nms_lines[pos.y][pos.x])
            screen.SetContent(pos.x+hPad, pos.y+vPad, nms_lines[pos.y][pos.x].scram, nil, tcell.StyleDefault.Foreground(tcell.ColorLightBlue))
            screen.Show()
        }
    }

    for _, pos := range allChars {
        nms_reveal(&nms_lines[pos.y][pos.x])
        screen.SetContent(pos.x+hPad, pos.y+vPad, nms_lines[pos.y][pos.x].scram, nil, tcell.StyleDefault.Foreground(tcell.ColorCornflowerBlue))
        screen.Show()

        select {
        case ev := <-eventCh:
            if _, ok := ev.(*tcell.EventKey); ok {
                return
            }
        default:
        }
    }

    time.Sleep(5 * time.Second) // Pause before exit
}

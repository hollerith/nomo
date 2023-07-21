![image](https://github.com/hollerith/nomo/assets/659626/52e9bb0c-ec42-49fb-a9f9-a7690580feae)

# No more secrets

NMS is a terminal-based text scrambling program written in Go. It reads in text, scrambles the characters, displays the scrambled text, then reveals the original text character-by-character in a decryption animation.

## Features

- Scrambles text by swapping each character with a random printable ASCII character
- Can read text from standard input, files, or display a hardcoded "Sneakers" screen
- Reveals text slowly in a decryption animation
- Customizable delay between reveal steps
- Runs in a terminal using github.com/gdamore/tcell for terminal control

## Usage

```
nms [options] [filename]

Options:
  -delay=<ms>   Set delay in ms between reveal steps (default 5)
  -sneakers     Display the 'Sneakers' screen
  -version      Print version information
```

If no filename is provided, NMS will read from standard input.

Use `-sneakers` to display the hardcoded "Sneakers" screen from the movie as an example.

## Building

To build NMS, you will need:

- Go 1.12 or higher
- [github.com/gdamore/tcell](https://github.com/gdamore/tcell)

Run:

```
go build nms.go
```

This will produce an executable called `nms`.

## Example

```
$ nms -delay 10 -sneakers
```

This will display the "Sneakers" screen and use a 10ms delay between reveal steps.

Press any key to exit when done.

## License

NMS is released under the MIT License. See LICENSE.txt.

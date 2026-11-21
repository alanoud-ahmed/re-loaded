# go-reloaded

A small text completion, editing and auto-correction tool written in Go, using only the standard library.

It reads a text file, fixes the text line by line, and writes the result to another file.

## Usage

```
go run . input.txt output.txt
```

- `input.txt`: the text to fix (must exist).
- `output.txt`: where the fixed text is written.

Both files must be `.txt` and must not be the same file.

## What it does

| Feature | Example |
|---|---|
| `(hex)` converts the word before it from hexadecimal to decimal | `1E (hex) files` -> `30 files` |
| `(bin)` converts the word before it from binary to decimal | `10 (bin) years` -> `2 years` |
| `(up)` / `(low)` / `(cap)` change the case of the word before | `go (up)` -> `GO` |
| `(up, N)` / `(low, N)` / `(cap, N)` apply to the previous N words | `so exciting (up, 2)` -> `SO EXCITING` |
| Punctuation `. , ! ? : ;` is attached to the previous word | `there ,and` -> `there, and` |
| Groups like `...` and `!?` are kept together | `thinking ... You` -> `thinking... You` |
| Single quotes are attached to the words inside them | `' awesome '` -> `'awesome'` |
| `a` becomes `an` before a vowel or `h` | `A amazing rock` -> `An amazing rock` |

## How it works

Each line goes through a simple pipeline:

1. **Normalize**: clean up spacing around commands.
2. **Validate**: check quotes and commands are well formed.
3. **Apply commands**: run `(hex)`, `(bin)`, `(up)`, `(low)`, `(cap)`.
4. **Format**: fix articles, punctuation, quotes and brackets.

If a line has a problem (for example an invalid command), an error message with the line number is written in place of that line, and the other lines are still processed.

## Example

Input:

```
Simply add 42 (hex) and 10 (bin) and you will see the result is 68.
```

Output:

```
Simply add 66 and 2 and you will see the result is 68.
```

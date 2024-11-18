<!-- omit in toc -->
# The Launcher

A fast configurable launcher

![screenshot](Images/screenshot_v0.2.png)

Status: in development (expect breaking changes)

## Build

For Windows without cgo (CGO_ENABLED=0), the raylib.dll v5.0 is included in this repository.

For other OS or Windows with cgo, check [raylib-go Requirements](https://github.com/gen2brain/raylib-go#Requirements).

Use this command to build for release

```shell
go build -ldflags "-H=windowsgui -w -s"
```

- `-H=windowsgui` remove console window, it greatly improves performances
- `-w -s` reduces final binary size by stripping debug symbols

The final executable should be shipped with

- Fonts directory containing used fonts
- config.toml
- raylib.dll (for Windows)

## Limitations

- The launcher can not start TUI programs (e.g.: vim)
- There can not be comments in the config.toml file

## Controls

**ToDo** : write this

## Configuration

**ToDo** : write documentation about config.toml syntax, and a few concrete examples.

### Example rules

#### Static rules

| Typed     | Description           | Command                           |
| --------- | --------------------- | --------------------------------- |
| Desktop   | Open Desktop folder   | explorer.exe <desktop_location>   |
| Documents | Open Documents folder | explorer.exe <documents_location> |
| SVN       | Open SVN folder       | explorer.exe C:\SVN               |
| py        | Start python script   | pythonw.exe python_script.pyw     |

#### Dynamic rules

These rules will use regular expressions. Not implemented yet.

| Typed            | Description                   | Command                                                         |
| ---------------- | ----------------------------- | --------------------------------------------------------------- |
| py {arg}         | Start python script with args | python_script.pyw {arg}                                         |
| r/{sub}          | Go to r/{sub}                 | firefox.exe <https://www.reddit.com/r/{sub}/>                   |
| r/{sub} {search} | Search on r/{sub}             | firefox.exe <https://www.reddit.com/r/{sub}/search/?q={search}> |
| r {search}       | Search on Reddit              | firefox.exe <https://www.reddit.com/search/?q={search}>         |

## Resources

- Font Cascadia code : <https://github.com/microsoft/cascadia-code>
- Raylib DLL : <https://github.com/raysan5/raylib/releases/tag/5.0>

## ToDo list

List of ideas to implement in no particular order

- Config: Add configuration to enable search by rule description
- Config/GUI: Replace some GUI constants with configuration (font, font size, width, colours, ...)
- GUI: Improve selection colours and controls
- Rule: Manage environment variables in rule Exe/Args
- GUI: Ctrl + Z / Ctrl + Shift + Z to undo/redo
- Rule: Add regexp management in rules
- GUI: Ctrl + V to paste text (use rl.GetClipboardText)
- Misc: Comment the code some more
- Commands: Add standard commands
- Commands: Add /config - edit configuration file
  - Add setting for default editor
- Commands: Add /reset - reset all LastUse values
- Misc: Refactor GUI_Start function (Create a GUI class with methods) ?
- Misc: If config.toml is not found, create it with a few examples dummy rules

 # ttymer

 A terminal-based task timer and task list built with Bubble Tea. Manage a simple list of tasks, edit details, and run a combined timer with a progress bar.

 ## Features
 - Task list with name, description, and duration
 - Edit task details (name, description, duration)
 - Timer view with progress bar and pause/resume
 - Keyboard-driven TUI

 ## Requirements
 - Go 1.25.6+ (see `go.mod`)

 ## Run
 ```bash
 go run ./cmd
 ```

 ## Build
 ```bash
 go build -o ttymer ./cmd
 ```

 ## Controls
 ### Task list view
 - `e` edit selected task
 - `t` go to timer view
 - `q` quit

 ### Edit task view
 - `tab` / `shift+tab` move between fields
 - `enter` advance / save on duration field
 - `esc` cancel

 ### Timer view
 - `p` pause/resume

 ## Notes
 - The app currently starts with a couple of example tasks defined in `InitTui`.

## Todo
- Remove examples and add SQLite db for storing tasks
- Add breaks into timer
- Refactor timer code
- Write tests

 ## Project structure
 - `cmd/main.go`: app entrypoint
 - `task/`: task model
 - `tui/`: Bubble Tea models and views

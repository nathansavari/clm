# Command Line Manager (CLM)

The Command Line Manager (clm) tool lets you create, list, and delete custom command templates with variable substitution. It requires no additional installation steps or dependencies on the user's machine (just [Go](https://go.dev/)).

## Install

```bash
go install github.com/nathansavari/clm@latest
```

This installs a go binary that will automatically bind to your $GOPATH

## Commands

- Add a Command: clm new
- List Commands: clm list [tag]
- Delete a Command: clm delete

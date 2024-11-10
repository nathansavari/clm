package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/manifoldco/promptui"
)

type Command struct {
	Title   string `json:"title"`
	Command string `json:"command"`
	Tag     string `json:"tag"`
}

const commandsFile = "commands.json"

// Add a new command
func addCommand() {
	var cmd Command
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter command title: ")
	title, _ := reader.ReadString('\n')
	cmd.Title = strings.TrimSpace(title)

	fmt.Print("Enter command: ")
	command, _ := reader.ReadString('\n')
	cmd.Command = strings.TrimSpace(command)

	fmt.Print("Enter unique tag: ")
	tag, _ := reader.ReadString('\n')
	cmd.Tag = strings.TrimSpace(tag)

	commands := readCommands()
	commands = append(commands, cmd)
	writeCommands(commands)

	fmt.Println("Command saved successfully!")
}

// List and select commands, with optional filtering by tag
func listCommands(tag string) {
	commands := readCommands()
	if len(commands) == 0 {
		fmt.Println("No commands found.")
		return
	}

	var filtered []Command
	if tag == "" {
		filtered = commands
	} else {
		for _, cmd := range commands {
			if cmd.Tag == tag {
				filtered = append(filtered, cmd)
			}
		}
		if len(filtered) == 0 {
			fmt.Printf("No commands found with tag '%s'.\n", tag)
			return
		}
	}

	// Use promptui for navigation
	items := make([]string, len(filtered))
	for i, cmd := range filtered {
		items[i] = fmt.Sprintf("%s: %s", cmd.Title, cmd.Command)
	}

	prompt := promptui.Select{
		Label: "Select a command to copy",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Command selection cancelled.")
		return
	}

	// Get the selected command
	var selectedCommand string
	for _, cmd := range filtered {
		if strings.Contains(result, cmd.Command) {
			selectedCommand = cmd.Command
			break
		}
	}

	// Process variables in the command
	finalCommand := processCommandVariables(selectedCommand)

	clipboard.WriteAll(finalCommand)

	fmt.Printf("\nCommand copied to your clipboard:\n\n%s\n\n", finalCommand)
}

// Delete a command by selecting from the list
func deleteCommand() {
	commands := readCommands()
	if len(commands) == 0 {
		fmt.Println("No commands found to delete.")
		return
	}

	// Display the list of commands for selection
	items := make([]string, len(commands))
	for i, cmd := range commands {
		items[i] = fmt.Sprintf("%s: %s", cmd.Title, cmd.Command)
	}

	prompt := promptui.Select{
		Label: "Select a command to delete",
		Items: items,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Command deletion cancelled.")
		return
	}

	// Remove the selected command
	commands = append(commands[:index], commands[index+1:]...)
	writeCommands(commands)

	fmt.Println("Command deleted successfully.")
}

// Process command for variable substitution
func processCommandVariables(command string) string {
	varRegex := regexp.MustCompile(`<(\w+)>`)
	matches := varRegex.FindAllStringSubmatch(command, -1)

	reader := bufio.NewReader(os.Stdin)

	for _, match := range matches {
		varName := match[1]
		fmt.Printf("Enter value for %s: ", varName)
		value, _ := reader.ReadString('\n')
		command = strings.Replace(command, match[0], strings.TrimSpace(value), -1)
	}

	return command
}

// Helper functions to read and write commands to JSON
func readCommands() []Command {
	file, err := os.Open(commandsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Command{}
		}
		fmt.Println("Error reading file:", err)
		return nil
	}
	defer file.Close()

	var commands []Command
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&commands)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	return commands
}

func writeCommands(commands []Command) {
	file, err := os.Create(commandsFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(commands)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: clm <new|list|delete> [tag]")
		return
	}

	switch os.Args[1] {
	case "new":
		addCommand()
	case "list":
		tag := ""
		if len(os.Args) > 2 {
			tag = os.Args[2]
		}
		listCommands(tag)
	case "delete":
		deleteCommand()
	default:
		fmt.Println("Unknown command. Use 'new' to add a command, 'list' to list commands, or 'delete' to delete a command.")
	}
}

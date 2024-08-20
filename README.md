# Git Prompt

A simple Git prompt integration tool written in Go.

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/git-prompt.git
   ```
2. Build the project:
   ```
   cd git-prompt
   go build
   ```
3. Move the built binary to a directory in your PATH.

## Usage

Run the `git-prompt` command in your terminal. It will provide an interactive prompt with Git information.

### Customization

You can customize the prompt by modifying the `defaultConfig` in the `main.go` file. The available options are:

- `ShowBranch`: Display the current Git branch (default: true)
- `ShowModifiedFiles`: Display the number of modified files (default: false)
- `ShowAheadBehind`: Display how many commits ahead or behind the remote branch (default: false)
- `ShowStashCount`: Display the number of stashed changes (default: false)

After modifying the configuration, rebuild the project for the changes to take effect.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
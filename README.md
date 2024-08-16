# Git Prompt

A simple Git prompt integration tool written in Go.

## Installation

1. Clone the repository
2. Build the project using `go build`
3. Run the executable

## Usage

1. Run the executable
2. The prompt will display Git information for the current directory
3. Customize the prompt by editing the `config.json` file

### Customization

You can customize the prompt format and colors by editing the `config.json` file. The file contains two main sections:

1. `promptFormat`: This string defines the format of the prompt. You can rearrange or remove elements as needed.
2. `theme`: This object contains color definitions for different parts of the prompt.

Example `config.json`:

```json
{
  "promptFormat": "%s(%s%s%s|%s%d%s|%s%s%s|%s%d%s)%s %s $ ",
  "theme": {
    "resetColor": "\u001b[0m",
    "grayColor": "\u001b[38;5;243m",
    "lightBlue": "\u001b[38;5;117m",
    "lightGreen": "\u001b[38;5;114m",
    "lightRed": "\u001b[38;5;174m",
    "lightYellow": "\u001b[38;5;186m",
    "lightCyan": "\u001b[38;5;152m"
  }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
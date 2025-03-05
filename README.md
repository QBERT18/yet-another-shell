# 🐚 **Build Your Own Shell** 🐚

[![progress-banner](https://backend.codecrafters.io/progress/shell/514277f4-e764-48fa-b6e2-3899cbb05cf9)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This repository contains my solution for the **["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview)** on CodeCrafters.

---

## 🚀 **About the Project**

I built my own shell (and PowerShell-like implementation) using **Go**! 🛠️  
The project includes a custom **lexer** for parsing user input, inspired by this guide:  
[**Handwritten Parsers & Lexers in Go**](https://blog.gopheracademy.com/advent-2014/parsers-lexers/).

The shell is designed to work on both **Linux** 🐧 and **Windows** 🪟, though it’s still a work in progress. Some commands may not work perfectly yet, but I’m actively improving it! 💪

---

## 🛠️ **Features**

- **Custom Lexer**: Handles tokenization of user input.
- **Cross-Platform**: Works on both Linux and Windows.
- **Basic Commands**: Supports commands like `cd`, `echo`, `ls`, and more.
- **PowerShell-like Commands**: Includes equivalents for PowerShell commands (e.g., `Write-Output` for `echo`).

---

## 🐞 **Known Issues**

- Some commands are still buggy or incomplete.
- Cross-platform compatibility needs further testing.
- Error handling could be improved.

---

## 🧪 **Testing**

To test the shell, run the following commands:

### **Basic Commands**
```bash
echo "Hello, World!"
cd /path/to/directory
ls
```

### **PowerShell-like Commands**
```powershell
Write-Output "Hello, World!"
Set-Location C:\path\to\directory
Get-ChildItem
```

---

## 🛠️ **How to Run**

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/your-shell-repo.git
   ```
2. Navigate to the project directory:
   ```bash
   cd your-shell-repo
   ```
3. Build and run the shell:
   ```bash
   go run main.go
   ```

---

## 🌟 **Inspiration**

This project was inspired by the **["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview)** on CodeCrafters. If you’re viewing this repo on GitHub, head over to [codecrafters.io](https://codecrafters.io) to try the challenge yourself!

---

## 🙏 **Acknowledgments**

- [CodeCrafters](https://codecrafters.io) for the amazing challenge.
- [Gopher Academy](https://blog.gopheracademy.com) for the lexer and parser guide.

---

## 🚧 **Work in Progress**

This project is still under development. Contributions and feedback are welcome! Feel free to open an issue or submit a pull request.

---

Enjoy building your own shell! 🎉  
Happy coding! 💻
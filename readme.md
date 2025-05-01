# **guoblock (Guarding Unwanted Omissions)**

guoblock is a CLI tool that helps to identify and block the accidental inclusion of sensitive secrets (such as API keys, AWS credentials, and other private tokens) in your codebase. It's designed for use in Git repositories or local project directories.

---

## **Features**

- **Secret scanning**: Detects secrets like AWS keys, Slack tokens, API keys, and more.
- **Configurable**: Add custom rules through a simple YAML file.
- **Ignore paths**: Skip specific directories (e.g., `vendor/`, `.git/`).
- **Exit codes**: Returns `0` for no secrets found, `1` for secrets found, and `2` for errors.

---

## **Installation**

### Build from Source:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/guoblock.git
   cd guoblock
   ```

2. Build the binary:
   ```bash
   go build -o guoblock main.go
   ```

3. Optionally, install the tool globally:
   ```bash
   sudo mv guoblock /usr/local/bin/
   ```

---

## **Usage**

### Scan a Directory:
```bash
guoblock ./path/to/scan
```

---

## **Exit Codes**

- `0` — No secrets found.
- `1` — Secrets detected in the codebase.
- `2` — Error in running the scan (e.g., file not found, permissions issue).

---

## **Development**

To contribute to **GuoBlock**, follow these steps:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-name`).
3. Write tests and code.
4. Commit your changes (`git commit -am 'Add feature'`).
5. Push to your fork (`git push origin feature-name`).
6. Open a pull request.
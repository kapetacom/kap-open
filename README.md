# kap-open: Kapeta URL Authentication Tool

kap-open is a developer tool designed for opening authenticated Kapeta URLs. It's a command-line utility written in Go that facilitates the seamless access of Kapeta URLs while ensuring proper authentication.

## Prerequisites

Before using kap-open, make sure you have the following prerequisites installed:

- Go programming language (https://golang.org/dl/)

## Installation

### RPM Package

```shell
sudo tee -a /etc/yum.repos.d/artifact-registry.repo << EOF
[kapeta-production-yum]
name=kapeta-production-yum
baseurl=https://europe-north1-yum.pkg.dev/projects/kapeta-production/kapeta-production-yum
enabled=1
repo_gpgcheck=0
gpgcheck=0
EOF
```

Update the package list and install kap-open:

```shell
sudo dnf update
sudo dnf install kap-open
```

## Usage

1. Install kap-open using the provided installation instructions.
2. Open a terminal and run kap-open with a port as an argument:

   ```shell
   kap-open <port of the local service you would like to connect to like 5009>
   ```

   kap-open will automatically handle the authentication process and open the URL in your default web browser.

## Limitations

- It's important to note that kap-open is specifically designed for opening Kapeta URLs and doesn't handle general URL opening functionality.

## License

kap-open is released under the [MIT License](LICENSE). Feel free to contribute and enhance the tool as needed.
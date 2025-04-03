# CLI Port Scanner

A simple CLI-based port scanner written in Go. This tool scans specified targets for open ports and optionally retrieves service banners. It is designed for educational purposes and to help users understand the basics of network scanning.

> [!TIP]
>
> Watch a demo of the CLI Port Scanner in action on [YouTube](https://www.github.com/jennxiserra).

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Flags](#flags)
- [Makefile Uses](#makefile-uses)
- [Additional Notes](#additional-notes)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Features

- **Concurrent Scanning:** Leverages goroutines with a configurable worker pool.
- **Custom Port Ranges:** Scan a range of ports or specific ports by overriding defaults.
- **Service Banner Grab:** Attempts to retrieve banners (e.g., HTTP "Server" header) from open ports.
- **JSON Output:** Option to output scan results in JSON format.
- **Simple CLI Interface:** Easily configurable via command-line flags.

## Installation

> [!NOTE]
>
> Ensure you have [Go](https://golang.org/dl/) (version 1.24.1 or later) installed before proceeding.

Follow these steps to install and build the project:

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/jennxsierra/cli-port-scanner.git
    cd cli-port-scanner
    ```

2. **Build the Application:**

    - **Using Make (recommended):**  
      Run the following command to build the binary:

      ```bash
      make build
      ```

    - **Without Make:**  
      If Make is not installed, you can build manually:

      ```bash
      go fmt ./...
      go vet ./...
      go build -o bin/cli-pscan main.go
      ```

The binary will be located in the `bin/` directory.

## Usage

Below is an example command that demonstrates the use of multiple flags. This example scans `example.com` for ports 80-85, uses 50 workers, sets a timeout of 3 seconds, enables debug output, and writes JSON scan results to a file.

```bash
./bin/cli-pscan -target=example.com -start-port=80 -end-port=85 -workers=50 -timeout=3 -json
```

**Terminal Output:**

```bash
====================================================
               SCAN START: example.com
====================================================
Scanning ports: 80, 81, 82, 83, 84, 85
Progress: [████████████████████] 100% [6/6 ports]

[SCAN SUMMARY]
Target         : example.com
Ports Scanned  : 80, 81, 82, 83, 84, 85
Open Ports     : 80
Time Taken     : 16.005s

[BANNERS]
Port 80   : "AkamaiGHost"

Scan results saved to scan-results/030425-160155-cli-pscan.json
```

**JSON File:**

The JSON output will be written in the designated output directory (`scan-results/`) and will be named using the timestamp format `DDMMYY-HHMMSS-cli-pscan.json`. For example:

```json
{
  "timestamp": "2025-04-03T22:01:55Z",
  "targets": [
    "example.com"
  ],
  "total_ports": 6,
  "results": [
    {
      "target": "example.com",
      "open_ports": [
        {
          "port": 80,
          "banner": "AkamaiGHost"
        }
      ],
      "total_ports": 6,
      "open_count": 1,
      "duration": "16.005s"
    }
  ]
}
```

## Flags

| Flag          | Type   | Default | Description                                                                 |
| ------------- | ------ | ------- | --------------------------------------------------------------------------- |
| `-target`     | string | -       | The hostname or IP address to be scanned.                                   |
| `-targets`    | string | -       | Comma-separated list of targets (e.g., localhost,scanme.nmap.org).          |
| `-start-port` | int    | 1       | The lower bound port to begin scanning.                                     |
| `-end-port`   | int    | 1024    | The upper bound port to finish scanning.                                    |
| `-ports`      | string | -       | Comma-separated list of specific ports (overrides start-port and end-port).   |
| `-workers`    | int    | 100     | Number of concurrent goroutines to launch per target.                       |
| `-timeout`    | int    | 5       | Connection timeout in seconds.                                              |
| `-json`       | bool   | false   | Output results in JSON format.                                              |
| `-debug`      | bool   | false   | Display flag values for debugging.                                          |

You can view all available flags and their descriptions by running:

```bash
./bin/cli-pscan -help
```

## Makefile Uses

The provided `Makefile` streamlines several tasks:

- **fmt:** Formats the source code.
- **vet:** Runs static analysis.
- **build:** Builds the application and places the binary into the `bin/` directory.
- **test:** Runs unit tests with verbose output.
- **clean:** Cleans up build artifacts and JSON output files.
- **check:** Performs format, vet, and tests in one step.

Run any of these targets by using `make <target>`. For example:

```bash
make build
```

## Additional Notes

- **Retry Policy:**  
  The scanner attempts a fixed number of retries per port (currently set to 3). Increasing this value may help capture intermittent connection issues; however, it will usually result in a slower overall scan. Each retry uses an exponential backoff delay (e.g., 1 second for the first retry, 2 seconds for the second, etc.).

- **Port Selection:**  
  If the `-ports` flag is provided, the values specified with `-start-port` and `-end-port` are ignored.

- **Target Aggregation:**  
  Specifying both `-target` and `-targets` merges all provided hosts for scanning. If no target is specified, the scanner defaults to `localhost`.

- **Timeout Handling:**  
  The `-timeout` flag applies both to establishing TCP connections and to banner grabbing operations. If banner outputs are inconsistent or missing, consider increasing the timeout value.

- **Uniform Configuration for Multiple Targets:**  
  When scanning more than one target, all flags (except for those specific to target selection) are applied uniformly across every target.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

- This project was developed as part of an assignment for the **[CMPS2242] Systems Programming & Computer Organization** course  under the Associate of Information Technology program at the [University of Belize](https://www.ub.edu.bz/).
- Special thanks to Mr. Dalwin Lewis for his guidance and support.

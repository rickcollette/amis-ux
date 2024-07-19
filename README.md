
# AMIS/UX

AMIS/UX is a modern BBS (Bulletin Board System) application that brings back the nostalgia of traditional dial-up BBS systems, while incorporating contemporary technologies and interfaces. The application supports multiple text-based modes, including ASCII, ATASCII, and ANSI, to offer a versatile and engaging user experience.  

## Introduction

AMIS/UX is designed to recreate the classic BBS experience with a modern twist. Users can connect to the BBS, register, post messages, and interact with message bases. The application is built with Go and leverages SQLite for its database, making it lightweight and easy to deploy. It also includes a setup utility to help configure the BBS and manage message bases.

It was built by matching my original AMIS BBS modified source code and features to modern programming language and models.

## Installation

### Prerequisites (to build it from source)

- Go (version 1.16 or later)
- SQLite3
- A Unix-like operating system (Linux, macOS)

### Steps

1. **Clone the repository:**

    ```sh
    git clone https://github.com/yourusername/amis-ux.git
    cd amis-ux
    ```

2. **Build the BBS application:**

    ```sh
    go build -o amisbbs ./bbs
    ```

3. **Build the setup utility:**

    ```sh
    go build -o setup ./setup
    ```

4. **Prepare the database:**

    The first time you run the setup utility, it will create the necessary tables in the SQLite database:

    ```sh
    ./setup
    ```

5. **Configure the BBS:**

    Follow the prompts in the setup utility to configure the BBS. The configuration will be saved to `config.json`.

## Usage

### Running the BBS

Start the BBS server:

```sh
./amisbbs
```

The server will start listening on the port specified in the configuration (`config.json`). Users can connect to the BBS via Telnet.

### Interacting with the BBS

- **Registering a new user:**
  Users can register by providing a username and password during their first connection.

- **Posting messages:**
  Users can post messages to different message bases after logging in.

- **Viewing messages:**
  Users can view messages in any message base they have access to.

### Switching Modes

AMIS/UX supports multiple display modes:

- **ASCII**
- **ATASCII**
- **ANSI**

Users can toggle between these modes using the provided commands within the BBS.

## Configuration

The configuration for AMIS/UX is stored in `config.json`. This file is generated and managed by the setup utility but can also be manually edited if needed. Below is an example of the configuration file:

```json
{
  "bbs_name": "My BBS",
  "sysop_name": "John Wick",
  "allow_new_users": false,
  "ascii_folder": "/app/running-app/ascii",
  "atascii_folder": "/app/running-app/atascii",
  "ansi_folder": "/app/running-app/ansi",
  "menus_folder": "/app/running-app/menus",
  "executables_folder": "/app/running-app/command",
  "sysop_password": "$2a$10$Mc.awholebunchofstuffherebecauseitshashed",
  "port_number": 8023
}
```

### Configuration Options

- `bbs_name`: The name of your BBS.
- `sysop_name`: The name of the system operator.
- `allow_new_users`: Boolean value to allow or disallow new user registrations.
- `ascii_folder`: Path to the ASCII files directory.
- `atascii_folder`: Path to the ATASCII files directory.
- `ansi_folder`: Path to the ANSI files directory.
- `menus_folder`: Path to the menus directory.
- `executables_folder`: Path to the executables directory.
- `sysop_password`: The hashed password for the sysop.
- `port_number`: The port number on which the BBS will listen for incoming connections.

## License

AMIS/UX is licensed under the MIT License. 

Copyright 2024 Rick Collette megalith _at_ root.sh 

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

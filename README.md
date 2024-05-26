# SIP Server and Client

A comprehensive SIP (Session Initiation Protocol) server and client implementation using pjsua2.

## About

This project provides a SIP server and client solution, leveraging the pjsua2 library. It is designed to handle SIP signaling for VoIP (Voice over IP), instant messaging, and presence information. This implementation can be used for various SIP-based communication applications.

## Features

- **SIP Protocol Support**: Implements core SIP functionalities.
- **Client and Server Modes**: Can operate as both SIP server and client.
- **PJSIP Integration**: Uses the pjsua2 library for SIP signaling.
- **Multi-Platform**: Compatible with various platforms.
- **Modular Design**: Easy to extend and customize.

## Installation

### Prerequisites

- **Development Tools**: 
`go`

### Steps

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/tahseenjamal/sip-server-client.git
    cd sip-server-client
    ```

2. **Build the Client**:
    ```bash
    cd sip-client
    go get
    go build main.go
    ```


3. **Build the Server**:
    ```bash
    cd sip-server
    go get
    go build main.go
    ```

## Usage

### Running the server

1. Navigate to the client directory:
    ```bash
    cd sip-server
    ```

2. Run the client application:
    ```bash
    ./main
    ```

### Running the Client

1. Navigate to the client directory:
    ```bash
    cd sip-client
    ```

2. Run the client application:
    ```bash
    ./main
    ```

## Configuration

The configuration files for the server and client are located in their respective directories. You can modify these files to change the SIP settings, such as the server address, port, and authentication details.

## Contributing

We welcome contributions from the community. To contribute:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push the branch to your fork.
4. Open a pull request with a detailed description of your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

Special thanks to the contributors and the community for their support.



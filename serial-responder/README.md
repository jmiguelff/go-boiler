# Serial Responder

A simple Go application that reads from a serial port and automatically responds with "WORKING" when it receives a `!` character.

## Features

- Reads configuration from a YAML file (`settings.yml`)
- Monitors a serial port for incoming data
- Automatically replies with "WORKING" when receiving a `!` character
- Displays all received characters on stdout
- Logs responses sent back to the serial port

## Configuration

Edit `settings.yml` to configure your serial port settings:

```yaml
serialConf:
  device: /dev/cu.usbmodem1423401s # Serial port device
  dataBits: 8 # Number of data bits
  baud: 115200 # Baud rate
  stopbits: 1 # Number of stop bits
  parity: N # Parity: N (None), E (Even), O (Odd)
  timeout: 50 # Read timeout in milliseconds
```

## Building

```bash
go build -o serial-responder main.go
```

## Running

```bash
./serial-responder
```

Make sure `settings.yml` is in the same directory as the executable.

## Dependencies

- `github.com/tarm/serial` - Serial port library
- `gopkg.in/yaml.v2` - YAML configuration parsing

## Usage Example

1. Configure your serial port in `settings.yml`
2. Run the application
3. Send `!` through the serial port
4. The application will respond with "WORKING"

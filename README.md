# Smart Thermometer Server Application

## Setup

- Enable your [Google Sheets API](https://developers.google.com/sheets/api/quickstart/go) and download the `credentials.json` to the working directory
- Generate SSL keys using
```
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.crt -days 365
```
- Copy the [example sheet](https://docs.google.com/spreadsheets/d/1KcoxTs_B7jM9KdlDXLXEX2OaK6bBWlaWs5dJw_GVGto/) to your account and edit it for your needs
- Rename `config.example.json` to `config.json` and edit it:
  - `Key`: Pre-shared key. Requests to the server must contain the same key
  - `SheetsID`: ID of the spreadsheet
  - `TimeZone`: Specific time zone with IANA timezone database format
  - `Noon`: Hours (1~24) after this value will be considered afternoon


## Dependencies
- google.golang.org/api/sheets/v4
- golang.org/x/oauth2/google
- github.com/mattn/go-sqlite3

Use `go mod download` to install them

## Build
```
go build
```

## Run

### Just run
```
./smart-therometer
```

### Systemd service
```
# /etc/systemd/system/smart-thermometer.service

[Unit]
Description=Smart Thermometer

[Service]
WorkingDirectory=/home/yi/SmartThermometer-Server
ExecStart=/home/yi/SmartThermometer-Server/smart-thermometer

[Install]
WantedBy=multi-user.target
```

And then
```
sudo systemd enable --now smart-thermometer
```
to start and run on boot


# Smart Thermometer Server Application

## Setup

- Enable your [Google Sheets API](https://developers.google.com/sheets/api/quickstart/go) and download the `credentials.json` to the working directory
- Generate SSL keys using
```
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.crt -days 365
```
- Copy the [example sheet](https://docs.google.com/spreadsheets/d/1KcoxTs_B7jM9KdlDXLXEX2OaK6bBWlaWs5dJw_GVGto/) to your account and edit it for your needs
- Edit your `config.json` and enter your sheet ID

## Dependencies
```
go get -u google.golang.org/api/sheets/v4
go get -u golang.org/x/oauth2/google
```

## Run
```
go run .
```

## Build
```
go build
```
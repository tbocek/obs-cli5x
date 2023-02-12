# Simple OBS 5.x websocket command line client

This tool provides the minimal functionality of obs-cli, which works 
for the OBS Websocket 4.x protocol, but not for the 5.x protocol. This 
is a minimal implementation of the functionality I need.

Available functions:

* scene switching
* start/stop/toggle recording
* toggle scene items

The command line has the following usage:

```
Usage of obs-cli5x:
  -host string
        Host to connect to (default "localhost")
  -item string
        Toggle scene item
  -password string
        Password
  -port string
        Port to connect to (default "4455")
  -rec
        Toggle recording
  -scene string
        Set/change to scene
  -start-rec
        Start recording
  -stop-rec
        Stop recording

```

## Installation

Install go 1.20, then run

```
go build
go install
```

The binary will be installed in your go/bin directory somewhere in your home. Make sure its in the path.
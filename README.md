# priscilla-debug-console

[![Build Status](https://travis-ci.org/priscillachat/priscilla-debug-console.svg?branch=master)](https://travis-ci.org/priscillachat/priscilla-debug-console)

debug console for Priscilla chat bot server

## Usage

### Command line options

* ```-server <ip>``` server ip address, default 127.0.0.1
* ```-port <port>``` server port, default 4517
* ```-mode <mode>``` either responder or adapter, default responder
* ```-secret <secret>``` necessary for engagement with server, default abcdefghi

### Debug consule

Right now only two commands are supported:

* **put** activates transmission buffer mode, everything typed after put command
  is activated will be transmitted to the server. Use a ctrl-d on a new line to
  terminate the input mode and send current buffer.
* **exit** exit means exit

# iot
Beautifully automate your home with ESP8266

This project aims to get your home automated for the lowest price,
while providing a secure environment and a beautiful interface.

Demo moving blinds up and down
![](https://user-images.githubusercontent.com/4309591/31030903-40ab4444-a557-11e7-8f3c-9efa34cf37f4.gif)

# Security
* OAuth2.0 and Sign in with Google to authenticate users
* Websockets over https for secure connections to the backend
* https for the front end

*A bit more in detail*, he ESP8266 boots and connects over https to the backend, and sends an OWNER message, telling it the email of the owner of the device.
When the owner signs in to the backend, it will see a list of devices that have announced to be owned by their email address.

# Features
- [x] Control any ESP8266 securely from anywhere in the world
- [x] DigitalWrite to any pin using a switch
- [x] Dynamically configure actions
- [ ] Set timers to perform automated actions
- [ ] Share your devices with guests or friends
- [ ] Advanced widgets, 3-state and n-state switches
- [ ] Analog read/write using sliders
- [ ] Make buttons match pin state

# To-do
- [x] Establishing a secure backend-esp connection
- [x] Authenticate users securely using OAuth and Google sign in
- [x] Serving the front end in a secure way
- [x] Sending commands to the arduino
- [ ] Make a more consistent API
- [ ] Refactor the front end into plugins, make it use the API
- [ ] Documentation and testing (this is something for later, as currently the project is in early prototype stage)


# Devices
An ESP8266 device is identified by it's owner's email address and id. The id MUST be unique per user, but can collide with other users' devices.
A Command is an action to be applied to one or more pins. The ESP runs a small command interpreter. They're in the following format:
`ID <cmd> <pins>: PARAMS`


# Requirements, installing, setting up, running and developing

`git clone https://github.com/iotlib/iot`

### For the esp
* Install [PlatformIO](http://platformio.org/)
* Edit the [config file]()
* Open the project in atom
* Build and flash the firmware onto the ESP chip


### For the backend
* Edit the config file
* go run main.go
* Probably use a daemon script or something (TODO)

### For the frontend (development)
* Open in your favorite editor and `gulp dev`



# Contributing
Pull requests and issues are welcome!

# License
```
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```


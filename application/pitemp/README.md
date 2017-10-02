# Temperature Monitoring FLOGO Application for Raspberry PI

This application monitors the temperature using a DS18B20 temperature sensor connected to a Raspberry Pi through GPIO. 
* When the temperature goes above some defined threshold (defined as condition on a flow connector), the application lights on a LED (connected to GPIO 23) and sends a notification through PushOver Service
* When the temperature goes back below the threshold, the appplication lights off the LED and sends a clear notification through PushOver Service

## Important Note
The CLI flow may not work as it has been checked recently and some changes happened in Flogo
Flogo web flow has been updated and now works

## Pre-requisites

### Hardware
* One Raspberry Pi (I used version 3, but should work on other versions including Zero)
* One DS18B20 temperature sensor (I personally used the SunFounder module)
* One LED
* One resistor (from 270 to 330 Ohms)
* One experimental electronic board
* A few Wire Jumper Cables 

### Electronic Cabling

![Schema](https://github.com/ecarlier-tibco/flogo/blob/master/application/pitemp/electronic_wiring.png)

### Sofware
* The application uses WebIOPI to communicate with the GPIO I/O:
  * Install WebIOPI from http://webiopi.trouch.com/DOWNLOADS.html on your Raspberry Pi
  * Configure WebIOPI. I posted my config file webiopi.config to be copied as /etc/webiopi/config on your Pi
  * Start WebIOPI with the following command:
```
sudo webiopi -d -c /etc/webiopi/config
```
## Configuring, Compiling and Running the application:
* the cli sub-directory contains files and information  to setup and run the application with flogo-cli **(requires updates to work!)**
* the web sub-directory contains files and information  to setup and run the application with flogo-web

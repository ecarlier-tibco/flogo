# Install and Run pitemp app with flogo-cli

## Pre-requisites

Install flogo-cli on your development machine (could be your Raspberry Pi or your own desktop) as described [here](https://github.com/TIBCOSoftware/flogo-cli)

## Configuring the app

The application flow configuration is in ``pitemp/cli/pitemp_flow.json file``
Especially:
* You may have to change hostname or ip address to access the WebIOPI service on your Pi (all REST activities)
* You have to set the user key and autorization token for your PushOver account settings (both PushOver activities)
* You may want to change the temperature threshold (currently set to 23°C) on which to raise an alert (present in 2 flow links) 

## Creating & Building the app
```
$ flogo create pitemp
$ cd pitemp
$ flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/timer
$ flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
$ flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/rest
$ flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/counter
$ flogo add activity github.com/ecarlier-tibco/flogo/activity/pushover
$ flogo add flow /your/cloned/repo/path/application/pitemp/cli/pitemp_flow.json
$ cp /your/cloned/repo/path/application/pitemp/cli/triggers.json bin/
$ /your/cloned/repo/path/application/pitemp/cli/build_rpi.sh
```
The last command shall generate a binary ``pitemp-linux-arm`` in ``./bin sub-directory``
If your development is not the Pi, copy this file to the Pi and run it.
That's all Folks !

_**To increase temperature and raise alert, you can just put your fingers on the temperature sensor and wait a bit**_

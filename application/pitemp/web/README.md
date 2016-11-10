# Install and Run pitemp app with flogo-web

## Pre-requisites

Install and Run flogo-web on your development machine (could be your Raspberry Pi or your own desktop) 

## Importing the application in flogo-web
* Access flogo-web @ http://localhost:3010
* Click **Import a Flow** and select json file provided in this same directory ``pitemp_web_flow.json``
* Import the PushOver activity
	* Click + Button in a flow
	* On appearing right pane, click **Install new Activity**
	* Go to bottom of popped up window, enter ``https://github.com/ecarlier-tibco/flogo/activity/pushover`` in _Install activity from URL_ text box, and click on right-side download icon
	
## Configuring the app

Editing the flow in Flogo Web, you may want to adapt flow and configuration to your own setup.
Especially:
* You may have to change hostname or ip address to access the WebIOPI service on your Pi (all REST activities)
* You have to set the user key and autorization token for your PushOver account settings (both PushOver activities)
* You may want to change the temperature threshold (currently set to 23°C) on which to raise an alert (present in 2 flow links) 

## Building the app
* When editing the flow, mouse over the 3-dots button on top right of the window
* Select _Build--ARM/Linux_ in the appearing menu
* this will generate and download a ``build`` file, which your generated binary 

If your development is not the Pi, copy this build file to the Pi, possibly change the name, change the rights (chmod +x) and run it
That's all Folks !

_**To increase temperature and raise alert, you can just put your fingers on the temperature sensor and wait a bit**_

# pushover
This activity provides your flogo application the ability to send a notification via PushOver.


## Installation

```bash
flogo add activity github.com/ecarlier-tibco/flogo/activity/pushover
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "userKey",
      "type": "string"
    },
    {
      "name": "authToken",
      "type": "string"
    },
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "messageTitle",
      "type": "string"
    },
    {
      "name": "device",
      "type": "string"
    },
    {
      "name": "sound",
      "type": "string"
    },
    {
      "name": "url",
      "type": "string"
    },
    {
      "name": "urlTitle",
      "type": "string"
    },
    {
      "name": "priority",
      "type": "int",
      "value": 0
    }
  ],
  "outputs": [
  	{
      "name": "status",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| userKey 	 	| Your PushOver account User Key |         
| authToken  	| The PushOver auth token allocated for your app |
| message    	| The message to send |
| messageTitle  | The title of the message |
| device    	| The device where to send the message |
| sound	    	| The sound to play at reception of the message |
| url	    	| If you want to link to an URL in your message |
| urlTitle    	| Title to set on your URL if you use it |
| priority    	| Priority of the message |


In the 'status' output, you may get the following values:
- 'OK' : the message was correctly sent
- 'PUSH_ERR' : an error on sending the message via PushOver
- 'CONNECT_ERR' : if there was an error connecting to PushOver
- 'NO_MSG' : if the input 'message' field is empty

## Configuration Example

```json
            {  
            	"id": 2,
            	"name": "PushOver Notification",
            	"type":1,
            	"activityType":"pushover",
            	"attributes":[  
    				{
      					"name": "userKey",
      					"value": "YOUR_KEY",
      					"type": "string"
    				},
    				{
      					"name": "authToken",
      					"value": "YOUR_TOKEN",
      					"type": "string"
    				},
    				{
      					"name": "message",
      					"value": "YOUR MESSAGE",
      					"type": "string"
    				}
            	]
         	},
```


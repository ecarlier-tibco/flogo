package pushover

var jsonMetadata = `{
  "name": "pushover",
  "version": "0.0.1",
  "description": "Simple PushOver Notification Activity",
  "author": "Eric Carlier <ecarlier@tibco.com>",
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
}`

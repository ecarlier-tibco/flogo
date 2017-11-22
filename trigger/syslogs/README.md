# syslogs trigger
This trigger provides your flogo application the ability to start a flow upon reception of a syslog event

## Flogo CLI Installation

```bash
flogo install github.com/ecarlier-tibco/flogo/trigger/syslogs
```
## Flogo Web Installation

Install from following link : https://github.com/ecarlier-tibco/flogo/trigger/syslogs

## Trigger Settings
```json
{
  "settings":[
    {
      "name": "port",
      "type": "integer",
      "value": "514"
    },
    {
      "name": "UDP",
      "type": "boolean",
      "value": "true"
    },
    {
      "name": "TCP",
      "type": "boolean",
      "value": "false"
    }
  ]
}
```
| Setting     | Description    | Default Value |
|:------------|:---------------|:--------------|
| port | The port to listen on | 514 |
| UDP | If UDP protocol shall be enabled or not | true |
| TCP | If TCP protocol shall be enabled or not | false |

## Handler Settings

There is no specific configuration settings on the Handler.
Also so far, the trigger shall support one and only one handler, as there is no point at this stage to have more

Therefore in your application json file, you just need to specify the actionId of your action (typically the id of the flow to trigger)

## Trigger Outputs

There are many outputs. They may not be all set by the trigger. The list of set fields will actually depend on the format (RFC3164 vs RFC5424 for instance) and content of the syslog event received

_(In DEBUG mode, the trigger will log fields actually received for each event)_

```json
{  "outputs": [
    {
      "name": "content",
      "type": "string"
    },
    {
      "name": "tag",
      "type": "string"
    },
    {
      "name": "priority",
      "type": "integer"
    },
    {
      "name": "severity",
      "type": "integer"
    },
    {
      "name": "hostname",
      "type": "string"
    },
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "facility",
      "type": "integer"
    },
    {
      "name": "version",
      "type": "integer"
    },
    {
      "name": "timestamp",
      "type": "string"
    },
    {
      "name": "app_name",
      "type": "string"
    },
    {
      "name": "proc_id",
      "type": "string"
    },
    {
      "name": "msg_id",
      "type": "string"
    },
    {
      "name": "client",
      "type": "string"
    },
    {
      "name": "tls_peer",
      "type": "string"
    }
  ]
}
```
## Trigger Expected Action Results

There are no expected results back from the action to the syslogs trigger

## Example Application File using Syslogs Trigger

This application integrate the syslogs trigger and maps 3 fields issued by the trigger into the flow : (priority, severity, content) which are mapped into flow parameters (priority, severity, syslogmsg)

The Flow itself just has one activity logging flow parameter (syslogmsg)

```json
{
  "name": "syslogsApp",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "",
  "triggers": [
    {
      "name": "Syslogs Trigger",
      "ref": "github.com/ecarlier-tibco/flogo/trigger/syslogs",
      "description": "Syslog Trigger",
      "settings": {
        "port": "11514",
        "UDP": "true",
        "TCP": "false"
      },
      "id": "syslogs_trigger",
      "handlers": [
        {
          "actionMappings": {
            "input": [
              {
                "type": 1,
                "mapTo": "priority",
                "value": "priority"
              },
              {
                "type": 1,
                "mapTo": "severity",
                "value": "severity"
              },
              {
                "type": 1,
                "mapTo": "syslogmsg",
                "value": "content"
              }
            ],
            "output": []
          },
          "settings": {},
          "actionId": "syslogs_flow"
        }
      ]
    }
  ],
  "actions": [
    {
      "metadata": {
        "input": [
          {
            "name": "syslogmsg",
            "type": "string"
          },
          {
            "name": "priority",
            "type": "integer"
          },
          {
            "name": "severity",
            "type": "integer"
          }
        ],
        "output": []
      },
      "data": {
        "flow": {
          "name": "syslogs flow",
          "type": 1,
          "attributes": [],
          "rootTask": {
            "id": 1,
            "type": 1,
            "tasks": [
              {
                "id": "log_2",
                "name": "Log Message",
                "description": "Simple Log Activity",
                "type": 1,
                "activityType": "github-com-tibco-software-flogo-contrib-activity-log",
                "activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
                "attributes": [
                  {
                    "name": "message",
                    "value": "",
                    "required": false,
                    "type": "string"
                  },
                  {
                    "name": "flowInfo",
                    "value": "false",
                    "required": false,
                    "type": "boolean"
                  },
                  {
                    "name": "addToFlow",
                    "value": "false",
                    "required": false,
                    "type": "boolean"
                  }
                ],
                "inputMappings": [
                  {
                    "type": 1,
                    "value": "$flow.syslogmsg",
                    "mapTo": "message"
                  }
                ]
              }
            ],
            "links": [],
            "attributes": []
          }
        }
      },
      "id": "syslogs_flow",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow"
    }
  ]
}
```
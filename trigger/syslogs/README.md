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
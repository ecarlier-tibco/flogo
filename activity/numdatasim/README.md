# Random Value Generator
This activity is targeted to be used for testing or simulation of sensor data values like temperature or humidity for instance.

## Installation
### Flogo Web
Use following Activity URL:

https://github.com/ecarlier-tibco/flogo/activity/numdatasim
### Flogo CLI
```bash
flogo add activity github.com/ecarlier-tibco/flogo/activity/numdatasim
```

## Schema
Inputs and Outputs:

```json
  "inputs":[
    {
      "name" : "distributionType",
      "type" : "string",
      "required" : true,
      "allowed" : ["Normal", "Uniform", "Exponential", "Triangular"]
    },
    {
      "name": "repeatSequence",
      "type": "boolean",
      "value": false
    },
    {
      "name" : "min",
      "type" : "number"
    },
    {
      "name" : "max",
      "type" : "number"
    },
    {
      "name" : "mean",
      "type" : "number"
    },
    {
      "name" : "sigma",
      "type" : "number"
    },
    {
      "name" : "mode",
      "type" : "number"
    },
    {
      "name" : "lambda",
      "type" : "number"
    }
  ],
  "outputs": [
    {
      "name": "output",
      "type": "number"
    }
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| distributionType      | True     | The distribution type used for generating random values (Allowed values are Normal, Uniform, Exponential, Triangular) |         
| repeatSequence         | False     | Flag indicating if same seed shall be re-used on restart, meaning repeating the same generated values across different runs. Default Value  = false|
| min       | False    | Minimum value. Required for Uniform & Triangular distribution |
| max       | False    | Maximum value. Required for Uniform & Triangular distribution |
| mode       | False    | Required for Triangular distribution |
| mean       | False    | Mean value. Required for Normal distribution |
| sigma       | False    | Standard dev value. Required for Normal distribution |
| sigma       | False    | lambda value. Required for Exponential distribution |





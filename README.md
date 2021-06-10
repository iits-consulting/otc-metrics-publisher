# OTC metrics publisher

The simple daemon runs scripts from the specific folder to get and sends recieved metrics to OTC Cloud Eye.


## Script output
Next fields are expected:
|field|status|
|---|---|
|name|mandatory|
|value|mandatory|
|unit|optional|


### multiline output example
```ini
name=cpu_load_5min
value=50.0
unit="%"
```
or
```ini
name = "mem_free"
value = "4096"
unit = "Mb"
```


### JSON output example
```json
{
  "name": "network_in",
  "value": 42
}
```

### Single line (CSV) output example
The delimeter is `;` (semicolon) first field is trated as name the second as value.
The third field might contain type and is optional.
```cs
network_out;42;bytes
```

## Configuration with environment variables

|Variable|Default value|Description|
|---|---|---|
|METRICS_AUTHENDPOINT|https://iam.eu-de.otc.t-systems.com/v3"|optional|
|METRICS_PROJECT_ID|""|Mandatory|
|METRICS_ACCESS_KEY|""|Mandatory|
|METRICS_SECRET_KEY|""|Mandatory|
|METRICS_USER     |""|not implemented|
|METRICS_PASSWORD |""|not implemented|
|METRICS_AUTH_METHOD|"aksk"|Only implemented auth method is AK/SK|
|METRICS_NAMESPACE|APP.node|[otc docs](https://docs.otc.t-systems.com/api/ces/en-us_topic_0171212508.html#EN-US_TOPIC_0171212508__en-us_topic_0022067719_section24282572112133)|
|METRICS_SEND_INTERVA|60||
|METRICS_GRAB_INTERVAL|10||
|METRICS_SCRIPTS_DIR|/opt/metric-scripts|Every *executable* script in diretory will be launched|
|METRICS_INSTANCE_ID|"undefined"||
|METRICS_FILE_CLOUD_INIT_INSTANCE_ID STRING |/run/cloud-init/.instance-id|we can read instance id from that file (if not defined)|
|METRICS_TTL|86400|one day|
|METRICS_LOG_LEVEL|error|possible values are: trace, debug, info, warn, warning, error, fatal, panic, quite, nolog|
|METRICS_LOG_FORMAT|json||
|METRICS_LOG_TIME_FORMAT|unix||

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

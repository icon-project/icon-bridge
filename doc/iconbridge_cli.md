# Btpsimple

## iconbridge

### Description
Command Line Interface of Relay for Blockchain Transmission Protocol

### Usage
` iconbridge [flags] `

### Options
|Name,shorthand | Environment Variable | Required | Default | Description|
|---|---|---|---|---|
| --base_dir | ICONBRIDGE_BASE_DIR | false |  |  Base directory for data |
| --config, -c | ICONBRIDGE_CONFIG | false |  |  Parsing configuration file |
| --console_level | ICONBRIDGE_CONSOLE_LEVEL | false | trace |  Console log level (trace,debug,info,warn,error,fatal,panic) |
| --dst.address | ICONBRIDGE_DST_ADDRESS | true |  |  BTP Address of destination blockchain (PROTOCOL://NID.BLOCKCHAIN/BMC) |
| --dst.endpoint | ICONBRIDGE_DST_ENDPOINT | true |  |  Endpoint of destination blockchain |
| --dst.options | ICONBRIDGE_DST_OPTIONS | false | [] |  Options, comma-separated 'key=value' |
| --key_password | ICONBRIDGE_KEY_PASSWORD | false |  |  Password of KeyStore |
| --key_secret | ICONBRIDGE_KEY_SECRET | false |  |  Secret(password) file for KeyStore |
| --key_store | ICONBRIDGE_KEY_STORE | false |  |  KeyStore |
| --log_forwarder.address | ICONBRIDGE_LOG_FORWARDER_ADDRESS | false |  |  LogForwarder address |
| --log_forwarder.level | ICONBRIDGE_LOG_FORWARDER_LEVEL | false | info |  LogForwarder level |
| --log_forwarder.name | ICONBRIDGE_LOG_FORWARDER_NAME | false |  |  LogForwarder name |
| --log_forwarder.options | ICONBRIDGE_LOG_FORWARDER_OPTIONS | false | [] |  LogForwarder options, comma-separated 'key=value' |
| --log_forwarder.vendor | ICONBRIDGE_LOG_FORWARDER_VENDOR | false |  |  LogForwarder vendor (fluentd,logstash) |
| --log_level | ICONBRIDGE_LOG_LEVEL | false | debug |  Global log level (trace,debug,info,warn,error,fatal,panic) |
| --log_writer.compress | ICONBRIDGE_LOG_WRITER_COMPRESS | false | false |  Use gzip on rotated log file |
| --log_writer.filename | ICONBRIDGE_LOG_WRITER_FILENAME | false |  |  Log file name (rotated files resides in same directory) |
| --log_writer.localtime | ICONBRIDGE_LOG_WRITER_LOCALTIME | false | false |  Use localtime on rotated log file instead of UTC |
| --log_writer.maxage | ICONBRIDGE_LOG_WRITER_MAXAGE | false | 0 |  Maximum age of log file in day |
| --log_writer.maxbackups | ICONBRIDGE_LOG_WRITER_MAXBACKUPS | false | 0 |  Maximum number of backups |
| --log_writer.maxsize | ICONBRIDGE_LOG_WRITER_MAXSIZE | false | 100 |  Maximum log file size in MiB |
| --offset | ICONBRIDGE_OFFSET | false | 0 |  Offset of MTA |
| --src.address | ICONBRIDGE_SRC_ADDRESS | true |  |  BTP Address of source blockchain (PROTOCOL://NID.BLOCKCHAIN/BMC) |
| --src.endpoint | ICONBRIDGE_SRC_ENDPOINT | true |  |  Endpoint of source blockchain |
| --src.options | ICONBRIDGE_SRC_OPTIONS | false | [] |  Options, comma-separated 'key=value' |

### Child commands
|Command | Description|
|---|---|
| [iconbridge save](#iconbridge-save) |  Save configuration |
| [iconbridge start](#iconbridge-start) |  Start server |
| [iconbridge version](#iconbridge-version) |  Print iconbridge version |

## iconbridge save

### Description
Save configuration

### Usage
` iconbridge save [file] [flags] `

### Options
|Name,shorthand | Environment Variable | Required | Default | Description|
|---|---|---|---|---|
| --save_key_store |  | false |  |  KeyStore File path to save |

### Inherited Options
|Name,shorthand | Environment Variable | Required | Default | Description|
|---|---|---|---|---|
| --base_dir | ICONBRIDGE_BASE_DIR | false |  |  Base directory for data |
| --config, -c | ICONBRIDGE_CONFIG | false |  |  Parsing configuration file |
| --console_level | ICONBRIDGE_CONSOLE_LEVEL | false | trace |  Console log level (trace,debug,info,warn,error,fatal,panic) |
| --dst.address | ICONBRIDGE_DST_ADDRESS | true |  |  BTP Address of destination blockchain (PROTOCOL://NID.BLOCKCHAIN/BMC) |
| --dst.endpoint | ICONBRIDGE_DST_ENDPOINT | true |  |  Endpoint of destination blockchain |
| --dst.options | ICONBRIDGE_DST_OPTIONS | false | [] |  Options, comma-separated 'key=value' |
| --key_password | ICONBRIDGE_KEY_PASSWORD | false |  |  Password of KeyStore |
| --key_secret | ICONBRIDGE_KEY_SECRET | false |  |  Secret(password) file for KeyStore |
| --key_store | ICONBRIDGE_KEY_STORE | false |  |  KeyStore |
| --log_forwarder.address | ICONBRIDGE_LOG_FORWARDER_ADDRESS | false |  |  LogForwarder address |
| --log_forwarder.level | ICONBRIDGE_LOG_FORWARDER_LEVEL | false | info |  LogForwarder level |
| --log_forwarder.name | ICONBRIDGE_LOG_FORWARDER_NAME | false |  |  LogForwarder name |
| --log_forwarder.options | ICONBRIDGE_LOG_FORWARDER_OPTIONS | false | [] |  LogForwarder options, comma-separated 'key=value' |
| --log_forwarder.vendor | ICONBRIDGE_LOG_FORWARDER_VENDOR | false |  |  LogForwarder vendor (fluentd,logstash) |
| --log_level | ICONBRIDGE_LOG_LEVEL | false | debug |  Global log level (trace,debug,info,warn,error,fatal,panic) |
| --log_writer.compress | ICONBRIDGE_LOG_WRITER_COMPRESS | false | false |  Use gzip on rotated log file |
| --log_writer.filename | ICONBRIDGE_LOG_WRITER_FILENAME | false |  |  Log file name (rotated files resides in same directory) |
| --log_writer.localtime | ICONBRIDGE_LOG_WRITER_LOCALTIME | false | false |  Use localtime on rotated log file instead of UTC |
| --log_writer.maxage | ICONBRIDGE_LOG_WRITER_MAXAGE | false | 0 |  Maximum age of log file in day |
| --log_writer.maxbackups | ICONBRIDGE_LOG_WRITER_MAXBACKUPS | false | 0 |  Maximum number of backups |
| --log_writer.maxsize | ICONBRIDGE_LOG_WRITER_MAXSIZE | false | 100 |  Maximum log file size in MiB |
| --offset | ICONBRIDGE_OFFSET | false | 0 |  Offset of MTA |
| --src.address | ICONBRIDGE_SRC_ADDRESS | true |  |  BTP Address of source blockchain (PROTOCOL://NID.BLOCKCHAIN/BMC) |
| --src.endpoint | ICONBRIDGE_SRC_ENDPOINT | true |  |  Endpoint of source blockchain |
| --src.options | ICONBRIDGE_SRC_OPTIONS | false | [] |  Options, comma-separated 'key=value' |

### Parent command
|Command | Description|
|---|---|
| [iconbridge](#iconbridge) |  BTP Relay CLI |

### Related commands
|Command | Description|
|---|---|
| [iconbridge save](#iconbridge-save) |  Save configuration |
| [iconbridge start](#iconbridge-start) |  Start server |
| [iconbridge version](#iconbridge-version) |  Print iconbridge version |

## iconbridge start

### Description
Start server

### Usage
` iconbridge start [flags] `

### Options
|Name,shorthand | Environment Variable | Required | Default | Description|
|---|---|---|---|---|
| --cpuprofile |  | false |  |  CPU Profiling data file |
| --memprofile |  | false |  |  Memory Profiling data file |

### Inherited Options
|Name,shorthand | Environment Variable | Required | Default | Description|
|---|---|---|---|---|
| --base_dir | ICONBRIDGE_BASE_DIR | false |  |  Base directory for data |
| --config, -c | ICONBRIDGE_CONFIG | false |  |  Parsing configuration file |
| --console_level | ICONBRIDGE_CONSOLE_LEVEL | false | trace |  Console log level (trace,debug,info,warn,error,fatal,panic) |
| --dst.address | ICONBRIDGE_DST_ADDRESS | true |  |  BTP Address of destination blockchain (PROTOCOL://NID.BLOCKCHAIN/BMC) |
| --dst.endpoint | ICONBRIDGE_DST_ENDPOINT | true |  |  Endpoint of destination blockchain |
| --dst.options | ICONBRIDGE_DST_OPTIONS | false | [] |  Options, comma-separated 'key=value' |
| --key_password | ICONBRIDGE_KEY_PASSWORD | false |  |  Password of KeyStore |
| --key_secret | ICONBRIDGE_KEY_SECRET | false |  |  Secret(password) file for KeyStore |
| --key_store | ICONBRIDGE_KEY_STORE | false |  |  KeyStore |
| --log_forwarder.address | ICONBRIDGE_LOG_FORWARDER_ADDRESS | false |  |  LogForwarder address |
| --log_forwarder.level | ICONBRIDGE_LOG_FORWARDER_LEVEL | false | info |  LogForwarder level |
| --log_forwarder.name | ICONBRIDGE_LOG_FORWARDER_NAME | false |  |  LogForwarder name |
| --log_forwarder.options | ICONBRIDGE_LOG_FORWARDER_OPTIONS | false | [] |  LogForwarder options, comma-separated 'key=value' |
| --log_forwarder.vendor | ICONBRIDGE_LOG_FORWARDER_VENDOR | false |  |  LogForwarder vendor (fluentd,logstash) |
| --log_level | ICONBRIDGE_LOG_LEVEL | false | debug |  Global log level (trace,debug,info,warn,error,fatal,panic) |
| --log_writer.compress | ICONBRIDGE_LOG_WRITER_COMPRESS | false | false |  Use gzip on rotated log file |
| --log_writer.filename | ICONBRIDGE_LOG_WRITER_FILENAME | false |  |  Log file name (rotated files resides in same directory) |
| --log_writer.localtime | ICONBRIDGE_LOG_WRITER_LOCALTIME | false | false |  Use localtime on rotated log file instead of UTC |
| --log_writer.maxage | ICONBRIDGE_LOG_WRITER_MAXAGE | false | 0 |  Maximum age of log file in day |
| --log_writer.maxbackups | ICONBRIDGE_LOG_WRITER_MAXBACKUPS | false | 0 |  Maximum number of backups |
| --log_writer.maxsize | ICONBRIDGE_LOG_WRITER_MAXSIZE | false | 100 |  Maximum log file size in MiB |
| --offset | ICONBRIDGE_OFFSET | false | 0 |  Offset of MTA |
| --src.address | ICONBRIDGE_SRC_ADDRESS | true |  |  BTP Address of source blockchain (PROTOCOL://NID.BLOCKCHAIN/BMC) |
| --src.endpoint | ICONBRIDGE_SRC_ENDPOINT | true |  |  Endpoint of source blockchain |
| --src.options | ICONBRIDGE_SRC_OPTIONS | false | [] |  Options, comma-separated 'key=value' |

### Parent command
|Command | Description|
|---|---|
| [iconbridge](#iconbridge) |  BTP Relay CLI |

### Related commands
|Command | Description|
|---|---|
| [iconbridge save](#iconbridge-save) |  Save configuration |
| [iconbridge start](#iconbridge-start) |  Start server |
| [iconbridge version](#iconbridge-version) |  Print iconbridge version |

## iconbridge version

### Description
Print iconbridge version

### Usage
` iconbridge version `

### Inherited Options
|Name,shorthand | Environment Variable | Required | Default | Description|
|---|---|---|---|---|
| --base_dir | ICONBRIDGE_BASE_DIR | false |  |  Base directory for data |
| --config, -c | ICONBRIDGE_CONFIG | false |  |  Parsing configuration file |
| --console_level | ICONBRIDGE_CONSOLE_LEVEL | false | trace |  Console log level (trace,debug,info,warn,error,fatal,panic) |
| --dst.address | ICONBRIDGE_DST_ADDRESS | true |  |  BTP Address of destination blockchain (PROTOCOL://NID.BLOCKCHAIN/BMC) |
| --dst.endpoint | ICONBRIDGE_DST_ENDPOINT | true |  |  Endpoint of destination blockchain |
| --dst.options | ICONBRIDGE_DST_OPTIONS | false | [] |  Options, comma-separated 'key=value' |
| --key_password | ICONBRIDGE_KEY_PASSWORD | false |  |  Password of KeyStore |
| --key_secret | ICONBRIDGE_KEY_SECRET | false |  |  Secret(password) file for KeyStore |
| --key_store | ICONBRIDGE_KEY_STORE | false |  |  KeyStore |
| --log_forwarder.address | ICONBRIDGE_LOG_FORWARDER_ADDRESS | false |  |  LogForwarder address |
| --log_forwarder.level | ICONBRIDGE_LOG_FORWARDER_LEVEL | false | info |  LogForwarder level |
| --log_forwarder.name | ICONBRIDGE_LOG_FORWARDER_NAME | false |  |  LogForwarder name |
| --log_forwarder.options | ICONBRIDGE_LOG_FORWARDER_OPTIONS | false | [] |  LogForwarder options, comma-separated 'key=value' |
| --log_forwarder.vendor | ICONBRIDGE_LOG_FORWARDER_VENDOR | false |  |  LogForwarder vendor (fluentd,logstash) |
| --log_level | ICONBRIDGE_LOG_LEVEL | false | debug |  Global log level (trace,debug,info,warn,error,fatal,panic) |
| --log_writer.compress | ICONBRIDGE_LOG_WRITER_COMPRESS | false | false |  Use gzip on rotated log file |
| --log_writer.filename | ICONBRIDGE_LOG_WRITER_FILENAME | false |  |  Log file name (rotated files resides in same directory) |
| --log_writer.localtime | ICONBRIDGE_LOG_WRITER_LOCALTIME | false | false |  Use localtime on rotated log file instead of UTC |
| --log_writer.maxage | ICONBRIDGE_LOG_WRITER_MAXAGE | false | 0 |  Maximum age of log file in day |
| --log_writer.maxbackups | ICONBRIDGE_LOG_WRITER_MAXBACKUPS | false | 0 |  Maximum number of backups |
| --log_writer.maxsize | ICONBRIDGE_LOG_WRITER_MAXSIZE | false | 100 |  Maximum log file size in MiB |
| --offset | ICONBRIDGE_OFFSET | false | 0 |  Offset of MTA |
| --src.address | ICONBRIDGE_SRC_ADDRESS | true |  |  BTP Address of source blockchain (PROTOCOL://NID.BLOCKCHAIN/BMC) |
| --src.endpoint | ICONBRIDGE_SRC_ENDPOINT | true |  |  Endpoint of source blockchain |
| --src.options | ICONBRIDGE_SRC_OPTIONS | false | [] |  Options, comma-separated 'key=value' |

### Parent command
|Command | Description|
|---|---|
| [iconbridge](#iconbridge) |  BTP Relay CLI |

### Related commands
|Command | Description|
|---|---|
| [iconbridge save](#iconbridge-save) |  Save configuration |
| [iconbridge start](#iconbridge-start) |  Start server |
| [iconbridge version](#iconbridge-version) |  Print iconbridge version |


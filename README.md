# nomon

#### Config example:
```
cpu_alert_threshold: 50 #percentage
ram_alert_threshold: 50 #percentage
disk_alert_threshold: 50 #percentage
check_every: 30  # time between alert checks in seconds
port: 8000
old_data_cleanup: 4  # when database cleans up: 0-23
log_level: 5  # https://github.com/sirupsen/logrus
urls: [ ]  # url examples: https://containrrr.dev/shoutrrr/0.7/
```
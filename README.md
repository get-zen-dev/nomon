# nomon

#### Config:
./data/config.yml
```
args:
  cpu_limit: 50
  cpu_cycles: 2
  ram_limit: 50
  ram_cycles: 2
  disk_limit: 50
  disk_cycles: 2
  duration: 300  # time between alert checks (at least 30 seconds)
  port: 8000
  db_clear_time: 14  # database clearing start hour 0-23 (UTC)
  monitor_log_level: 5
report: 
  service: "matrix"
  matrix_access_token: "matrix_token"
  matrix_room_id: "matrix_room_id"
  matrix_host_server: "matrix_server"

```

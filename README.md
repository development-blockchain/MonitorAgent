# MonitorAgent
本模块提供了系统信息的自动监控，以及提供了可接入kafka的日志模块.

# 如何运行
```
# go build 
# ./MonitorAgent
```
运行后通过访问 `http://127.0.0.1:9000/metrics` 可以获得监控信息.
```curl
curl http://127.0.0.1:9000/metrics
```
目前已经包含的信息有：
```
cpu count
cpu load
mem total
mem available
mem free
mem used
mem used_percent
disk total
disk free
disk used
disk used_percent
```
信息如下,其中monitor开头的是本模块增加的监控内容，其他为prometheus自动添加的监控数据：
```
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
...
# HELP monitor_systeminfo_disk_usage_free 
# TYPE monitor_systeminfo_disk_usage_free gauge
monitor_systeminfo_disk_usage_free{hostname="myself",path="/"} 5.3633957888e+10
monitor_systeminfo_disk_usage_free{hostname="myself",path="/boot/efi"} 2.29752832e+08
monitor_systeminfo_disk_usage_free{hostname="myself",path="/home/luxq/opt"} 1.21935384576e+11
# HELP monitor_systeminfo_disk_usage_total 
# TYPE monitor_systeminfo_disk_usage_total gauge
monitor_systeminfo_disk_usage_total{hostname="myself",path="/"} 2.0528734208e+11
monitor_systeminfo_disk_usage_total{hostname="myself",path="/boot/efi"} 2.68435456e+08
monitor_systeminfo_disk_usage_total{hostname="myself",path="/home/luxq/opt"} 2.06032830464e+11
# HELP monitor_systeminfo_disk_usage_used 
# TYPE monitor_systeminfo_disk_usage_used gauge
monitor_systeminfo_disk_usage_used{hostname="myself",path="/"} 1.41150797824e+11
monitor_systeminfo_disk_usage_used{hostname="myself",path="/boot/efi"} 3.8682624e+07
monitor_systeminfo_disk_usage_used{hostname="myself",path="/home/luxq/opt"} 7.355695104e+10
# HELP monitor_systeminfo_disk_usage_used_percent 
# TYPE monitor_systeminfo_disk_usage_used_percent gauge
monitor_systeminfo_disk_usage_used_percent{hostname="myself",path="/"} 72.46501262793852
monitor_systeminfo_disk_usage_used_percent{hostname="myself",path="/boot/efi"} 14.410400390625
monitor_systeminfo_disk_usage_used_percent{hostname="myself",path="/home/luxq/opt"} 37.626514005380656
# HELP monitor_systeminfo_host_cpu_count 
# TYPE monitor_systeminfo_host_cpu_count gauge
monitor_systeminfo_host_cpu_count{hostname="myself"} 8
# HELP monitor_systeminfo_host_cpu_load 
# TYPE monitor_systeminfo_host_cpu_load gauge
monitor_systeminfo_host_cpu_load{hostname="myself"} 2.1813138317989553
# HELP monitor_systeminfo_host_processor_count 
# TYPE monitor_systeminfo_host_processor_count gauge
monitor_systeminfo_host_processor_count{hostname="myself"} 446
# HELP monitor_systeminfo_host_uptime 
# TYPE monitor_systeminfo_host_uptime gauge
monitor_systeminfo_host_uptime{hostname="myself"} 4962
# HELP monitor_systeminfo_mem_available 
# TYPE monitor_systeminfo_mem_available gauge
monitor_systeminfo_mem_available{hostname="myself"} 2.2306103296e+10
# HELP monitor_systeminfo_mem_total 
# TYPE monitor_systeminfo_mem_total gauge
monitor_systeminfo_mem_total{hostname="myself"} 3.3340608512e+10
# HELP monitor_systeminfo_mem_used 
# TYPE monitor_systeminfo_mem_used gauge
monitor_systeminfo_mem_used{hostname="myself"} 8.015663104e+09
# HELP monitor_systeminfo_mem_used_percent 
# TYPE monitor_systeminfo_mem_used_percent gauge
monitor_systeminfo_mem_used_percent{hostname="myself"} 24.041742072928844
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.01
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
...
```

# 日志模块如何使用
参考示例代码，日志模块一定会将日志存储到本地文件，如果配置了kafka，那么同时还会将日志发往kafka对应的topic中.
```go
func TestLocalSystemLogger(T *testing.T) {
	config := LogConfig{
		Kafka: nil,
		Path:  "logs",
		Level: "debug",
	}
	InitLog(config)
	l := Entry()
	l.Info("this is info")
	l.Debug("this is debug")
	l.Error("this is error")
}

func TestKafkaLogger(T *testing.T) {
	config := LogConfig{
		Kafka: &KafkaConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "nodelog",
		},
		Path:  "logs",
		Level: "debug",
	}
	InitLog(config)
	l := Entry()
	l.Info("this is info")
	l.Debug("this is debug")
	l.Error("this is error")
	time.Sleep(time.Second)
}

```
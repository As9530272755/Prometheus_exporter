1 Configure the config.yaml file

```yaml
# serivce config
service:
  name: "systemd-timesyncd"
  config_path: "/etc/systemd/timesyncd.conf"

# web 
web:
  #The default address is 19090
  addr:

#Log Configuration

#Default log path: log/time under the current path_ exporter.log

#There are 7 levels of log, including panic total error warn info debug trace
log: 
  filename: log/time_exporter.log
  max_age: 10
```



2 Start program

```bash
[17:32:05 root@go ntp]#./timeService_exporter 
```



3 Browser corresponding indicators

http://10.0.0.135:19090/metrics

![images](https://github.com/As9530272755/Prometheus_exporter/blob/v1.0.beta/Time_Service_Exporter/image-20230412173307594.png)



4 Prometheus view docked

![images](https://github.com/As9530272755/Prometheus_exporter/blob/v1.0.beta/Time_Service_Exporter/image-20230412173444276.png)



5 alertmanager 

![image-20230412173518738](ntp-exporter 开发.assets/image-20230412173518738.png)



6 Stop the service and modify the configuration file

```bash
[17:36:39 root@go ntp]#systemctl stop systemd-timesyncd.service 
[17:37:19 root@go ntp]#echo "sdad" >> /etc/systemd/timesyncd.conf 
```



7 alerts have turned red

![image-20230412173833216](ntp-exporter 开发.assets/image-20230412173833216.png)



8 Enterprise WeChat Receiving Alarm

![image-20230412173842108](ntp-exporter 开发.assets/image-20230412173842108.png)

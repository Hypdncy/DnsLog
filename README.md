# DnsLog
增加basic auth


安装
---

假设你的域名在阿里云，为example.com，VPS的IP为1.1.1.1

1. 在阿里云配置A记录，ns1.example.com 为1.1.1.1
2. 在阿里云配置NS记录，*.log.example.com 为ns1.example.com
3. 配置json文件config.DNS.domain的值为log.example.com（与上面对应)
4. 访问web：http://1.1.1.1:8053/
5. 密码在config.USER

# 1.获取发行版
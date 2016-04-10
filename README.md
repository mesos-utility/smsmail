smsmail:短信/邮件发送接口
================================
[![Build Status](https://travis-ci.org/mesos-utility/smsmail.png?branch=master)](https://travis-ci.org/mesos-utility/smsmail)

系统需求
--------------------------------
操作系统：Linux

主要逻辑
--------------------------------
配置短信及邮箱发送接口,以http方式调用该接口发送短信及邮件

使用方法
--------------------------------
1. 根据实际部署情况，配置mail及sms发送网关;
 * mail: "addr": "10.10.10.10:25", "username": "", "password": ""
 * sms:  "url": "10.10.10.10:20230"

2. 测试： ./control build && ./control start    OR
 * glide create                            # Start a new workspace
 * glide get github.com/soarpenguin/smsmail # Get a package and add to glide.yaml
 * glide install                           # Install packages and dependencies
 * go build                                # Go tools work normally
 * glide up                                # Update to newest versions of the package

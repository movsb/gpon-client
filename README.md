# 天翼网关客户端

这个程序用来在命令行下控制天翼网关路由器的一些常用配置。

## 帮助

```bash
$ ./gpon-client
A GPON (Tiānyì Gateway) client used to modify router configurations.

Usage: ./gpon-client command [sub-command] [<arguments>...]

All Available Command List:

        portmaps list
        portmaps create  <name> <protocol> <outer-port> <inner-ip> <inner-port>
        portmaps delete  <name>
        portmaps enable  <name>
        portmaps disable <name>

```

## 运行环境初始化配置

请先在命令行下导出几个环境变量。

```bash
# 路由器IP地址，默认 192.168.1.1
$ export IP=192.168.1.1

# 路由器用户名，默认 useradmin
$ export USERNAME=useradmin

# 路由器密码，无默认
$ export PASSWORD=
```

## 示例使用

### 端口转发配置

#### 枚举端口映射

```bash
$ ./gpon-client portmaps list
ID   Name            Protocol    OuterPort   InnerIP             InnerPort   Enable
-----------------------------------------------------------------------------------
1    test            BOTH        4321        192.168.1.6         1234        1     
2    nginx           TCP         8888        192.168.1.6         8888        1     
3    ssh             TCP         10022       192.168.1.6         22          1     
4    bt              BOTH        8999        192.168.1.6         8999        1     
5    https           TCP         443         192.168.1.6         443         1  
```

#### 添加端口映射

```bash
$ ./gpon-client portmaps create nginx BOTH 443 192.168.1.6 443
```

#### 删除端口映射

```bash
$ ./gpon-client portmaps delete nginx
```

#### 启用端口映射

```bash
$ ./gpon-client portmaps enable nginx
```

#### 禁用端口映射

```bash
$ ./gpon-client portmaps disable nginx
```

## 版权

Copyright (C) 2019 movsb

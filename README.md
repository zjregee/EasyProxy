## EasyProxy

这是一个简易的代理服务，用Go语言实现，该代理服务器实现了反向代理以及对静态资源的代理

在TCP的基础上实现对HTTP报文的解析和构建

代理和日志配置文件仿照nginx

未使用http相关标准库以及第三方库

TODO：

* 对配置文件的读取
* 对HTTP更多的支持
* 静态代理的实现
* 负载均衡的实现
* bug的修复
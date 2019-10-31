# VLOG
使用go语言开发的vlog框架

[![github](https://badgen.net/badge/golang/1.13/green)](https://github.com/golang/go)
[![github](https://badgen.net/badge/build/passing/green)](#)
[![github](https://badgen.net/badge/license/GUN/green)](https://github.com/Rejudge-F/ShortLink/blob/master/LICENSE)

# 接口定义
## vlog播放接口
- 接口描述：获取一条vlog视频
- 请求方法：GET
- 请求地址：http://domain/video/{video_name.mp4}
- 状态码：  200 successful

## vlog上传接口
- 接口描述：上传一条vlog，视频大小不能超过10m
- 请求方法：POST
- 请求地址：http://domain/api/upload
- 状态码：  200 successful 500 Internal

## vlog列表接口
- 接口描述：列出所有的vlog
- 请求方法：GET
- 请求地址：http://domain/api/list
- 状态码：  200 successful

# 视图HTML
使用SUI框架开发，通过AJAX与后端数据进行数据交换，文件通过multipart/form-data来传输

# 如何部署
将view视图与controllers中的所有文件中localhost替换为你所需要的部署机器的ip即可

# 测试网址

47.107.54.187：8088

由于部署机器为阿里云低配机，网速比较慢，需要等待一会才能将视频传输完毕

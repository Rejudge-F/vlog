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

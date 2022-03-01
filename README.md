# KSUV

**一个自动化上传快手视频的工具**

## 特别文件/目录

| 文件/目录           | 说明                                                |
| ------------------- | --------------------------------------------------- |
| config.example.json | 配置文件模板，需改成"config.json"，程序才可读取配置 |
| title.txt           | 标题文件,每个标题一行                               |

## config.json

```json
{
    "version": "1.0",
    "name": "快手自动上传视频",
    "authors": "NaturalGao",
    "userConfig": {
        "cookie": "", // 用户 Cookie（必填）
        "webApiPh": "" // 客户端标识符（必填）
    },
    "titleFileUrl": "./title.txt", // 标题文件路径（必填）
    "videoFileUrl": "./video/", // 视频文件路径（必填）
    "secondDomain": "其它描述" // 视频其它描述
}
```



## cookie 和 webApiPh 获取方式

登录 ”快手创作者服务平台“，打开开发者工具，通常快捷键F12



<img src="https://i.bmp.ovh/imgs/2022/02/1e11a828931ad05f.png" style="zoom: 50%;" />



上图画红色框的就是需要的 Cookie，webApiPh 就是 Cookie 下的



## 程序运行效果

配置正常的情况下：

1. 选择是否开始程序

<img src="https://i.bmp.ovh/imgs/2022/02/17916e6f64f8c5fa.png" style="zoom: 67%;" />



2. 选择关联的商品

   <img src="https://i.bmp.ovh/imgs/2022/02/261a44472480b8d3.png" style="zoom:67%;" />

3. 上传

   <img src="https://i.bmp.ovh/imgs/2022/02/3c0e5575e89261fd.png" style="zoom:67%;" />


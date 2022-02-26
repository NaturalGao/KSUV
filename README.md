# KSUV

**一个自动化上传快手视频的工具**

## 特别文件/目录

| 文件/目录 | 说明                               |
| --------- | ---------------------------------- |
| config    | 配置文件,使用时需加上 '.json' 后缀 |
| title.txt | 标题文件,每个标题一行              |

## config.json

```json
{
    "version": "1.0",
    "name": "快手自动上传视频",
    "authors": "NaturalGao",
    "userConfig": {
        "cookie": "", // 用户 Cookie（必填）
        "uploadToken": "", // 上传 Token（必填）
        "webApiPh": "" // 客户端标识符（必填）
    },
    "titleFileUrl": "./title.txt", // 标题文件路径（必填）
    "videoFileUrl": "./video/", // 视频文件路径（必填）
    "secondDomain": "其它描述" // 视频其它描述
}
```

其它相关信息待更新...

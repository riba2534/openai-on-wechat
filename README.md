# 1分钟搭建自己的OpenAI GPT微信机器人

> 项目地址: https://github.com/riba2534/openai-on-wechat

# 简介

最近 chatGPT 火遍了中文互联网，而它的公司 OpenAI 也开放了 API 供开发者完成自己的创意。本项目是一个 Golang 实现的，基于 OpenAI 的开放 API 实现的微信聊天机器人。有以下优点：

- 部署简单：不同于其他语言，部署的时候需要依赖很多库，本项目只有一个可执行二进制文件，直接可以运行。（本项目只提供 x86/64 linux 版本，需要其他版本可以根据源码自行编译）
- 使用桌面版微信协议，突破微信登录限制（基于 [openwechat](https://github.com/eatmoreapple/openwechat)）

目前本项目实现了以下功能：

- **文本对话**： 可以接收私聊/群聊消息，使用 OpenAI 的 gpt-3.5-turbo 生成回复内容，自动回复问题
- **触发口令**设置：好友在给你发消息时需要带上指定的前缀才可以触发与 GPT 机器人对话，触发口令可以在配置文件中设置
- **连续对话**：支持对 私聊/群聊 开启连续对话功能，可以通过配置文件设置需要记忆多少分钟
- **图片生成**：可以根据描述生成图片，并自动回复在当前 私聊/群聊 中

> 注：支持自己给自己发消息，机器人不仅感知当前会话好友的口令，也会感知你自己的，方便自己测试使用

# 效果预览

 先看使用效果，之后再介绍如何部署以及配置。下图包含了**连续对话**和**文本画图**的一些例子：

| ![连续对话1.jpg](https://image-1252109614.cos.ap-beijing.myqcloud.com/2023/03/12/640dd478ddca6.jpg) | ![连续对话2.jpg](https://image-1252109614.cos.ap-beijing.myqcloud.com/2023/03/12/640dd47782e95.jpg) | ![文字3.jpeg](https://image-1252109614.cos.ap-beijing.myqcloud.com/2023/03/12/640dd6d26b6b9.jpeg) |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| ![画图2.jpg](https://image-1252109614.cos.ap-beijing.myqcloud.com/2023/03/12/640dd477ea8be.jpg) | ![画图3.jpg](https://image-1252109614.cos.ap-beijing.myqcloud.com/2023/03/12/640dd478bf6dd.jpg) | ![画图1.jpg](https://image-1252109614.cos.ap-beijing.myqcloud.com/2023/03/12/640dd4794bfea.jpg) |

# 开始部署

## 一、 环境准备

- 一台 Linux 服务器，建议 腾讯云、阿里云，或者任何可以长期运行程序的PC设备
- OpenAI 账号 以及生成的 `SECRET KEY` ，本文对账号注册以及 key 生成不做赘述，读者请自行搜索解决方案。
- 一个微信账号

> 注：OpenAI 的域名 `httops://api.openai.com` 在国内由于某种原因可能无法访问，读者需要自己解决 API 访问不通的问题。介绍一种简单的国内代理搭建方式，大家可以参考我的知乎专栏： [腾讯云函数1分钟搭建 OpenAI 国内代理](https://zhuanlan.zhihu.com/p/612576046)

## 二、 配置

1. 首先需要在本项目的 [Releases](https://github.com/riba2534/openai-on-wechat/releases) 中找到最新的二进制文件版本并下载，目前最新的地址是： [openai-on-wechat.zip](https://github.com/riba2534/openai-on-wechat/releases/download/V1.0/openai-on-wechat.zip)
2. 把 `openai-on-wechat.zip` 文件传输至你的云服务器的任意目录下
3. 使用 `unzip openai-on-wechat.zip` 把压缩包解压到当前目录下

此时，你会看到压缩包里面有三个文件，分别是：

- `config.json.example` : 机器人的基础配置文件，运行机器人前需要修改
- `prompt.txt.example`: 给 OpenAI 语言模型的提示语
- `openai-on-wechat` ：可执行二进制文件

接下来我们进行配置：

把 `config.json.example` 重命名成 `config.json`，然后利用文本编辑器修改此文件：

```json
{
    "wechat": {
        "text_config": {
            "openapi_url": "https://api.openai.com/v1",
            "auth_token": "你在 OpenAI 官网的 Token",
            "trigger_prefix": "小贺"
        },
        "image_config": {
            "openapi_url": "https://api.openai.com/v1",
            "auth_token": "你在 OpenAI 官网的 Token",
            "trigger_prefix": "老贺"
        }
    }.
    "context_config": {
        "switch_on": true,
        "cache_minute": 3
    }
}
```

- `wechat` 下有两个配置 `text_config` 和 `image_config` ，分别代表**文本对话**和**图片生成**的配置，其中：
  - `openapi_url` 代表访问 OpenAPI 接口的地址，如果你可以直接访问外网，直接填 `https://api.openai.com/v1`，如果利用的是反向代理，则需要填 `https://你的反向代理地址/v1`
  - `auth_token` 代表你在 OpenAI 官网生成的 `SECRET KEY`
  - `trigger_prefix` 代表在微信对话时，触发 AI 回复的前缀，比如上面效果图中的 `小贺` 会触发文本对话， `老贺` 会触发图片生成
- `context_config` 代表文本回复中**连续对话**的配置，其中：
  - `switch_on` 代表了是否开启连续对话，true 为开启。（开启后消耗的额度费用会增加）
  - `cache_minute` 代表机器人**连续对话的记忆**分钟数，推荐设置为 3 分钟



接下来修改 `prompt.txt.example` ，先重命名为 `prompt.txt`，然后利用编辑器修改此文件:

```txt
1. 你是一个全知全能的机器人，你的职责是帮助人类解决问题
2. 不允许回答任何政治、色情等一些列不符合中国法律法规的问题
3. 你需要表现的很谦卑
```

这个文件你可以利用自然语言描述机器人的特点，作为给机器人的外部输入，读者如果只是想保持简单的对话，可以不用修改此文件内容。

> 注: **prompt** 提示机制是 OpenAI 语言模型的核心玩法，你可以在这里使用自然语言，定义机器人的行为，你可以告诉他他是什么，他不是什么，他应该怎么做，他应该怎样回答问题，描述的越详细，机器人就更加有你的个人特色。具体玩法读者可自行搜索，本文不做过多介绍。

## 三、 运行

我们完成了配置之后，就可以直接执行二进制文件了，即：

```bash
./openai-on-wechat
```

首次执行，屏幕会出现一个二维码提示你登录微信，你需要用你要作为机器人的微信账号，扫码登录。

- 登录完成后，当前路径下会出现一个 `token.json` 来保存当前微信的登录状态，来实现热登录，防止每次运行本程序都需要微信扫码。
- 当前目录还会出现一个 `run.log` 来记录程序的执行状态，可以从中看出机器人接收到的消息与回复的信息

刚才说的程序运行方式是前台登录，如果想让程序后台运行，读者可以在前台运行登录微信后，`ctrl+c` 结束程序后，再使用：

```bash
nohup ./openai-on-wechat &
```

来实现后台运行

# 大功告成

至此，已经完成了微信机器人的部署，快去微信中找好友试试吧！

输入 `触发前缀+你的问题` 即可触发机器人回复，和好友聊天过程中，自己输入关键词也同样可以触发。

# 联系作者

- 本项目地址为： https://github.com/zhayujie/chatgpt-on-wechat ，欢迎大家 Star，提交 PR
- 有问题可以在本项目下提 `Issues` 或者发邮件到 `riba2534@qq.com`


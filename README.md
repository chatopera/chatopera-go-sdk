# chatopera-go-sdk

企业聊天机器人-Go 开发工具包

<p align="center">
    <b>Chatopera开发者平台：809987971， <a href="https://jq.qq.com/?_wv=1027&k=5S51T2a" target="_blank">点击链接加入群聊</a></b><br>
    <img src="https://user-images.githubusercontent.com/3538629/48105854-0bfcca00-e274-11e8-8eb4-ffb46a2c9179.png" width="200">
  </p>
  
  
# [chatopera-go-sdk](https://github.com/chatopera/chatopera-go-sdk)
企业聊天机器人-Go开发工具包

本教程介绍如何使用 Chatopera 机器人开发者平台的[Go SDK](https://github.com/chatopera/chatopera-go-sdk)与机器人进行集成，阅读本教程需要 10 分钟时间。

## 安装

```
go get github.com/chatopera/chatopera-go-sdk
```

## 使用文档

快速开始，类接口定义和实例化文档等，参考 [文档中心](https://docs.chatopera.com/products/chatbot-platform/integration.html)：

[https://docs.chatopera.com/products/chatbot-platform/integration.html](https://docs.chatopera.com/products/chatbot-platform/integration.html)

## 示例程序

假设您已经:

1. 准备好**ClientId**和**Secret**了；

2. 安装了[chatopera-go-sdk](https://github.com/chatopera/chatopera-go-sdk)，

那么，可以用以下代码测试。

```go
import (
	"github.com/chatopera/chatopera-go-sdk"
)
var chatbot = chatopera.Chatbot("YOUR CLIENT ID", "YOUR SECRET")
```

更多参考代码：

[https://github.com/chatopera/chatopera-go-sdk/blob/master/chatopera_test.go](https://github.com/chatopera/chatopera-go-sdk/blob/master/chatopera_test.go)

## 单元测试

单元测试: [chatopera_test.go](https://github.com/chatopera/chatopera-go-sdk/blob/master/chatopera_test.go)

```
cd chatopera-go-sdk
export GOPATH=YOUR_PROJECT_GOPATH:$GOPATH
go test -run ''
```

## 开源许可协议

Copyright (2018-2020) [北京华夏春松科技有限公司](https://www.chatopera.com/)

[Apache License Version 2.0](./LICENSE)

Copyright 2018-2020, [北京华夏春松科技有限公司](https://www.chatopera.com/). All rights reserved. This software and related documentation are provided under a license agreement containing restrictions on use and disclosure and are protected by intellectual property laws. Except as expressly permitted in your license agreement or allowed by law, you may not use, copy, reproduce, translate, broadcast, modify, license, transmit, distribute, exhibit, perform, publish, or display any part, in any form, or by any means. Reverse engineering, disassembly, or decompilation of this software, unless required by law for interoperability, is prohibited.

[![chatoper banner][co-banner-image]][co-url]

[co-banner-image]: https://user-images.githubusercontent.com/3538629/42383104-da925942-8168-11e8-8195-868d5fcec170.png
[co-url]: https://www.chatopera.com

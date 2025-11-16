# Swagger API 文档使用指南

## 概述

项目已集成 Swagger/OpenAPI 文档，提供了完整的 API 接口文档和在线测试功能。

## 访问 Swagger UI

启动服务后，通过浏览器访问：

```
http://localhost:8080/swagger/index.html
```

## 功能特性

- ✅ 完整的 API 接口列表
- ✅ 请求/响应参数说明
- ✅ 在线测试功能
- ✅ 交互式 API 文档
- ✅ 自动生成的代码示例

## API 分组

Swagger 文档按功能模块分组：

- **STT**: 语音识别相关接口（REST API - 批量处理）
- **TTS**: 语音合成相关接口（REST API - 批量处理）
- **系统**: 健康检查、统计信息、监控等系统接口

## 重要说明：流式处理 vs 批量处理

### REST API（批量处理）

Swagger 文档中展示的是 **REST API**，主要用于批量处理：

- **STT**: 
  - `POST /api/v1/stt/recognize` - 文件上传识别（批量）
  - `POST /api/v1/stt/batch` - 批量识别多个文件
- **TTS**: 
  - `POST /api/v1/tts/synthesize` - 文本合成（批量）
  - `POST /api/v1/tts/batch` - 批量合成多个文本

### WebSocket API（流式处理）

**流式处理**通过 **WebSocket** 接口实现，不在 Swagger 文档中展示：

- **STT 流式识别**: `ws://localhost:8080/ws/stt`
  - 实时接收音频流
  - 实时返回识别结果
  - 支持部分结果（partial）和最终结果（final）

- **TTS 流式合成**: `ws://localhost:8080/ws/tts`
  - 实时接收文本
  - 实时返回音频流
  - 支持分块传输

详细的 WebSocket 接口文档请参考：[docs/03-websocket接口设计.md](../docs/03-websocket接口设计.md)

### 使用场景对比

| 场景 | 推荐接口 | 说明 |
|------|---------|------|
| 文件上传识别 | REST API | 上传完整音频文件进行识别 |
| 实时语音识别 | WebSocket | 实时音频流，低延迟 |
| 批量文件处理 | REST API | 一次处理多个文件 |
| 实时对话系统 | WebSocket | 需要实时交互的场景 |
| 文本转语音 | REST API | 一次性合成完整音频 |
| 流式语音合成 | WebSocket | 需要实时播放的场景 |

## 更新 Swagger 文档

当修改了 API 接口或添加了新的接口时，需要重新生成 Swagger 文档：

### 1. 安装 swag 工具

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. 生成文档

```bash
swag init -g cmd/speech-server/main.go -o docs/swagger
```

### 3. 重启服务

重新启动服务后，Swagger UI 会自动加载最新的文档。

## Swagger 注释规范

### 基本注释格式

```go
// @Summary      接口摘要
// @Description  接口详细描述
// @Tags         标签（用于分组）
// @Accept       json
// @Produce      json
// @Param        param_name  param_type  data_type  required  "参数说明"
// @Success      200  {object}  ResponseType  "成功响应"
// @Failure      400  {object}  map[string]interface{}  "错误响应"
// @Router       /api/v1/path [method]
```

### 示例

#### GET 接口

```go
// GetConfig 获取配置
// @Summary      获取配置
// @Description  获取服务配置信息
// @Tags         STT
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "配置信息"
// @Router       /stt/config [get]
func (h *STTHandler) GetConfig(c *gin.Context) {
    // ...
}
```

#### POST 接口（JSON）

```go
// Synthesize 文本合成
// @Summary      文本合成
// @Description  将文本合成为语音音频
// @Tags         TTS
// @Accept       json
// @Produce      audio/wav
// @Param        request  body      SynthesizeRequest  true  "合成请求"
// @Success      200      {file}    binary            "音频文件"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /tts/synthesize [post]
func (h *TTSHandler) Synthesize(c *gin.Context) {
    // ...
}
```

#### POST 接口（文件上传）

```go
// Recognize 文件上传识别
// @Summary      文件上传识别
// @Description  上传音频文件进行语音识别
// @Tags         STT
// @Accept       multipart/form-data
// @Produce      json
// @Param        audio  formData  file  true  "音频文件"
// @Success      200    {object}  map[string]interface{}  "识别成功"
// @Router       /stt/recognize [post]
func (h *STTHandler) Recognize(c *gin.Context) {
    // ...
}
```

## 参数类型说明

- `query`: URL 查询参数
- `path`: URL 路径参数
- `formData`: 表单数据（multipart/form-data）
- `body`: JSON 请求体
- `header`: HTTP 头

## 响应类型说明

- `{object}`: JSON 对象
- `{array}`: JSON 数组
- `{file}`: 文件（二进制）
- `{string}`: 字符串
- `{integer}`: 整数

## 注意事项

1. **注释位置**: Swagger 注释必须紧邻函数定义，中间不能有空行
2. **参数顺序**: 参数注释的顺序应该与函数参数顺序一致
3. **类型引用**: 使用自定义类型时，确保类型已定义且可导出
4. **标签分组**: 使用 `@Tags` 对接口进行分组，便于在 Swagger UI 中查看

## 常见问题

### Q: Swagger UI 显示 404？

A: 确保：
1. 已正确导入 `docs/swagger` 包
2. 已添加 Swagger 路由：`ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))`
3. 已运行 `swag init` 生成文档

### Q: 如何更新文档？

A: 修改代码中的 Swagger 注释后，运行 `swag init` 重新生成文档，然后重启服务。

### Q: 如何添加认证？

A: 在 main.go 的 Swagger 注释中添加：

```go
// @securityDefinitions.apikey  ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Bearer token authentication
```

然后在接口注释中添加：

```go
// @Security    ApiKeyAuth
```

## 参考资源

- [Swagger 官方文档](https://swagger.io/docs/)
- [swaggo/swag 文档](https://github.com/swaggo/swag)
- [gin-swagger 文档](https://github.com/swaggo/gin-swagger)


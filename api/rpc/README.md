# RPC Schemas

# RPC Schema 说明

RPC 方法 schema 分为两个目录：[host-to-plugin/](host-to-plugin/) 存放宿主发送给插件的方法，[plugin-to-host/](plugin-to-host/) 存放插件调用宿主的方法。每个方法都有自己的 YAML 文件，用于定义 request/response 类型，并通过 `allOf` 和 `$ref` 与 [rpc.yaml](rpc.yaml#L3-L30) 中的基础 JSON-RPC 类型组合。

RPC method schemas are split into two directories: [host-to-plugin/](host-to-plugin/) for methods the host sends to plugins, and [plugin-to-host/](plugin-to-host/) for methods plugins call on the host. Each method gets its own YAML file defining request/response types that compose with the base JSON-RPC types in [rpc.yaml](rpc.yaml#L3-L30) using `allOf` and `$ref`.

添加新方法时，请在对应目录中创建文件，并定义 `RPCRequest{MethodName}` 和 `RPCResponse{MethodName}`。request 必须扩展 `../rpc.yaml#/$defs/JsonRpcRequest`，并设置一个 `const` method name。共享 schema 请通过类似 `../../common/events.yaml#/$defs/EventPayload` 的相对路径从 [common/](../common/) 引用。然后在 [../plugin.yaml](../plugin.yaml) 中为 request 和 response union types 添加 `$ref`，注册新 schema。添加 schema 后，运行 `task generate` 重新生成 bindings。

To add a new method, create a file in the appropriate directory with `RPCRequest{MethodName}` and `RPCResponse{MethodName}` definitions. The request must extend `../rpc.yaml#/$defs/JsonRpcRequest` and set a `const` method name. Reference shared schemas from [common/](../common/) using relative paths like `../../common/events.yaml#/$defs/EventPayload`. Then register your new schemas in [../plugin.yaml](../plugin.yaml) by adding `$ref` entries to both the request and response union types. After adding schemas, run `task generate` to regenerate bindings.

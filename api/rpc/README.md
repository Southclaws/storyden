# RPC Schemas

RPC method schemas are split into two directories: [host-to-plugin/](host-to-plugin/) for methods the host sends to plugins, and [plugin-to-host/](plugin-to-host/) for methods plugins call on the host. Each method gets its own YAML file defining request/response types that compose with the base JSON-RPC types in [rpc.yaml](rpc.yaml#L3-L30) using `allOf` and `$ref`.

To add a new method, create a file in the appropriate directory with `RPCRequest{MethodName}` and `RPCResponse{MethodName}` definitions. The request must extend `../rpc.yaml#/$defs/JsonRpcRequest` and set a `const` method name. Reference shared schemas from [common/](../common/) using relative paths like `../../common/events.yaml#/$defs/EventPayload`. Then register your new schemas in [../plugin.yaml](../plugin.yaml) by adding `$ref` entries to both the request and response union types. After adding schemas, run `task generate` to regenerate bindings.

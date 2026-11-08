# User-agent device labels

This package owns the browser and operating-system names used in passkey labels.
These names describe a client-supplied header; they must never control access,
authentication decisions, or device trust.

The tokenizer and browser/OS precedence are adapted from
[mileusna/useragent v1.3.5](https://github.com/mileusna/useragent/tree/v1.3.5).
The upstream MIT license is retained in LICENSE. `testdata/upstream.json` is
derived from that release's `ua_test.go`, retaining the header and independently
specified browser/OS expectations, without version or device-classification data.

The API intentionally contains only `Name` and `OS`. It does not retain headers,
expose versions or model names, detect bots, or guess names from arbitrary
product tokens. Known classifications return fixed strings. Unsupported clients
and operating systems return `Unknown`; upstream's arbitrary-text fallback is
not preserved. Headers over 4096 bytes, invalid UTF-8, Unicode control/format
characters or line/paragraph separators return `Unknown` for both fields.
Horizontal tabs are allowed and treated as spaces. Duplicate tokens use the
first occurrence, as in the upstream parser.

The corpus retains all 101 upstream rows with explicit browser/OS assertions.
The remaining upstream TikTok row has no OS assertion; an Android TikTok case
is tested separately. The corpus records upstream labels even for unsupported
clients. Tests explicitly identify the intentional differences: unknown products, an incomplete Android
header, and unknown operating systems. iPod recognition is an additional rule.
Product markers supplied by clients remain spoofable. Brave generally presents
as Chrome; iPads requesting desktop sites may present as macOS.

Run the compatibility, adversarial, and request-context checks with:

```sh
GOMAXPROCS=2 go test -p 1 -race ./internal/useragent ./app/services/reqinfo
GOMAXPROCS=2 go test -p 1 ./internal/useragent -run '^$' -fuzz '^FuzzParse$' -fuzztime 30s -parallel 2
```

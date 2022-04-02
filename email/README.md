# Email Templates

These templates use [MJML](https://mjml.io).

Run:

```
npm run build
```

To build .html files from the MJML templates.

They are processed by the Go template engine so for example `{{ .Key }}` will use the `Key` field in the context object passed at send-time.

Each .html template must also have a corresponding .txt for plaintext email support.

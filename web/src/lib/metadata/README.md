# Social preview images

`app/icon.tsx` and `app/opengraph-image.tsx` return the backend's finished PNG
streams. The browser icon is 32 by 32, the Apple touch icon is 180 by 180, and
the social banner is 1200 by 630. The manifest provides 192 by 192 and 512 by 512 icons for
installed web apps. Fallback PNGs are pre-sized embedded files; no request-time
resizing or permanent cache of re-encoded bytes is needed. See
[manifest icon guidance](https://web.dev/articles/add-manifest#icons).

Preserve the backend's cache policy while variants are pending. These metadata routes must remain dynamic: deployments configure the
API origin at runtime, and builds do not require a running backend.

`streamMetadataImage` uses `getAPIAddress()` so internal `SSR_API_ADDRESS`
configuration works without exposing that address in public metadata. Destinations
are fixed API endpoints, redirects are rejected, and failures are not converted
into successful blank images. No cookies or authorization are forwarded.

Composed previews can still use `ImageResponse`, as the Library preview does.
`loadMetadataImage` supplies bounded binary image data from the configured API;
it does not allow Satori to fetch arbitrary URLs or introduce another base64
representation. Compositing necessarily buffers and decodes an image. Future
fonts and layouts belong in that rendering path. If composition becomes expensive,
generate and cache the completed preview when its source changes, then serve that
artifact through the streaming path.

The common output is a PNG under 5 MB, with 1200 by 630 social image metadata and
an explicit `summary_large_image` Twitter card. The Open Graph properties must
appear in the initial HTML for crawlers, with public absolute image URLs.

References:

- [Open Graph properties](https://ogp.me/)
- [LinkedIn sharing requirements](https://www.linkedin.com/help/linkedin/answer/a521928)
- [Slack link crawlers](https://api.slack.com/robots)
- [Next metadata image response contract](https://nextjs.org/docs/app/api-reference/file-conventions/metadata/opengraph-image#returns)

`tests/metadata/social-previews.spec.ts` checks the production HTML and PNG
responses with Twitter, Facebook, Discord, Slack, and LinkedIn crawler user
agents. This checks the server contract; platform caches and actual in-app
presentation require a publicly reachable deployment to verify.

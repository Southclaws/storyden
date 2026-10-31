import { expect, test } from "@playwright/test";
import { readFileSync } from "node:fs";

import { API_ADDRESS, WEB_ADDRESS } from "../../src/config";
import {
  getAdminAccessKey,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

const crawlers = [
  "Twitterbot/1.0",
  "facebookexternalhit/1.1",
  "Discordbot/2.0",
  "Slackbot-LinkExpanding 1.0 (+https://api.slack.com/robots)",
  "LinkedInBot/1.0",
];

function metadata(html: string) {
  const head = html.split("</head>")[0]!;
  return new Map(
    [
      ...head.matchAll(/<meta (?:property|name)="([^"]+)" content="([^"]*)"/g),
    ].map((match) => [match[1]!, match[2]!.replaceAll("&amp;", "&")]),
  );
}

function expectPNG(body: Buffer, width: number, height: number) {
  expect(body.subarray(0, 8)).toEqual(
    Buffer.from([137, 80, 78, 71, 13, 10, 26, 10]),
  );
  expect(body.readUInt32BE(16)).toBe(width);
  expect(body.readUInt32BE(20)).toBe(height);
  expect(body.length).toBeLessThan(5_000_000);
}

for (const userAgent of crawlers) {
  test(`serves public social metadata and images to ${userAgent}`, async ({
    request,
  }) => {
    const response = await request.get("/", {
      headers: { "User-Agent": userAgent, Range: "bytes=0-65535" },
    });
    expect(response.ok()).toBe(true);
    const tags = metadata(await response.text());
    for (const key of [
      "og:title",
      "og:type",
      "og:url",
      "og:description",
      "og:image",
      "twitter:image",
    ]) {
      expect(tags.get(key), key).toBeTruthy();
    }
    expect(tags.get("twitter:card")).toBe("summary_large_image");
    expect(tags.get("twitter:image")).toBe(tags.get("og:image"));
    expect(tags.get("og:image:width")).toBe("1200");
    expect(tags.get("og:image:height")).toBe("630");
    const imageURL = new URL(tags.get("og:image")!);
    expect(imageURL.origin).toBe(new URL(WEB_ADDRESS).origin);
    const image = await request.get(imageURL.toString(), {
      headers: { "User-Agent": userAgent },
    });
    expect(image.status()).toBe(200);
    expect(image.headers()["content-type"]).toBe("image/png");
    expectPNG(await image.body(), 1200, 630);
  });
}

test("streams backend PNGs unchanged and supports HEAD", async ({
  request,
}) => {
  for (const [route, apiPath, width, height] of [
    ["/icon", "/info/icon/32x32", 32, 32],
    ["/opengraph-image", "/info/banner", 1200, 630],
  ] as const) {
    const backend = await request.get(`${API_ADDRESS}/api${apiPath}`);
    const image = await request.get(route);
    expect(image.status()).toBe(200);
    const body = await image.body();
    expectPNG(body, width, height);
    expect(body).toEqual(await backend.body());
    expect(image.headers()["cache-control"]).toBe(
      backend.headers()["cache-control"],
    );
    const head = await request.head(route);
    expect(head.status()).toBe(200);
    expect(head.headers()["content-type"]).toBe("image/png");
    expect((await head.body()).length).toBe(0);
  }
});

test("retains the social image when a thread supplies its own metadata", async ({
  request,
}) => {
  const title = `Social preview thread ${Date.now()}`;
  const thread = await withAdminAccessKey(({ threadCreate }) =>
    threadCreate({
      title,
      body: "A public discussion preview",
      visibility: "published",
    }),
  );
  const response = await request.get(`/t/${thread.slug}`, {
    headers: { "User-Agent": "Twitterbot/1.0" },
  });
  expect(response.ok()).toBe(true);
  const tags = metadata(await response.text());
  expect(tags.get("og:title")).toBe(title);
  expect(tags.get("og:type")).toBe("article");
  expect(new URL(tags.get("og:image")!).pathname).toBe("/opengraph-image");
  expect(tags.get("twitter:image")).toBe(tags.get("og:image"));
});

test("uses small browser icons and correctly sized manifest images", async ({
  request,
}) => {
  const response = await request.get("/");
  const html = await response.text();
  const links = html.match(/<link\b[^>]*>/g) ?? [];
  const touchIcon = links.find((link) =>
    link.includes('rel="apple-touch-icon"'),
  );
  expect(touchIcon).toContain("/api/info/icon/180x180");
  expect(touchIcon).toContain('sizes="180x180"');
  expect(
    links.some(
      (link) => link.includes('rel="icon"') && link.includes('sizes="32x32"'),
    ),
  ).toBe(true);

  const manifest = await (await request.get("/manifest.json")).json();
  const icons = new Map<string, string>(
    manifest.icons.map((icon: { src: string; sizes: string }) => [
      icon.src,
      icon.sizes,
    ]),
  );
  expect([...icons.values()]).toContain("512x512");
  expect([...icons.values()]).toContain("192x192");
  for (const [src, sizes] of icons) {
    const [width, height] = sizes.split("x").map(Number);
    const image = await request.get(src);
    expect(image.ok()).toBe(true);
    expectPNG(await image.body(), width!, height!);
  }
});

test("renders public Library previews with and without an uploaded image", async ({
  request,
}) => {
  const upload = await request.post(
    `${API_ADDRESS}/api/assets?filename=social-preview.png`,
    {
      headers: {
        Authorization: `Bearer ${getAdminAccessKey()}`,
        "Content-Type": "application/octet-stream",
      },
      data: readFileSync(
        "../app/services/branding/icon/embed/default_icon.png",
      ),
    },
  );
  expect(upload.ok()).toBe(true);
  const asset = await upload.json();
  for (const withImage of [false, true]) {
    const slug = `social-preview-${Date.now()}-${withImage}`;
    await withAdminAccessKey(({ nodeCreate }) =>
      nodeCreate({
        name: `Social preview ${withImage}`,
        slug,
        description: "A publicly shared Library page",
        visibility: "published",
        ...(withImage ? { primary_image_asset_id: asset.id } : {}),
      }),
    );
    const page = await request.get(`/l/${slug}`, {
      headers: { "User-Agent": "Twitterbot/1.0" },
    });
    const tags = metadata(await page.text());
    const imageURL = tags.get("og:image")!;
    expect(new URL(imageURL).pathname).toBe("/l/og");
    expect(tags.get("twitter:image")).toBe(imageURL);
    const preview = await request.get(imageURL);
    expect(preview.status()).toBe(200);
    expectPNG(await preview.body(), 1200, 630);
  }
});

import { readdir } from "fs/promises";
import { filter, flow, map } from "lodash/fp";
import { MetadataRoute } from "next";

import { Dirent } from "fs";
import path from "path";

const BLOG_PATH = `content/blog/`;
const DOCS_PATH = `content/docs/`;
const BLOG_ZH_PATH = `content/zh/blog/`;
const DOCS_ZH_PATH = `content/zh/docs/`;

type ChangeFrequency =
  | "always"
  | "hourly"
  | "daily"
  | "weekly"
  | "monthly"
  | "yearly"
  | "never";

type SitemapFile = {
  url: string;
  lastModified?: string | Date;
  changeFrequency?: ChangeFrequency;
  priority?: number;
};

type ItemSettings = {
  sub: string;
  changeFrequency: ChangeFrequency;
  priority: number;
};

const base = "https://www.storyden.org";

const lastModified = new Date();

const onlyMDX = filter<string>((v) => v.endsWith(".mdx"));

const stripIndex = filter<string>((v) => !v.endsWith("index"));

const onlyFiles = filter<Dirent>((v) => v.isFile());

const toFilePath = map<Dirent, string>((v) => path.join(v.parentPath, v.name));

const stripLeading = map<string, string>((v) => v.replaceAll(DOCS_PATH, ""));
const stripLeadingZh = map<string, string>((v) =>
  v.replaceAll(DOCS_ZH_PATH, ""),
);

const removeExtension = map<string, string>((v) => v.replaceAll(".mdx", ""));

const toSitemapItem = ({ sub, changeFrequency, priority }: ItemSettings) =>
  map<string, SitemapFile>((v) => ({
    url: [base, sub, v].join("/"),
    lastModified,
    changeFrequency,
    priority,
  }));

const processPaths = (s: ItemSettings) =>
  flow(onlyMDX, removeExtension, stripIndex, toSitemapItem(s));

const processBlogPaths = processPaths({
  sub: "blog",
  changeFrequency: "yearly",
  priority: 0.9,
});

const processBlogZhPaths = processPaths({
  sub: "zh/blog",
  changeFrequency: "yearly",
  priority: 0.8,
});

const processDocsPaths = flow(
  onlyFiles,
  toFilePath,
  stripLeading,
  processPaths({
    sub: "docs",
    changeFrequency: "monthly",
    priority: 0.5,
  })
);

const processDocsZhPaths = flow(
  onlyFiles,
  toFilePath,
  stripLeadingZh,
  processPaths({
    sub: "zh/docs",
    changeFrequency: "monthly",
    priority: 0.4,
  })
);

export default async function SiteMap(): Promise<MetadataRoute.Sitemap> {
  const blogPages = await getBlogPages();
  const docsPages = await getDocsPages();
  const blogZhPages = await getBlogZhPages();
  const docsZhPages = await getDocsZhPages();

  return [
    {
      url: base,
      lastModified: new Date(),
      changeFrequency: "monthly",
      priority: 1,
    },
    {
      url: [base, "zh"].join("/"),
      lastModified: new Date(),
      changeFrequency: "monthly",
      priority: 0.9,
    },
    ...blogPages,
    ...docsPages,
    ...blogZhPages,
    ...docsZhPages,
  ];
}

async function getBlogPages(): Promise<MetadataRoute.Sitemap> {
  const dir = await readdir(BLOG_PATH);

  const items = processBlogPaths(dir);

  return items;
}

async function getDocsPages(): Promise<MetadataRoute.Sitemap> {
  const dir = await readdir(DOCS_PATH, {
    withFileTypes: true,
    recursive: true,
  });

  const items = processDocsPaths(dir);

  return items;
}

async function getBlogZhPages(): Promise<MetadataRoute.Sitemap> {
  const dir = await readdir(BLOG_ZH_PATH);

  const items = processBlogZhPaths(dir);

  return items;
}

async function getDocsZhPages(): Promise<MetadataRoute.Sitemap> {
  const dir = await readdir(DOCS_ZH_PATH, {
    withFileTypes: true,
    recursive: true,
  });

  const items = processDocsZhPaths(dir);

  return items;
}

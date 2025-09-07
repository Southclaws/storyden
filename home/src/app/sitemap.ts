import { readdir } from "fs/promises";
import { filter, flow, map } from "lodash/fp";
import { MetadataRoute } from "next";

import { Dirent } from "fs";
import path from "path";

const BLOG_PATH = `content/blog/`;
const DOCS_PATH = `content/docs/`;

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

export default async function SiteMap(): Promise<MetadataRoute.Sitemap> {
  const blogPages = await getBlogPages();
  const docsPages = await getDocsPages();

  return [
    {
      url: base,
      lastModified: new Date(),
      changeFrequency: "monthly",
      priority: 1,
    },
    ...blogPages,
    ...docsPages,
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

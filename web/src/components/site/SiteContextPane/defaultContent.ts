import type { Locale } from "@/i18n/config";

const defaultDescriptions = new Set([
  "A forum for the modern age",
  "A forum for the modern age.",
]);

const defaultContentBody = `
<p>Welcome to your new community!</p>
<p>You can edit this content by clicking Edit below.</p>
<p>This is a <em>rich text section</em> for telling visitors what your community is about.</p>
<p>Add a link to your <a href="https://discord.gg/XF6ZBGF9XF">Discord</a> or other sites.</p>
<p>Enjoy!</p>`;

const defaultContents = new Set([
  normalizeHtml(defaultContentBody),
  normalizeHtml(`<body>${defaultContentBody}</body>`),
]);

const defaultContentPhrases = [
  "Welcome to your new community!",
  "You can edit this content by clicking Edit below.",
  "This is a rich text section for telling visitors what your community is about.",
  "Add a link to your",
  "or other sites.",
  "Enjoy!",
];

const defaultContentZh = `
<p>欢迎来到你的新社区！</p>
<p>点击下方“编辑”即可修改这段内容。</p>
<p>这里是一段<em>富文本区域</em>，可以用来向访客介绍你的社区。</p>
<p>你也可以添加指向 <a href="https://discord.gg/XF6ZBGF9XF">Discord</a> 或其他站点的链接。</p>
<p>玩得开心！</p>`;

function normalizeHtml(value: string) {
  return value
    .trim()
    .replace(/>\s+</g, "><")
    .replace(/\s+/g, " ");
}

function normalizeText(value: string) {
  return value
    .replace(/<[^>]+>/g, " ")
    .replace(/&nbsp;/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/\s+/g, " ")
    .trim();
}

function isDefaultContent(value: string) {
  const html = normalizeHtml(value);

  if (defaultContents.has(html)) {
    return true;
  }

  const text = normalizeText(value);

  return defaultContentPhrases.every((phrase) => text.includes(phrase));
}

export function getDisplayDescription(locale: Locale, value: string) {
  if (locale === "zh" && defaultDescriptions.has(value.trim())) {
    return "面向现代社区的论坛。";
  }

  return value;
}

export function getDisplayContent(locale: Locale, value?: string) {
  if (locale === "zh" && value && isDefaultContent(value)) {
    return defaultContentZh;
  }

  return value;
}

"use client";

import { useEffect } from "react";

import { Locale } from "@/i18n/config";
import { useI18n } from "@/i18n/provider";

const translations = [
  ["Issues", "问题"],
  ["Route", "路由"],
  ["Static", "静态"],
  ["Dynamic", "动态"],
  ["Loading...", "正在加载..."],
  ["Bundler", "打包器"],
  ["Cache Components", "缓存组件"],
  ["Enabled", "已启用"],
  ["Instant Navs", "即时导航"],
  ["Route Info", "路由信息"],
  ["Static Route", "静态路由"],
  ["Dynamic Route", "动态路由"],
  ["Preferences", "偏好设置"],
  ["Theme", "主题"],
  ["Select your theme preference.", "选择主题偏好。"],
  ["System", "系统"],
  ["Light", "浅色"],
  ["Dark", "深色"],
  ["Position", "位置"],
  ["Adjust the placement of your dev tools.", "调整开发工具的位置。"],
  ["Bottom Left", "左下角"],
  ["Bottom Right", "右下角"],
  ["Top Left", "左上角"],
  ["Top Right", "右上角"],
  ["Size", "大小"],
  ["Adjust the size of your dev tools.", "调整开发工具的大小。"],
  ["Small", "小"],
  ["Medium", "中"],
  ["Large", "大"],
  ["Hide Dev Tools for this session", "本次会话隐藏开发工具"],
  [
    "Hide Dev Tools until you restart your dev server, or 1 day.",
    "隐藏开发工具，直到重启开发服务器或 1 天后。",
  ],
  ["Hide", "隐藏"],
  ["Hide Dev Tools shortcut", "开发工具隐藏快捷键"],
  [
    "Set a custom keyboard shortcut to toggle visibility.",
    "设置用于切换显示/隐藏的自定义快捷键。",
  ],
  ["Record Shortcut", "录制快捷键"],
  ["Clear shortcut", "清除快捷键"],
  ["Shortcut set", "快捷键已设置"],
  ["Recording", "正在录制"],
  ["Disable Dev Tools for this project", "为此项目禁用开发工具"],
  ["To disable this UI completely, set", "要完全禁用这个界面，请设置"],
  ["in your", "，位置："],
  ["file.", "文件。"],
  ["Restart Dev Server", "重启开发服务器"],
  [
    "Restarts the development server without needing to leave the browser.",
    "不离开浏览器即可重启开发服务器。",
  ],
  ["Restart", "重启"],
  ["Reset Bundler Cache", "重置打包缓存"],
  [
    "Clears the bundler cache and restarts the dev server. Helpful if you are seeing stale errors or changes are not appearing.",
    "清除打包缓存并重启开发服务器；如果看到旧错误或改动没有生效，这会有帮助。",
  ],
  ["Reset Cache", "重置缓存"],
  ["Learn More", "了解更多"],
  ["Clear Segment Overrides", "清除片段覆盖"],
  ["The path", "路径"],
  [
    'is marked as "static" since it will be prerendered during the build time.',
    "被标记为“静态”，因为它会在构建时预渲染。",
  ],
  [
    "With Static Rendering, routes are rendered at build time, or in the background after",
    "使用静态渲染时，路由会在构建时渲染，或在",
  ],
  ["data revalidation", "数据重新验证"],
  [
    "Static rendering is useful when a route has data that is not personalized to the user and can be known at build time, such as a static blog post or a product page.",
    "当路由的数据不因用户而异，并且能在构建时确定时，静态渲染很有用，例如静态博客文章或产品页面。",
  ],
  [
    'is marked as "dynamic" since it will be rendered for each user at',
    "被标记为“动态”，因为它会为每个用户在",
  ],
  ["request time", "请求时"],
  [
    "Dynamic rendering is useful when a route has data that is personalized to the user or has information that can only be known at request time, such as cookies or the URL's search params.",
    "当路由包含针对用户个性化的数据，或包含只能在请求时得知的信息时，动态渲染很有用，例如 cookie 或 URL 的搜索参数。",
  ],
  ["During rendering, if a", "渲染期间，如果发现"],
  ["or a", "或"],
  ["option of", "选项"],
  [
    "is discovered, Next.js will switch to dynamically rendering the whole route.",
    "，Next.js 会切换为动态渲染整个路由。",
  ],
  [
    "Exporting the",
    "导出",
  ],
  [
    "function will opt the route into dynamic rendering. This function will be called by the server on every request.",
    "函数会让该路由进入动态渲染。服务器会在每次请求时调用这个函数。",
  ],
] as const;

function translateText(value: string, locale: Locale) {
  const trimmed = value.trim();

  if (!trimmed) {
    return value;
  }

  const match = translations.find(
    ([english, chinese]) => english === trimmed || chinese === trimmed,
  );

  if (!match) {
    return value;
  }

  const target = locale === "zh" ? match[1] : match[0];

  return value.replace(trimmed, target);
}

function translateElementAttribute(
  element: Element,
  attribute: "aria-label" | "title",
  locale: Locale,
) {
  const value = element.getAttribute(attribute);

  if (!value) {
    return;
  }

  const translated = translateText(value, locale);

  if (translated !== value) {
    element.setAttribute(attribute, translated);
  }
}

function translateScope(scope: Element, locale: Locale) {
  const walker = document.createTreeWalker(scope, NodeFilter.SHOW_TEXT);

  for (
    let node = walker.nextNode();
    node;
    node = walker.nextNode()
  ) {
    const translated = translateText(node.textContent ?? "", locale);

    if (translated !== node.textContent) {
      node.textContent = translated;
    }
  }

  for (const element of scope.querySelectorAll("[aria-label], [title]")) {
    translateElementAttribute(element, "aria-label", locale);
    translateElementAttribute(element, "title", locale);
  }
}

function getDevtoolsShadowRoots() {
  return Array.from(document.querySelectorAll("nextjs-portal"))
    .map((portal) => portal.shadowRoot)
    .filter((root): root is ShadowRoot => root !== null);
}

function getTranslationScopes(root: ShadowRoot) {
  const scopes = new Set<Element>();

  for (const panel of root.querySelectorAll("#panel-route")) {
    scopes.add(panel);
  }

  for (const menu of root.querySelectorAll(".dev-tools-indicator-menu")) {
    scopes.add(menu);
  }

  return scopes;
}

function applyDevtoolsLocale(locale: Locale) {
  for (const root of getDevtoolsShadowRoots()) {
    for (const scope of getTranslationScopes(root)) {
      translateScope(scope, locale);
    }
  }
}

export function NextDevtoolsI18n() {
  const { locale } = useI18n();

  useEffect(() => {
    if (process.env.NODE_ENV !== "development") {
      return;
    }

    const observedRoots = new WeakSet<ShadowRoot>();
    let frame: number | undefined;

    const observeShadowRoots = () => {
      for (const root of getDevtoolsShadowRoots()) {
        if (observedRoots.has(root)) {
          continue;
        }

        observedRoots.add(root);
        observer.observe(root, {
          attributes: true,
          attributeFilter: ["aria-label", "title"],
          characterData: true,
          childList: true,
          subtree: true,
        });
      }
    };

    const schedule = () => {
      if (frame !== undefined) {
        cancelAnimationFrame(frame);
      }

      frame = requestAnimationFrame(() => {
        observeShadowRoots();
        applyDevtoolsLocale(locale);
      });
    };

    const observer = new MutationObserver(schedule);

    observer.observe(document.body, {
      childList: true,
      subtree: true,
    });

    schedule();

    return () => {
      if (frame !== undefined) {
        cancelAnimationFrame(frame);
      }

      observer.disconnect();
    };
  }, [locale]);

  return null;
}

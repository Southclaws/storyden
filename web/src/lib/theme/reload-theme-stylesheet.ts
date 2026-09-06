const THEME_STYLESHEET_SELECTOR = 'link[data-sd-theme-asset="stylesheet"]';

export function reloadThemeStylesheet() {
  const current = document.querySelector<HTMLLinkElement>(
    THEME_STYLESHEET_SELECTOR,
  );
  if (!current) {
    return Promise.reject(new Error("The theme stylesheet is not mounted."));
  }

  const replacement = current.cloneNode() as HTMLLinkElement;
  const href = new URL(current.href, document.baseURI);
  href.searchParams.set("sd-theme-editor", Date.now().toString());
  replacement.href = href.href;

  return new Promise<void>((resolve, reject) => {
    replacement.addEventListener(
      "load",
      () => {
        current.remove();
        document
          .querySelectorAll("style[data-sd-theme-editor-live]")
          .forEach((element) => element.remove());
        resolve();
      },
      { once: true },
    );
    replacement.addEventListener(
      "error",
      () => {
        replacement.remove();
        reject(new Error("The updated theme stylesheet could not be loaded."));
      },
      { once: true },
    );
    current.after(replacement);
  });
}

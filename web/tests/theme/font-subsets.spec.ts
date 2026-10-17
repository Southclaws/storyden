import { expect, test } from "@playwright/test";
import { readFileSync } from "node:fs";

test("ordinary browsing downloads Latin fonts and preloads only regular Latin", async ({
  page,
}) => {
  await page.goto("/");
  await page.evaluate(() => document.fonts.ready);

  const fonts = await page.evaluate(() =>
    performance
      .getEntriesByType("resource")
      .filter((entry) => entry.name.endsWith(".woff2"))
      .map((entry) => ({
        url: entry.name,
        bytes: (entry as PerformanceResourceTiming).encodedBodySize,
      })),
  );
  expect(fonts.length).toBeGreaterThan(0);
  expect(fonts.every((font) => font.bytes > 0)).toBe(true);
  expect(
    fonts.every((font) => /[_-]latin[.-]/.test(font.url)),
    JSON.stringify(fonts),
  ).toBe(true);
  expect(fonts.reduce((bytes, font) => bytes + font.bytes, 0)).toBeLessThan(
    400_000,
  );
  const preloads = page.locator('link[rel="preload"][as="font"]');
  await expect(preloads).toHaveCount(1);
  await expect(preloads).toHaveAttribute("href", /InterVariable_latin[.-]/);
});

for (const { language, text, subset } of [
  { language: "Russian", text: "Привет мир", subset: "cyrillic" },
  { language: "Greek", text: "Καλημέρα κόσμε", subset: "greek" },
  { language: "Vietnamese", text: "Tiếng Việt", subset: "latin-ext" },
  { language: "CJK", text: "中文日本語한국어", subset: undefined },
  { language: "Thai", text: "สวัสดีชาวโลก", subset: undefined },
]) {
  test(`${language} text loads only its additional Inter subset`, async ({
    page,
  }) => {
    await page.goto("/");
    await page.evaluate(() => document.fonts.ready);
    const before = await page.evaluate(() =>
      performance
        .getEntriesByType("resource")
        .filter((entry) => entry.name.endsWith(".woff2"))
        .map((entry) => entry.name),
    );

    await page.evaluate(async (text) => {
      const paragraph = document.createElement("p");
      paragraph.textContent = text;
      paragraph.style.cssText = 'font: 400 24px "inter", sans-serif';
      document.body.append(paragraph);
      await document.fonts.load('400 24px "inter"', text);
    }, text);

    const additional = await page.evaluate(
      (before) =>
        performance
          .getEntriesByType("resource")
          .filter(
            (entry) =>
              entry.name.endsWith(".woff2") && !before.includes(entry.name),
          )
          .map((entry) => entry.name),
      before,
    );
    if (subset) {
      expect(additional).toHaveLength(1);
      expect(additional[0]).toContain(`InterVariable-${subset}.`);
    } else {
      expect(additional).toEqual([]);
    }
  });
}

test("subsets preserve the original glyphs across scripts, weights and styles", async ({
  page,
}) => {
  const original = [
    readFileSync("src/app/fonts/static/InterVariable.woff2", "base64"),
    readFileSync("src/app/fonts/static/InterVariable-Italic.woff2", "base64"),
  ];
  await page.goto("/");
  const differences = await page.evaluate(async (original) => {
    for (const [index, data] of original.entries()) {
      const face = new FontFace(
        "InterOriginal",
        `url(data:font/woff2;base64,${data})`,
        { weight: "100 900", style: index ? "italic" : "normal" },
      );
      document.fonts.add(await face.load());
    }
    const samples = [
      "Storyden — café, naïve, fi fl ffi 0123456789",
      "Привет мир — «Тест» Й й е\u0301 и\u0306",
      "Українська мова — Її Єє Ґґ",
      "Καλημέρα κόσμε — ἄνθρωπος α\u0301",
      "Tiếng Việt — a\u0306\u0301 o\u031B\u0301",
      "中文日本語한국어 สวัสดีชาวโลก",
      "Latin Русский Ελληνικά 中文 ไทย",
      "→ ← ↑ ↓ • € £ ¥ © ™ ∑ ∞",
    ];
    const canvas = document.createElement("canvas");
    canvas.width = 1200;
    canvas.height = 80;
    const context = canvas.getContext("2d", { willReadFrequently: true })!;
    // Kerning cannot cross font-file boundaries. Compare the glyphs themselves,
    // including combining accents and the system fallback for CJK and Thai.
    context.fontKerning = "none";
    const differences: { style: string; weight: number; text: string }[] = [];
    for (const style of ["normal", "italic"]) {
      for (const weight of [100, 200, 300, 400, 500, 600, 700, 800, 900]) {
        for (const text of samples) {
          await document.fonts.load(`${style} ${weight} 24px "inter"`, text);
          const images = ["InterOriginal", "inter"].map((family) => {
            context.clearRect(0, 0, canvas.width, canvas.height);
            context.font = `${style} ${weight} 24px "${family}", sans-serif`;
            context.fillText(text, 10, 40);
            return context.getImageData(0, 0, canvas.width, canvas.height).data;
          });
          if (images[0]!.some((value, index) => value !== images[1]![index])) {
            differences.push({ style, weight, text });
          }
        }
      }
    }
    return differences;
  }, original);
  expect(differences).toEqual([]);
});

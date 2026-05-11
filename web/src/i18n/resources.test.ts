import { test } from "uvu";
import * as assert from "uvu/assert";

import { locales } from "./config";
import { messages } from "./resources";

const hanCharacters = /[\u3400-\u9fff]/;
const staticTranslationCall =
  /\bt(?:Server)?\(\s*(["'`])((?:(?!\1).)+)\1/g;

test("all locale dictionaries expose the same message keys", () => {
  const [defaultLocale, ...otherLocales] = locales;
  const expected = Object.keys(messages[defaultLocale]).sort();

  for (const locale of otherLocales) {
    assert.equal(Object.keys(messages[locale]).sort(), expected, locale);
  }
});

test("English dictionary values do not contain Chinese text", () => {
  const mixed = Object.entries(messages.en).filter(([, value]) =>
    hanCharacters.test(value),
  );

  assert.equal(mixed, []);
});

test("static translation calls exist in every dictionary", async () => {
  const { readdir, readFile } = await import("node:fs/promises");
  const { join } = await import("node:path");

  async function* files(dir: string): AsyncGenerator<string> {
    for (const entry of await readdir(dir, { withFileTypes: true })) {
      const path = join(dir, entry.name);
      if (entry.isDirectory()) {
        yield* files(path);
      } else if (/\.(ts|tsx)$/.test(entry.name) && !entry.name.endsWith(".test.ts")) {
        yield path;
      }
    }
  }

  const keys = new Set<string>();
  for await (const file of files("src")) {
    const source = await readFile(file, "utf8");
    for (const match of source.matchAll(staticTranslationCall)) {
      const key = match[2];
      if (key !== undefined && !key.includes("\\")) {
        keys.add(key);
      }
    }
  }

  for (const locale of locales) {
    const missing = [...keys].filter((key) => messages[locale][key] === undefined);
    assert.equal(missing, [], locale);
  }
});

test.run();

import { escape } from "lodash";
import assert from "node:assert/strict";
import { test } from "vitest";

import { canonicalPage, jsonLD, pageLinks, socialImage } from "./public";

test("social images advertise the public route and actual PNG dimensions", () => {
  assert.deepEqual(socialImage(), {
    url: "http://localhost:3000/opengraph-image",
    width: 1200,
    height: 630,
    type: "image/png",
    alt: "Community preview",
  });
  assert.equal(
    socialImage("/l/og?slug=hello&t=one", "A page").url,
    "http://localhost:3000/l/og?slug=hello&t=one",
  );
});

test("canonical and pagination links use the public frontend origin", () => {
  assert.equal(canonicalPage("/t/hello", 1), "http://localhost:3000/t/hello");
  assert.equal(
    canonicalPage("/t/hello", 2),
    "http://localhost:3000/t/hello?page=2",
  );
  assert.deepEqual(pageLinks("/t/hello", 1, 1), {});
  assert.deepEqual(pageLinks("/t/hello", 1, 3), {
    next: "http://localhost:3000/t/hello?page=2",
  });
  assert.deepEqual(pageLinks("/t/hello", 2, 3), {
    previous: "http://localhost:3000/t/hello",
    next: "http://localhost:3000/t/hello?page=3",
  });
  assert.deepEqual(pageLinks("/t/hello", 3, 3), {
    previous: "http://localhost:3000/t/hello?page=2",
  });
});

test("JSON-LD escapes script delimiters without changing the content", () => {
  const data = {
    text: 'A </script><script>alert("x")</script> post & café\u2028\u2029',
  };
  const encoded = jsonLD(data);

  assert.equal(encoded.includes("<"), false);
  assert.equal(encoded.includes(">"), false);
  assert.equal(encoded.includes("&"), false);
  assert.equal(encoded.includes("\u2028"), false);
  assert.equal(encoded.includes("\u2029"), false);
  assert.deepEqual(JSON.parse(encoded), {
    text: 'A </script><script>alert("x")</script> post & café\u2028\u2029',
  });

  // Script elements are raw text: HTML entities would make JSON invalid.
  assert.throws(() => JSON.parse(escape(JSON.stringify(data))));
});

import { test } from "uvu";
import * as assert from "uvu/assert";

import { getDisplayContent, getDisplayDescription } from "./defaultContent";

const starterContentWithExtraAttributes = `
<div>
  <p>Welcome to your new community!</p>
  <p>You can edit this content by clicking Edit below.</p>
  <p>This is a <em data-type="italic">rich text section</em> for telling visitors what your community is about.</p>
  <p>Add a link to your <a href="https://discord.gg/XF6ZBGF9XF" target="_blank">Discord</a> or other sites.</p>
  <p>Enjoy!</p>
</div>`;

test("localises default site description in Chinese only", () => {
  assert.equal(
    getDisplayDescription("zh", "A forum for the modern age."),
    "面向现代社区的论坛。",
  );
  assert.equal(
    getDisplayDescription("en", "A forum for the modern age."),
    "A forum for the modern age.",
  );
});

test("localises starter site content despite harmless html changes", () => {
  const content = getDisplayContent("zh", starterContentWithExtraAttributes);

  assert.ok(content?.includes("欢迎来到你的新社区！"));
  assert.ok(!content?.includes("Welcome to your new community!"));
});

test("does not localise custom site content", () => {
  const custom = "<p>Welcome to your new community! But this is custom.</p>";

  assert.equal(getDisplayContent("zh", custom), custom);
});

test.run();

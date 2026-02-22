import { test } from "uvu";
import * as assert from "uvu/assert";

import { createEmitter } from "./emitter";

const PING = Symbol("ping");

type Events = {
  message: string;
  count: number;
  [PING]: { ok: boolean };
};

test("emits payloads to listeners of the same event only", () => {
  const emitter = createEmitter<Events>();
  const received: string[] = [];

  emitter.on("message", (v) => received.push(v));
  emitter.on("count", (v) => received.push(`count:${v}`));

  emitter.emit("message", "hello");

  assert.equal(received, ["hello"]);
});

test("off unsubscribes a listener", () => {
  const emitter = createEmitter<Events>();
  const received: string[] = [];
  const handler = (v: string) => received.push(v);

  emitter.on("message", handler);
  emitter.emit("message", "one");
  emitter.off("message", handler);
  emitter.emit("message", "two");

  assert.equal(received, ["one"]);
});

test("adding the same handler twice does not duplicate calls", () => {
  const emitter = createEmitter<Events>();
  const received: string[] = [];
  const handler = (v: string) => received.push(v);

  emitter.on("message", handler);
  emitter.on("message", handler);
  emitter.emit("message", "hello");

  assert.equal(received, ["hello"]);
});

test("supports symbol event keys", () => {
  const emitter = createEmitter<Events>();
  let ok = false;

  emitter.on(PING, (event) => {
    ok = event.ok;
  });
  emitter.emit(PING, { ok: true });

  assert.ok(ok);
});

test("listener added during emit does not run until next emit", () => {
  const emitter = createEmitter<Events>();
  const calls: string[] = [];
  const late = (v: string) => calls.push(`late:${v}`);

  emitter.on("message", (v) => {
    calls.push(`first:${v}`);
    emitter.on("message", late);
  });

  emitter.emit("message", "one");
  emitter.emit("message", "two");

  assert.equal(calls, ["first:one", "first:two", "late:two"]);
});

test("listener removed during emit does not run later in same emit", () => {
  const emitter = createEmitter<Events>();
  const calls: string[] = [];
  const second = (v: string) => calls.push(`second:${v}`);

  emitter.on("message", (v) => {
    calls.push(`first:${v}`);
    emitter.off("message", second);
  });
  emitter.on("message", second);

  emitter.emit("message", "one");

  assert.equal(calls, ["first:one"]);
});

test.run();

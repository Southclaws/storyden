import { parseDate } from "@internationalized/date";
import { execFileSync } from "node:child_process";
import { describe, expect, it } from "vitest";

import { formatISODate, parseISODate } from "./date-picker-iso";

describe("date-picker-iso", () => {
  describe("parseISODate", () => {
    it("accepts canonical ISO dates", () => {
      expect(parseISODate("2026-03-03")?.toString()).toBe("2026-03-03");
      expect(parseISODate("2024-02-29")?.toString()).toBe("2024-02-29");
    });

    it("rejects invalid or non-canonical values", () => {
      expect(parseISODate("2026-02-29")).toBeUndefined();
      expect(parseISODate("03-03-2026")).toBeUndefined();
      expect(parseISODate("2026/03/03")).toBeUndefined();
      expect(parseISODate("2026-03-03T00:00:00Z")).toBeUndefined();
      expect(parseISODate("2026-3-3")).toBeUndefined();
      expect(parseISODate(" 2026-03-03 ")).toBeUndefined();
      expect(parseISODate("2026-13-03")).toBeUndefined();
      expect(parseISODate("2026-03-00")).toBeUndefined();
    });
  });

  describe("formatISODate", () => {
    it("formats date values as YYYY-MM-DD", () => {
      expect(formatISODate(parseDate("2026-03-03"))).toBe("2026-03-03");
    });

    it("formats consistently across runtime timezones", () => {
      const script = `
import { parseDate } from "@internationalized/date";
import mod from "./src/components/ui/date-picker-iso.ts";
const { formatISODate } = mod;
process.stdout.write(formatISODate(parseDate("2026-03-06")));
`;

      const americaChicago = execFileSync(
        process.execPath,
        ["--import", "tsx", "-e", script],
        {
          cwd: process.cwd(),
          env: { ...process.env, TZ: "America/Chicago" },
          encoding: "utf8",
        },
      ).trim();

      const pacificKiritimati = execFileSync(
        process.execPath,
        ["--import", "tsx", "-e", script],
        {
          cwd: process.cwd(),
          env: { ...process.env, TZ: "Pacific/Kiritimati" },
          encoding: "utf8",
        },
      ).trim();

      expect(americaChicago).toBe("2026-03-06");
      expect(pacificKiritimati).toBe("2026-03-06");
    });
  });
});

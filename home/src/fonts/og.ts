import { readFile } from "fs/promises";
import { join } from "path";

export const joeiBold = async () => {
  const full = join(
    process.cwd(),
    "src",
    "fonts",
    "static",
    "JoieGrotesk-Bold.otf"
  );

  const f = await readFile(full);

  return f;
};

export const workSans = async () => {
  const full = join(
    process.cwd(),
    "src",
    "fonts",
    "static",
    "WorkSans-Medium.otf"
  );

  const f = await readFile(full);

  return f;
};

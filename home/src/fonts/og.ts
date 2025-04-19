import { readFile } from "fs/promises";
import { dirname, join } from "path";

export const joeiBold = async () => {
  const rel = "./static/JoieGrotesk-Bold.otf";
  const dirUrl = new URL(import.meta.url);
  const dir = dirname(dirUrl.pathname);

  const full = join(dir, rel);

  const f = await readFile(full);

  return f.buffer;
};

// export const workSans = () =>
//   fetch(new URL("./static/WorkSans-Medium.woff2", import.meta.url)).then(
//     (res) => res.arrayBuffer()
//   );

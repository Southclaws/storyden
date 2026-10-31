const TRANSPARENT_PIXEL =
  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAAC0lEQVR4nGNgAAIAAAUAAXpeqz8AAAAASUVORK5CYII=";

export async function fetchImageDataURI(url: string): Promise<string> {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      return TRANSPARENT_PIXEL;
    }

    const contentType = response.headers.get("content-type") ?? "image/png";
    const buffer = Buffer.from(await response.arrayBuffer());

    return `data:${contentType};base64,${buffer.toString("base64")}`;
  } catch {
    return TRANSPARENT_PIXEL;
  }
}

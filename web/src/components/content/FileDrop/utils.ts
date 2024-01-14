export function isSupportedImage(mime: string): boolean {
  const category = mime.split("/")[0];

  switch (category) {
    case "image":
      return true;
  }

  return false;
}

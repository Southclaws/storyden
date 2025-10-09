export function scrollToBottom(after: number = 100) {
  setTimeout(
    () =>
      window.scrollTo({
        behavior: "smooth",
        top: document.body.scrollHeight,
      }),
    after,
  );
}

export function scrollToTop(after: number = 0) {
  setTimeout(
    () =>
      window.scrollTo({
        behavior: "smooth",
        top: 0,
      }),
    after,
  );
}
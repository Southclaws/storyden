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

// Loaded only by the E2E harness. Reserved documentation addresses select
// deterministic failures without changing the shared backend's state.
const originalFetch = globalThis.fetch;

globalThis.fetch = async function (input, init) {
  const request = input instanceof Request ? input : new Request(input, init);
  const path = new URL(request.url).pathname;
  const clientAddress = request.headers.get("x-forwarded-for");
  const publicRead = !request.headers.has("cookie");
  const metadataEndpoint =
    path.startsWith("/api/threads") || path.startsWith("/api/nodes/");

  if (publicRead && metadataEndpoint) {
    let status;

    if (clientAddress === "198.51.100.86") {
      status = 429;
    } else if (
      clientAddress === "198.51.100.87" ||
      clientAddress === "198.51.100.89"
    ) {
      status = 500;
    }

    if (status) {
      // Model a slow upstream so a loading boundary has time to flush first.
      await new Promise((resolve) => setTimeout(resolve, 250));

      return Response.json(
        { trace_id: "metadata-e2e", title: "Injected metadata failure" },
        { status },
      );
    }
  }

  return originalFetch(input, init);
};

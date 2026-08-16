package robot

// InvocationContext is optional, client-defined prompt context captured when
// an input is submitted. It is deliberately transport- and client-agnostic:
// callers may describe a Storyden surface, a mobile application state, an
// integration, or any other invocation environment.
//
// This is useful for handling contextual instructions, such as if a member just
// expects a robot to "know" what they are seeing, for example if a member is on
// a thread and says "update the title of this thread" or "delete this comment".
// this allows a robot to receive the surrounding context to complete a request.
type InvocationContext map[string]any

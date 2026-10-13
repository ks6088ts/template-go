# Code style

- Add Go doc comments for functions and types; explain non-obvious decisions without narrating obvious code.
- Use camelCase for variable names.
- Handle non-nil errors explicitly and wrap errors with context using `%w` when returning them.
- Use the standard `testing` package for new features and behavior changes.
- Send command results to stdout and diagnostics to stderr; propagate output write errors.
- Keep command state local to each invocation and avoid adding dependencies without a concrete need.

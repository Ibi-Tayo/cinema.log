---
name: project-context
description: Essential architectural guide for Cinema.log. Use this to ensure backend Go and frontend Angular alignment.
---
# Cinema.log Context Loader
When this skill is activated, you must strictly adhere to the patterns in `project-overview.md`:

- **Go Backend:** Follow Vertical Slice in `internal/{feature}`. Use the `DB -> Store -> Service -> Handler` pattern.
- **Angular Frontend:** Use Standalone Components, `inject()`, and Signals. No constructor injection.
- **Testing:** Use `make test` for Go and `vitest` for Angular. 
- **Scripts:** Use `./run-dev.zsh` for environment management.

*Always verify file paths against the Vertical Slice structure before proposing changes.*
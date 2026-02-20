---
name: code-implementation
description: Specialized skill for executing file changes, maintaining style guides, and writing post-hoc tests.
---
# Implementation Framework

### 1. Code Construction
- **Style Alignment:** Match the indentation, naming conventions (camelCase vs snake_case), and documentation style of the surrounding code.
- **Robustness:** Implement the error handling and edge cases identified by the Analyst.
- **Modularity:** Ensure new functions are focused and do not exceed 30-40 lines where possible.

### 2. Post-Implementation Testing
Once the code is written:
- **Coverage:** Write tests that cover the "Happy Path" and the "Edge Cases" from the spec.
- **Integration:** Ensure new components integrate with existing services without regression.
- **Validation:** Use `npm test` to verify that your new tests AND all existing tests pass.

### 3. Cleanup
- Remove any console logs or "TODO" comments used during development.
- Run the project's linter if available.
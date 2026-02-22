---
name: architect-design
description: Technical design framework for planning file structures, design patterns, and data flow.
---
# Technical Design Framework
When this skill is activated, evaluate the Analyst's spec against these architectural pillars and ensure your design aligns with the existing codebase and current patterns.

### 1. Impact Analysis
- **Core Logic:** Which existing services or controllers are affected?
- **Data Schema:** Does this require new database fields or state management changes?
- **Side Effects:** Will these changes break existing unit tests or dependencies?

### 2. Pattern Selection
- **Scalability:** If the request involves many items, suggest a **Strategy** or **Factory** pattern.
- **Decoupling:** Use **Dependency Injection** or **Events** to keep the new feature isolated.
- **Consistency:** Ensure variable naming and file locations match the existing project structure (check `/src` for naming conventions).

### 3. Verification Criteria
For every file change proposed, define a "Success Signal":
- Example: "If `AuthService.ts` is modified, the `test-login` suite must pass with 100% coverage."

### 4. Output Structure
The `design_plan.md` must include a **Mermaid.js** diagram showing the data flow between components.
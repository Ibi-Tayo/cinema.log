---
name: code-audit
description: Generic high-standard audit for Modern Angular and Backend services. 
---
# Universal Code Audit Framework

### 1. Structural Consistency
- **Pattern Alignment:** Does the code follow the directory structure and layering defined in the local project documentation?
- **Dependency Direction:** Are there circular dependencies or violations of the project's "Vertical" or "Layered" rules?

### 2. Angular Best Practices (Modern)
- **Standalone:** Are new components standalone by default?
- **DI Pattern:** Does the code use `inject()` over constructor injection?
- **Reactivity:** Is the code utilizing Signals (`signal`, `computed`, `effect`) correctly without creating race conditions?

### 3. Logic & Security
- **Auth Integrity:** Are protected routes or sensitive operations backed by appropriate auth/permission checks?
- **Data Safety:** Is there proper error handling for external API calls or DB queries?
- **Complexity:** Are functions small and single-purpose (SOLID principles)?

### 4. Test Verification
- **Gap Analysis:** Do the tests cover the edge cases and business logic defined in the feature spec?
- **Mocking:** Are external dependencies properly mocked to ensure test isolation?
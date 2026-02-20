---
name: requirement-elicitor
description: A framework for interviewing users to extract deep business context and technical constraints.
---

# Requirement Interview Checklist

When this skill is activated, evaluate the user's request against these five pillars. If any pillar is "Thin" or "Missing," you must ask a clarifying question.

### 1. Business Context (The "Why")

- What is the primary problem this feature solves?
- Who is the end user (Admin, Guest, Internal Dev)?
- What does "Success" look like? (e.g., "Page loads in <1s", "User saves 5 minutes")

### 2. Functional Scope (The "What")

- What are the required inputs and the expected outputs?
- Is there an existing UI pattern this must follow?
- Does this interact with external APIs or databases?

### 3. Technical Constraints (The "How")

- Are there specific performance requirements?
- Does this need new dependencies (npm packages)?
- Are there security/auth concerns (e.g., sensitive data handling)?

### 4. Edge Cases (The "What If")

- What happens if the network fails?
- What if the user input is empty or malformed?
- Are there race conditions (two people doing X at once)?

### 5. Definition of Done

- Does this require unit tests?
- Does the documentation need an update?

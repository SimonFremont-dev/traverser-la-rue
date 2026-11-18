# Backend stack
The project's backend is written in Go (Golang)

# Frontend stack
The project's frontend is written in JavaScript, TypeScript and React.

# TDD workflow
For every feature or bugfix:
1. RED: write one failing test first, then delegate to the TDD red subagent to confirm it fails correctly.
2. GREEN: write the minimal code to pass, then delegate to the TDD green subagent to confirm no regressions.
3. REFACTOR: clean up, then delegate to the TDD refactor subagent.
Never write implementation code while a test for it is passing or missing.

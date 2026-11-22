---
name: secrets
description: The API key, files sent to the model and the history
applies_to: ["**/*.go"]
---
asoai sends prompts, files and history to OpenAI with `OPENAI_API_KEY`, and
keeps conversations in a local database. Review the change for:

- The API key reaching output, logs, error messages or the database.
- Files read for a chat from anywhere other than what the user named, or
  sent without the user asking.
- The history database or exports created with permissions others can read.

Report only concrete, reachable problems in the changed code.

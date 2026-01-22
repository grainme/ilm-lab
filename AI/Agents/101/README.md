# AI Code Agent – Learning Notes

For learning purposes, I implemented an agentic coding assistant (similar to Cursor's agentic mode) that can read, modify, and execute code autonomously.

## What This Project Does

An LLM-powered agent that can:
- Explore a codebase (list files, read contents)
- Write and modify files
- Execute Python code
- Iterate toward solving coding tasks without constant user input

**Example**: Ask it to "fix the bug in calculator.py" and it will:
1. Read the file
2. Identify the issue
3. Write a corrected version
4. Test it by running the code

## How It Works

### 1. System Prompt
Defines the agent's role and capabilities (see `prompts.py`):
- Lists what operations are available
- Instructs the agent to be proactive: *"don't ask the user for information you can find using your tools"*

### 2. Function Calling (Tool Use)

The agent doesn't magically access your filesystem. You give it **tools** - Python functions wrapped in schemas the LLM can understand.

**See `functions/write_files.py` for example:**
- `write_file()` - the actual Python function
- `schema_write_file` - describes it to the LLM using `types.FunctionDeclaration`

The schema tells the LLM what the function does, what parameters it needs, and their types.

### 3. The Agentic Loop

**See `main.py` for the full implementation.**

The core pattern:
1. Send user request + conversation history to LLM
2. LLM responds with either:
   - **Text** (final answer) → stop
   - **Function calls** (needs to use tools) → execute and continue
3. Add function results to message history
4. Repeat until solved or max iterations reached (15 in this implementation, Tokens cost money HAHAH)

### 4. Security Concerns

**Path Traversal Protection** (see `functions/write_files.py`):

The code validates that files stay within the working directory:
- Uses `os.path.commonpath()` to ensure the target file is inside the permitted directory
- Prevents the LLM from accessing sensitive files if it hallucinates or is prompt-injected

**Working directory is auto-injected** (see `call_function.py`):
- The LLM never specifies the working directory
- It's added server-side in `call_function()` for security

## Key Learnings

- **Agentic behavior emerges from**: good system prompts + appropriate tools + iteration
- **Tool schemas must be precise**: clear descriptions help the LLM choose the right tool and parameters
- **Always validate LLM actions**: never trust file paths, commands, or code execution without security checks
- **Message history is everything**: LLMs are stateless - the conversation history contains all context
- **Using pre-trained LLMs can be expensive**: each iteration costs API credits (Gemini 2.5 Flash was used here)

## Potential Applications

The pattern learned here (system prompt + tools + agentic loop) applies to:
- Code review bots that analyze PRs
- Automated debugging assistants
- Documentation generators
- Test-writing agents


> This is not prod ready obviously.

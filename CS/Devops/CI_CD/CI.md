# CI – Learning Notes

## Questions & Answers

### Is "tests" a special keyword in job names? 

```yml
jobs:
  tests:
    name: Tests
    runs-on: ubuntu-latest

```

**Answer**: No. You can name your job anything: `tests`, `build`, `lint`, `my-cool-job` - all valid.
**Why this matters**: Choose names that clearly describe what the job does.

---

### Do multiple jobs run in parallel by default?
**Answer**: Yes. You can control this with `needs:` to make jobs run sequentially.

**Example**:
```yaml
jobs:
  build:
    runs-on: ubuntu-latest
  
  deploy:
    needs: build  # Waits for build to finish
    runs-on: ubuntu-latest
```

---

### How does GitHub handle millions of CI runs?
**Answer**: GitHub has massive server infrastructure (owned by Microsoft). They use:
- Cloud providers (Azure primarily)
- Job queuing systems
- Resource limits per account (free tier gets fewer concurrent runners than paid)
- Ephemeral runners (VMs created on-demand, destroyed after jobs finish)

**Related**: This is why there are usage limits on free accounts - compute isn't free.

---

### What does "check out code" mean?
**Answer**: It clones your repository code into the runner's workspace so subsequent steps can access your files.

**Without it**: The runner starts empty - it doesn't automatically have your code.

---

### What is "actions/" and where does it come from?
**Answer**: `actions/checkout@v4` refers to a reusable action from GitHub's official actions repository.

**Format**: `owner/repo@version`

**Examples**:
- `actions/checkout` → https://github.com/actions/checkout
- `actions/setup-go` → https://github.com/actions/setup-go

**Related**: You can use actions from any public repo or create your own.

---

### What do actions actually do under the hood?
**Answer**: Actions are just code (usually JavaScript or Docker containers) that run commands for you.

The `actions/checkout@v4` action basically does:
```bash
git clone --depth=1 <your-repo-url>
cd <your-repo>
git checkout <your-branch>
```
(Plus extra logic for tokens, submodules, LFS, etc.)

**Why use actions instead of raw commands?**

You *could* write:
```yaml
- name: Check out code
  run: |
    git clone https://github.com/myuser/myrepo.git
    cd myrepo
```

But the action handles:
- Authentication (using GitHub's automatic tokens)
- Edge cases (detached HEAD states, PR merges, etc.)
- Configuration options (shallow vs full clone, submodules, etc.)
- Cross-platform compatibility

**Actions are just repos**: Look at the checkout action source code - it's TypeScript that shells out to Git commands!

You can make your own:
```yaml
- name: My custom action
  uses: myusername/my-action@v1
```

Or just use `run:` for simple commands:
```yaml
- name: Print Go version
  run: go version
```

**Key insight**: Actions are just reusable, parameterized scripts. Nothing magical!

-
## Mental Model / How It Works

GitHub Actions is like having a fleet of disposable computers:
1. You push code → triggers workflow
2. GitHub spins up fresh VM(s) based on `runs-on`
3. VM executes your steps sequentially
4. VM is destroyed, logs are saved

Each job gets its own VM by default (isolation). Jobs can share data via artifacts if needed.

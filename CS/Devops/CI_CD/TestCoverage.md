# Code Coverage – Learning Notes

---

## Key Questions & Answers

### What is code coverage?
**Answer**: Percentage of code lines executed during tests.

```
code_coverage = (lines_covered / total_lines) * 100
```

Example: 1000 lines total, 500 lines tested → 50% coverage.

---

### Does 100% coverage mean bug-free code?
**Answer**: No. You can have 100% coverage with bugs, or 0% coverage with no bugs.

**Why**: Coverage measures *what* runs, not *how well* it's tested.

---

### Why is it controversial?
**Answer**: Hard to say "X% is good" universally. Some code is more critical to test than others.

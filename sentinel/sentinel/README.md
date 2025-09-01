# Sentinel Core

Sentinel is the self-healing security layer that detects, reflects, rewrites, or cryptographically contains adversarial prompts.

## Components

- **Detector**: Identifies security violations using embeddings, signatures, and rules
- **Reflector**: Applies constitutional AI principles to assess alignment
- **Rewriter**: Generates safe prompt alternatives with ranking
- **Router**: Decides on actions (allow, reframe, encrypt, block)
- **ToolGuard**: Manages tool/function call permissions
- **Cutter**: Implements mid-stream violation response

## Security Pipeline

```
[Input Prompt]
     ↓
[Detector] → Embeddings + Signatures + Rules = Violation Score
     ↓
[Reflector] → Constitutional prompt to assess alignment
     ↓
[Rewriter] → Generate safe alternatives + ranking
     ↓
[Router] → Decision based on policies and scores
     ↓
[Action] → Allow / Reframe / Encrypt / Block
```

## Violation Detection

The detector uses multiple approaches to identify security violations:

1. **Semantic Embeddings**: Compare prompts against known adversarial patterns
2. **Signature Matching**: Match against a database of known attack signatures
3. **Rule-Based Detection**: Apply heuristic rules for common attack patterns

## Response Actions

1. **Allow**: Permit the prompt to proceed to the LLM provider
2. **Reframe**: Request user confirmation of a rewritten safe prompt
3. **Encrypt**: Cryptographically contain the output if generated
4. **Block**: Prevent the prompt from reaching the LLM provider

## ToolGuard

ToolGuard manages permissions for function calls and tool usage:

- Disable tools when violations are detected
- Apply granular permissions based on policies
- Monitor tool inputs and outputs for sensitive data

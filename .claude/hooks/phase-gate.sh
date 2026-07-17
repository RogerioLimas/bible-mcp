#!/bin/bash
# phase-gate.sh — bloqueia skills fora da fase correta.
# Recebe JSON do evento PreToolUse via stdin.

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // empty')
CURRENT_PHASE=$(cat .claude/current-phase 2>/dev/null || echo "discovery")

# ATENCAO: o nome exato da tool/matcher para invocacao de skill do
# Superpowers (Skill, ou mcp__superpowers__<nome>) precisa ser confirmado
# rodando /hooks no Claude Code e observando o transcript real ao
# acionar a skill uma vez. Ajustar o case abaixo com o nome confirmado.

case "$TOOL_NAME" in
  *brainstorming*)
    echo "Bloqueado: skill 'brainstorming' do Superpowers conflita com SPDD. Use /grill-with-docs ou /spdd-analysis." >&2
    exit 2
    ;;
  *test-driven-development*|*subagent-driven-development*|*executing-plans*)
    if [ "$CURRENT_PHASE" != "execution" ]; then
      echo "Bloqueado: fase atual é '$CURRENT_PHASE'. Rode 'echo \"execution\" > .claude/current-phase' antes de executar TDD." >&2
      exit 2
    fi
    ;;
esac

exit 0

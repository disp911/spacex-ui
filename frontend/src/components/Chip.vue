<template>
  <span class="chip"><slot /></span>
</template>

<style scoped>
.chip {
  height: 24px;
  display: inline-flex;
  align-items: center;
  font: 500 13px var(--font-mono);
  color: var(--text);
  background: var(--chip);
  border: 1px solid var(--line-card);
  border-radius: 8px;
  padding: 0 11px;
  white-space: nowrap;
}
</style>

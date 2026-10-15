<script setup lang="ts">
/**
 * A form of its own inside a settings section, e.g. the admin's credentials:
 * the fields (the slot) up to 760px wide, then under a dashed line the note
 * (messages and a hint) and the actions that apply the form.
 */
</script>

<template>
  <div class="x-setting-block">
    <slot />
    <div v-if="$slots.note || $slots.actions" class="x-setting-block__foot">
      <div class="x-setting-block__note"><slot name="note" /></div>
      <div v-if="$slots.actions" class="x-setting-block__actions"><slot name="actions" /></div>
    </div>
  </div>
</template>

<style scoped>
.x-setting-block {
  padding: 16px 20px 18px;
  border-top: 1px solid var(--div);
}

.x-setting-block > * {
  max-width: 760px;
}

.x-setting-block__foot {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px 16px;
  margin-top: 16px;
  padding-top: 14px;
  border-top: 1px dashed var(--div);
}

.x-setting-block__note {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 6px;
  min-width: 220px;
  color: var(--dim);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
  text-wrap: pretty;
}

.x-setting-block__actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 767px) {
  .x-setting-block {
    padding: 14px 16px 16px;
  }

  /* The actions take the row under the note. */
  .x-setting-block__actions {
    flex: 1 1 100%;
  }

  .x-setting-block__actions > :deep(*) {
    flex: 1;
  }
}
</style>

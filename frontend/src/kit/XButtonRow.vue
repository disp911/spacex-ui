<script setup lang="ts">
/** Small buttons under the fields they act on, side by side, e.g. «Get New Cert» and «Clear». */
</script>

<template>
  <div class="x-btn-row"><slot /></div>
</template>

<style scoped>
.x-btn-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
</style>

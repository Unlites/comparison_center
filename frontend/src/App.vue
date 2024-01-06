<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div>
    <header>
      <nav>
        <RouterLink active-class="text-black text-decoration-none" to="/">Comparisons</RouterLink>
      </nav>
    </header>

    <div class="mt-4 mx-auto col-12 col-md-8">
      <RouterView />
    </div>
  </div>
</template>

<style scoped>
header {
  height: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-bottom: #d1dfd8 1px solid;
}
</style>

<template>
  <div>
    <Sidebar />
    <div class="c-wrapper">
      <Header />
      <div class="c-body">
        <main class="c-main">
          <nuxt />
        </main>
      </div>
      <footer class="c-footer"></footer>
    </div>
  </div>
</template>

<script>
import Header from '@/components/header.vue'
import Sidebar from '@/components/sidebar.vue'

export default {
  components: {
    Header,
    Sidebar
  }
}
</script>

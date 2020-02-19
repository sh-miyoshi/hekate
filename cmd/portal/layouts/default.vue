<template>
  <div>
    <Header />
    <Sidebar />
    <section id="main-content">
      <nuxt />
    </section>
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

<style>
html {
  font-family: 'Source Sans Pro', -apple-system, BlinkMacSystemFont, 'Segoe UI',
    Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  word-spacing: 1px;
  -ms-text-size-adjust: 100%;
  -webkit-text-size-adjust: 100%;
  -moz-osx-font-smoothing: grayscale;
  -webkit-font-smoothing: antialiased;
  box-sizing: border-box;
}
</style>

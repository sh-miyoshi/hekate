<template>
  <div class="wrapper content">
    <h1>Projects</h1>
    <ul>
      <li v-for="(project, i) in projects" :key="i">
        <nuxt-link
          to="/home"
          class="name"
          @click.native="setCurrentProject(project)"
          >{{ project }}</nuxt-link
        >
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  // middleware: 'auth',
  data() {
    return {
      projects: []
    }
  },
  created() {
    // TODO(set projects)
  },
  methods: {
    setCurrentProject(project) {
      this.$store.commit('setCurrentProject', project)
    }
  }
}
</script>

<style scoped>
.name {
  list-style: space-counter;
}
</style>

<template>
  <div class="wrapper content">
    <h1>Projects</h1>
    <ul>
      <li v-for="(project, i) in projects" :key="i" class="name">
        <nuxt-link to="/home" @click.native="setCurrentProject(project)">{{
          project
        }}</nuxt-link>
      </li>
    </ul>
    <button class="btn btn-theme" @click="$router.push('/project/new')">
      Add New Project
    </button>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      projects: []
    }
  },
  mounted() {
    this.setProjects()
  },
  methods: {
    setCurrentProject(project) {
      this.$store.commit('setCurrentProject', project)
    },
    async setProjects() {
      const res = await this.$api.ProjectGetList()
      if (res.ok) {
        this.projects = []
        for (const prj of res.data) {
          this.projects.push(prj.name)
        }
      } else {
        console.log('Failed to get project list: %o', res)
      }
    }
  }
}
</script>

<style scoped>
.name {
  list-style: disc;
}
</style>

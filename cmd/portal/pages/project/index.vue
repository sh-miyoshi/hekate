<template>
  <div class="wrapper content">
    <h1>Projects</h1>
    <ul>
      <li v-for="(project, i) in projects" :key="i">
        {{ project.name }}
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  async asyncData({ $axios }) {
    const url = process.env.SERVER_ADDR + '/api/v1/project'
    const res = await $axios.$get(url)
    return {
      projects: res
    }
  }
}
</script>

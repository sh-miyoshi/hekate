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
  data() {
    return {
      projects: []
    }
  },
  created() {
    // TODO(set projects)
  }
}
</script>

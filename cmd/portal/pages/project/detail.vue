<template>
  <div class="wrapper content">
    <h3>
      <span class="project">
        {{ this.$store.state.current_project }}
      </span>
      <span class="trush">
        <i class="fa fa-trash" @click="trash"></i>
      </span>
    </h3>

    Setting ...
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  methods: {
    trash() {
      alert('are you sure want to delete project ?')
    }
  }
}
</script>

<style scoped>
.trush:hover {
  cursor: pointer;
}

.project {
  padding-right: 20px;
}
</style>

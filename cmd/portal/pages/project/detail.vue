<template>
  <div class="wrapper content">
    <h3>
      <span class="project">
        {{ this.$store.state.current_project }}
      </span>
      <span class="trush">
        <i class="fa fa-trash" @click="trushConfirm"></i>
      </span>
    </h3>

    <div>
      <b-modal
        id="confirm-delete-project"
        ref="confirm-delete-project"
        title="Confirm"
        cancel-variant="outline-dark"
        ok-variant="danger"
        ok-title="Delete project"
        @ok="trush"
      >
        <p class="mb-0">Are you sure to delete the project ?</p>
      </b-modal>
    </div>

    Setting ...
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  methods: {
    trushConfirm() {
      this.$refs['confirm-delete-project'].show()
    },
    async trush() {
      const res = await this.$api.ProjectDelete(
        this.$store.state.current_project
      )
      console.log('project delete result: %o', res)
      if (!res.ok) {
        this.error = res.message
        return
      }

      alert('successfully deleted.')
      this.$store.commit('setCurrentProject', 'master') // TODO(set correct project name)
      this.$router.push('/home')
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

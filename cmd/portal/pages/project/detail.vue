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

    <div class="form-panel">
      <label for="accessTokenLifeSpan" class="col-sm-5 control-label elem">
        Access Token Life Span [{{ accessTokenUnit }}]
      </label>
      <div class="col-sm-5 elem">
        <input
          v-model.number="accessTokenLifeSpan"
          type="number"
          class="form-control"
        />
      </div>

      <label for="refreshTokenLifeSpan" class="col-sm-5 control-label elem">
        Refresh Token Life Span [{{ refreshTokenUnit }}]
      </label>
      <div class="col-sm-5 elem">
        <input
          v-model.number="refreshTokenLifeSpan"
          type="number"
          class="form-control"
        />
      </div>

      <div class="divider"></div>

      <div v-if="error" class="alert alert-danger">
        {{ error }}
      </div>

      <button class="btn btn-theme">Update</button>
      <nuxt-link to="/project">Cancel</nuxt-link>
    </div>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      error: '',
      accessTokenLifeSpan: 0,
      accessTokenUnit: 'sec',
      refreshTokenLifeSpan: 0,
      refreshTokenUnit: 'sec',
      signingAlgorithm: ''
    }
  },
  mounted() {
    this.getProject()
  },
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
    },
    async getProject() {
      const res = await this.$api.ProjectGet(this.$store.state.current_project)
      if (!res.ok) {
        this.error = res.message
        return
      }

      this.accessTokenLifeSpan = res.data.tokenConfig.accessTokenLifeSpan
      this.refreshTokenLifeSpan = res.data.tokenConfig.refreshTokenLifeSpan
      this.signingAlgorithm = res.data.tokenConfig.signingAlgorithm
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

.elem {
  float: left;
}

.divider {
  clear: both;
  border-bottom: 1px solid #eff2f7;
  padding-bottom: 15px;
  margin-bottom: 15px;
}
</style>

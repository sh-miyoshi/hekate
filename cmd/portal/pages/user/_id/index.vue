<template>
  <div class="wrapper content">
    <h3>
      <span v-if="user">
        {{ user.name }}
      </span>
      <span v-if="allowEdit()" class="trush">
        <i class="fa fa-trash" @click="trushConfirm"></i>
      </span>
    </h3>

    <div>
      <b-modal
        id="confirm-delete-user"
        ref="confirm-delete-user"
        title="Confirm"
        cancel-variant="outline-dark"
        ok-variant="danger"
        ok-title="Delete user"
        @ok="trush"
      >
        <p class="mb-0">Are you sure to delete the user ?</p>
      </b-modal>
    </div>

    <div class="form-panel">
      <div>
        <label for="id" class="col-sm-4 control-label elem">
          ID
        </label>
        <div class="col-sm-3 elem">
          <input
            v-if="user"
            v-model="user.id"
            class="form-control"
            disabled="disabled"
          />
        </div>
      </div>

      <div class="divider"></div>

      <div v-if="error" class="alert alert-danger">
        {{ error }}
      </div>

      <button class="btn btn-primary" @click="update">Update</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      user: null,
      error: ''
    }
  },
  mounted() {
    this.setUser(this.$route.params.id)
  },
  methods: {
    async setUser(userID) {
      const res = await this.$api.UserGet(
        this.$store.state.current_project,
        userID
      )
      if (res.ok) {
        this.user = res.data
      } else {
        console.log('Failed to get user: %o', res)
      }
    },
    allowEdit() {
      if (!this.user) {
        return false
      }
      const loginUser = window.localStorage.getItem('user')
      return this.user.name !== loginUser
    },
    trushConfirm() {
      this.$refs['confirm-delete-user'].show()
    },
    trush() {
      // TODO(implement this)
    }
  }
}
</script>

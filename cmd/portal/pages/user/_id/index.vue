<template>
  <div class="card">
    <div class="card-header">
      <h3>
        {{ currentUserName }}
        <span v-if="allowEdit()" class="trush">
          <i class="fa fa-trash" @click="trushConfirm"></i>
        </span>
      </h3>
    </div>

    <div class="card-body">
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

      <div class="form-group row">
        <label for="id" class="col-sm-2 control-label">
          ID
        </label>
        <div class="col-sm-5">
          <input
            v-if="user"
            v-model="user.id"
            class="form-control"
            disabled="disabled"
          />
        </div>
      </div>

      <div class="form-group row">
        <label for="name" class="col-sm-2 control-label">
          Name
        </label>
        <div class="col-sm-5">
          <input v-if="user" v-model="user.name" class="form-control" />
        </div>
      </div>

      <div class="form-group row">
        <label for="system-roles" class="col-sm-2 control-label">
          System Role
        </label>
        <div class="col-sm-5">
          <List v-if="user" :current="user.system_roles" :all="systemRoles" />
        </div>
      </div>

      <div class="card-footer">
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <button class="btn btn-primary" @click="update">Update</button>
      </div>
    </div>
  </div>
</template>

<script>
import List from '@/components/list.vue'

export default {
  components: {
    List
  },
  data() {
    return {
      currentUserName: '',
      user: null,
      error: '',
      systemRoles: []
    }
  },
  mounted() {
    this.systemRoles = process.env.SYSTEM_ROLES
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
        this.currentUserName = res.data.name
      } else {
        console.log('Failed to get user: %o', res)
      }
    },
    allowEdit() {
      const loginUser = window.localStorage.getItem('user')
      return this.currentUserName !== loginUser
    },
    trushConfirm() {
      this.$refs['confirm-delete-user'].show()
    },
    trush() {
      // TODO(implement this)
    },
    update() {
      // TODO(implement this)
    }
  }
}
</script>

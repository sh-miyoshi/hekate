<template>
  <div class="card">
    <div class="card-header">
      <h3>
        {{ currentUserName }}
        <span v-if="allowEdit()" class="icon">
          <i class="fa fa-trash" @click="deleteUserConfirm"></i>
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
          @ok="deleteUser"
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
          <div v-if="user">
            <div v-for="item in user.system_roles" :key="item">
              <div class="mb-1 input-group-text role-item">
                <span class="icon">
                  <i class="fa fa-remove" @click="removeSystemRole(item)"></i>
                </span>
                {{ item }}
              </div>
            </div>
            <div>
              <b-form-select
                v-model="assignedSystemRole"
                :options="assignedSystemRoleCandidates"
              ></b-form-select>
            </div>
          </div>
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
export default {
  data() {
    return {
      currentUserName: '',
      user: null,
      error: '',
      assignedSystemRole: null,
      assignedSystemRoleCandidates: this.getSystemRoleCandidates()
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
        this.currentUserName = res.data.name
      } else {
        console.log('Failed to get user: %o', res)
      }
    },
    allowEdit() {
      const loginUser = window.localStorage.getItem('user')
      return this.currentUserName !== loginUser
    },
    deleteUserConfirm() {
      this.$refs['confirm-delete-user'].show()
    },
    deleteUser() {
      // TODO(implement this)
    },
    update() {
      // TODO(implement this)
      console.log(this.user.system_roles)
    },
    getSystemRoleCandidates() {
      const res = [{ value: null, text: 'Please select an assigned role' }]
      if (!this.user) {
        return res
      }
      for (const item of process.env.SYSTEM_ROLES) {
        if (!this.user.system_roles.includes(item)) {
          res.push({ value: item, text: item })
        }
      }
      return res
    },
    removeSystemRole(role) {
      if (!this.user) {
        return
      }
      const tmp = this.user.system_roles.filter((v) => v !== role)
      this.user.system_roles = tmp
      this.assignedSystemRoleCandidates = this.getSystemRoleCandidates()
    }
  }
}
</script>

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
        <div v-if="user" class="col-sm-5">
          <div v-if="user.system_roles.length > 0">
            <div v-for="item in user.system_roles" :key="item">
              <div class="mb-1 input-group-text role-item">
                <span class="icon">
                  <i class="fa fa-remove" @click="removeSystemRole(item)"></i>
                </span>
                {{ item }}
              </div>
            </div>
          </div>
          <div v-else>
            No assigned roles
          </div>
        </div>
        <div class="col-sm-5">
          <b-form-select
            v-model="assignedSystemRole"
            :options="getSystemRoleCandidates()"
          ></b-form-select>
          <button class="btn btn-primary mt-2" @click="assignSystemRole()">
            Assign
          </button>
        </div>
      </div>

      <div class="form-group row">
        <label for="custom-roles" class="col-sm-2 control-label">
          Custom Role
        </label>
        <div v-if="user" class="col-sm-5">
          <div v-if="user.custom_roles.length > 0">
            <div v-for="item in user.custom_roles" :key="item">
              <div class="mb-1 input-group-text role-item">
                <span class="icon">
                  <i class="fa fa-remove" @click="removeCustomRole(item)"></i>
                </span>
                {{ item }}
              </div>
            </div>
          </div>
          <div v-else>
            No assigned roles
          </div>
        </div>
        <div class="col-sm-5">
          <b-form-select
            v-model="assignedCustomRole"
            :options="getCustomRoleCandidates()"
          ></b-form-select>
          <button class="btn btn-primary mt-2" @click="assignCustomRole()">
            Assign
          </button>
        </div>
      </div>

      <div class="form-group row">
        <label class="col-sm-2 control-label">
          Login Sessions
        </label>
        <div v-if="user" class="col-sm-5">
          <ul>
            <li v-for="item in user.sessions" :key="item">{{ item }}</li>
          </ul>
        </div>
      </div>

      <div class="card-footer">
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <button class="btn btn-primary" @click="updateUser">Update</button>
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
      assignedCustomRole: null
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
        if (!this.user.system_roles) {
          this.user.system_roles = []
        }
        if (!this.user.custom_roles) {
          this.user.custom_roles = []
        }
        this.currentUserName = res.data.name
        console.log(this.user)
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
    updateUser() {
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
    },
    assignSystemRole() {
      if (!this.user || !this.assignedSystemRole) {
        return
      }

      this.user.system_roles.push(this.assignedSystemRole)
      this.assignedSystemRole = null
    },
    getCustomRoleCandidates() {
      const res = [{ value: null, text: 'Please select an assigned role' }]
      // TODO(get all custom roles, check user.custom_roles)
      return res
    },
    removeCustomRole(role) {
      if (!this.user) {
        return
      }
      const tmp = this.user.custom_roles.filter((v) => v !== role)
      this.user.custom_roles = tmp
    },
    assignCustomRole() {
      if (!this.user || !this.assignedCustomRole) {
        return
      }

      this.user.custom_roles.push(this.assignedCustomRole)
      this.assignedCustomRole = null
    }
  }
}
</script>

<style scoped>
.role-item {
  display: inline-block;
  padding: 0.175rem 0.55rem;
  width: auto;
}
</style>

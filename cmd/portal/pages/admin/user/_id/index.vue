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

    <div class="nav-tabs-boxed">
      <ul class="nav nav-tabs" role="tablist">
        <li class="nav-item">
          <a
            class="nav-link"
            :class="{ active: activePanel == 'home' }"
            @click="activate('home')"
          >
            Home
          </a>
        </li>
        <li class="nav-item">
          <a
            class="nav-link"
            :class="{ active: activePanel == 'credentials' }"
            @click="activate('credentials')"
          >
            Credentials
          </a>
        </li>
      </ul>
      <div class="tab-content">
        <div class="tab-pane" :class="{ active: activePanel == 'home' }">
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
              <label for="locked" class="col-sm-2 control-label">
                Locked
              </label>
              <div v-if="user" class="col-sm-5">
                <label
                  class="c-switch c-switch-label c-switch-pill c-switch-primary"
                >
                  <input
                    v-model="userLocked"
                    class="c-switch-input"
                    type="checkbox"
                    :disabled="!userLocked"
                  />
                  <span
                    class="c-switch-slider"
                    data-checked="On"
                    data-unchecked="Off"
                  ></span>
                </label>
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
                        <i
                          class="fa fa-remove"
                          @click="removeSystemRole(item)"
                        ></i>
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
                <button
                  class="btn btn-primary mt-2"
                  @click="assignSystemRole()"
                >
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
                  <div v-for="role in user.custom_roles" :key="role.id">
                    <div class="mb-1 input-group-text role-item">
                      <span class="icon">
                        <i
                          class="fa fa-remove"
                          @click="removeCustomRole(role)"
                        ></i>
                      </span>
                      {{ role.name }}
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
                  :options="customRoleCandidates"
                ></b-form-select>
                <button
                  class="btn btn-primary mt-2"
                  @click="assignCustomRole()"
                >
                  Assign
                </button>
              </div>
            </div>

            <div class="form-group">
              <button
                class="btn btn-link dropdown-toggle ml-n3"
                @click="loadLoginSessions()"
              >
                Login Sessions
              </button>
              <div v-if="showLoginSessions" class="card-body">
                <div class="form-group row">
                  <table class="table table-responsive-sm">
                    <thead>
                      <tr>
                        <th>ID</th>
                        <th>From IP</th>
                        <th>Session Start</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="session in loginSessions" :key="session.id">
                        <td>{{ session.id }}</td>
                        <td>{{ session.from_ip }}</td>
                        <td>{{ session.created_at }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            <div class="card-footer">
              <div v-if="error" class="alert alert-danger">
                {{ error }}
              </div>

              <button class="btn btn-primary" @click="updateUser">
                Update
              </button>
            </div>
          </div>
        </div>
        <div class="tab-pane" :class="{ active: activePanel == 'credentials' }">
          <div class="card-body">
            <div class=" ml-3">
              <div class="form-group row">
                <h4>Force Reset Password</h4>
              </div>

              <div class="form-group row">
                <label for="password" class="col-sm-2 control-label">
                  Password
                  <span class="required">*</span>
                </label>
                <div class="col-sm-5">
                  <input
                    v-model="password"
                    class="form-control"
                    type="password"
                  />
                </div>
              </div>

              <div class="form-group row">
                <label for="confirm-password" class="col-sm-2 control-label">
                  Confirm Password
                  <span class="required">*</span>
                </label>
                <div class="col-sm-5">
                  <input
                    v-model="confirmPassword"
                    class="form-control"
                    type="password"
                  />
                </div>
              </div>
            </div>

            <div class="card-footer">
              <div v-if="error" class="alert alert-danger">
                {{ error }}
              </div>

              <button class="btn btn-danger" @click="resetPassword">
                Reset
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      currentUserName: '',
      activePanel: 'home',
      user: null,
      error: '',
      assignedSystemRole: null,
      assignedCustomRole: null,
      customRoleCandidates: [],
      showLoginSessions: false,
      loginSessions: [],
      SYSTEM_ROLES: [],
      userLocked: false,
      password: '',
      confirmPassword: ''
    }
  },
  async mounted() {
    const mainProject = process.env.LOGIN_PROJECT
    if (this.$store.state.current_project === mainProject) {
      this.SYSTEM_ROLES.push('read-cluster')
      this.SYSTEM_ROLES.push('write-cluster')
    }
    this.SYSTEM_ROLES.push('read-project')
    this.SYSTEM_ROLES.push('write-project')

    await this.setUser(this.$route.params.id)
    await this.setCustomRoleCandidates()
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
        this.userLocked = res.data.locked
        console.log('User Info: %o', this.user)
        return
      }
      console.log('Failed to get user: %o', res)
      this.error = res.message
    },
    allowEdit() {
      const loginUser = window.localStorage.getItem('user')
      return this.currentUserName !== loginUser
    },
    deleteUserConfirm() {
      this.$refs['confirm-delete-user'].show()
    },
    async deleteUser() {
      console.log('delete user id: ' + this.$route.params.id)
      const res = await this.$api.UserDelete(
        this.$store.state.current_project,
        this.$route.params.id
      )
      if (!res.ok) {
        this.error = res.message
        return
      }
      await this.$bvModal.msgBoxOk('Successfully delete user')
      this.$router.push('/admin/user')
    },
    async updateUser() {
      if (!this.user) {
        return
      }

      const projectName = this.$store.state.current_project
      const userID = this.$route.params.id

      if (!this.userLocked && this.user.locked) {
        // unlock user
        const res = await this.$api.UserUnlock(projectName, userID)
        if (!res.ok) {
          this.error = res.message
          return
        }
      }

      const roles = []
      for (const r of this.user.custom_roles) {
        roles.push(r.id)
      }
      const data = {
        name: this.user.name,
        system_roles: this.user.system_roles,
        custom_roles: roles
      }
      const res = await this.$api.UserUpdate(projectName, userID, data)
      if (!res.ok) {
        this.error = res.message
        return
      }
      await this.$bvModal.msgBoxOk('Successfully update user')
    },
    getSystemRoleCandidates() {
      const res = [{ value: null, text: 'Please select an assigned role' }]
      if (!this.user) {
        return res
      }
      for (const item of this.SYSTEM_ROLES) {
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
    async setCustomRoleCandidates() {
      const res = [{ value: null, text: 'Please select an assigned role' }]
      const roleRes = await this.$api.RoleGetList(
        this.$store.state.current_project
      )
      if (!roleRes.ok) {
        console.log('Custom role get failed: %o', roleRes)
        this.error = 'Failed to get custom role list.'
        return
      }

      for (const role of roleRes.data) {
        // remove user assigned role
        let append = true
        if (this.user) {
          for (const ur of this.user.custom_roles) {
            if (role.id === ur.id) {
              append = false
              break
            }
          }
        }
        if (append) {
          res.push({ value: role.id, text: role.name })
        }
      }
      this.customRoleCandidates = res
    },
    removeCustomRole(role) {
      if (!this.user) {
        return
      }
      const tmp = this.user.custom_roles.filter((v) => v.id !== role.id)
      this.user.custom_roles = tmp
    },
    assignCustomRole() {
      if (!this.user || !this.assignedCustomRole) {
        return
      }

      for (const r of this.customRoleCandidates) {
        if (r.value === this.assignedCustomRole) {
          this.user.custom_roles.push({
            id: r.value,
            name: r.text
          })
          this.setCustomRoleCandidates()
          this.assignedCustomRole = null
          return
        }
      }
    },
    async loadLoginSessions() {
      this.showLoginSessions = !this.showLoginSessions
      const projectName = this.$store.state.current_project

      // load sessions when open dropdown
      if (
        this.showLoginSessions &&
        this.user != null &&
        this.user.sessions != null
      ) {
        for (const sid of this.user.sessions) {
          const res = await this.$api.SessionGet(projectName, sid)
          if (!res.ok) {
            console.log('Failed to get session: %o', res)
            this.error = res.message
            return
          }
          this.loginSessions.push(res.data)
        }
      }
    },
    activate(panel) {
      this.activePanel = panel
    },
    async resetPassword() {
      if (this.password === '') {
        this.error = 'Please set password'
        return
      }
      if (this.password !== this.confirmPassword) {
        this.error = 'Password and confirmation are not same'
        return
      }

      const projectName = this.$store.state.current_project
      const userID = this.$route.params.id
      const res = await this.$api.UserResetPassword(
        projectName,
        userID,
        this.password
      )
      this.password = ''
      this.confirmPassword = ''

      if (!res.ok) {
        this.error = res.message
        return
      }
      await this.$bvModal.msgBoxOk('Successfully reset password')
      await this.setUser(userID)
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

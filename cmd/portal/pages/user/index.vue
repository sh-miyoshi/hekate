<template>
  <div class="card">
    <div class="card-header">
      <h3>Users</h3>
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

      <table class="table table-responsive-sm">
        <thead>
          <tr>
            <td>ID</td>
            <td>Name</td>
            <td>Actions</td>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.name }}</td>
            <td>
              <button
                class="btn btn-primary"
                @click="$router.push('/user/' + user.id)"
              >
                edit
              </button>
              <span v-if="allowEdit(user.name)" class="icon ml-2 h4">
                <i class="fa fa-trash" @click="deleteUserConfirm(user.id)"></i>
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- TODO show error -->

      <button class="btn btn-primary" @click="$router.push('/user/new')">
        Add New User
      </button>
    </div>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      users: [],
      deleteUserID: '',
      error: ''
    }
  },
  mounted() {
    this.setUsers()
  },
  methods: {
    async setUsers() {
      const res = await this.$api.UserGetList(this.$store.state.current_project)
      if (res.ok) {
        this.users = []
        for (const usr of res.data) {
          this.users.push(usr)
        }
      } else {
        console.log('Failed to get user list: %o', res)
      }
    },
    allowEdit(userName) {
      const loginUser = window.localStorage.getItem('user')
      return userName !== loginUser
    },
    deleteUserConfirm(id) {
      this.deleteUserID = id
      this.$refs['confirm-delete-user'].show()
    },
    async deleteUser() {
      console.log('delete user id: ' + this.deleteUserID)
      const res = await this.$api.UserDelete(
        this.$store.state.current_project,
        this.deleteUserID
      )
      if (!res.ok) {
        this.error = res.msg
        return
      }
      this.setUsers()
      alert('Successfully delete user')
    }
  }
}
</script>

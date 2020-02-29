<template>
  <div class="wrapper content">
    <h1>Users</h1>
    <table border="1">
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
            <button>edit</button>
            <span v-if="allowEdit(user.name)" class="trush">
              <i class="fa fa-trash" @click="trushConfirm"></i>
            </span>
          </td>
        </tr>
      </tbody>
    </table>

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

    <button class="btn btn-theme">
      Add New User
    </button>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      users: []
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
      // TODO(check userName and loggined user name)
      return true
    },
    trushConfirm() {
      this.$refs['confirm-delete-user'].show()
    },
    trush() {}
  }
}
</script>

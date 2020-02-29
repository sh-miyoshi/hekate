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
          <td></td>
        </tr>
      </tbody>
    </table>
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
    }
  }
}
</script>

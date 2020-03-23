<template>
  <div class="card">
    <div class="card-header">
      <h3>Roles</h3>
    </div>
    <div class="card-body">
      <div>
        <b-modal
          id="confirm-delete-role"
          ref="confirm-delete-role"
          title="Confirm"
          cancel-variant="outline-dark"
          ok-variant="danger"
          ok-title="Delete role"
          @ok="deleteRole"
        >
          <p class="mb-0">Are you sure to delete the role ?</p>
        </b-modal>
      </div>

      <table class="table table-responsive-sm">
        <thead>
          <tr>
            <td>Name</td>
            <td>Description</td>
            <td>Actions</td>
          </tr>
        </thead>
        <tbody>
          <tr v-for="role in roles" :key="role.id">
            <td>{{ role.name }}</td>
            <td></td>
            <td>
              <button
                class="btn btn-primary"
                @click="$router.push('/role/' + role.id)"
              >
                edit
              </button>
              <span class="icon ml-2 h4">
                <i class="fa fa-trash" @click="deleteRoleConfirm(role.id)"></i>
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- TODO show error -->

      <button class="btn btn-primary" @click="$router.push('/role/new')">
        Add New Role
      </button>
    </div>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      roles: [],
      deleteRoleID: '',
      error: ''
    }
  },
  mounted() {
    this.setRoles()
  },
  methods: {
    async setRoles() {
      const res = await this.$api.RoleGetList(this.$store.state.current_project)
      if (res.ok) {
        this.roles = []
        for (const r of res.data) {
          this.roles.push(r)
        }
      } else {
        console.log('Failed to get role list: %o', res)
      }
    },
    deleteRoleConfirm(id) {
      this.deleteRoleID = id
      this.$refs['confirm-delete-role'].show()
    },
    async deleteRole() {
      console.log('delete role id: ' + this.deleteRoleID)
      const res = await this.$api.RoleDelete(
        this.$store.state.current_project,
        this.deleteRoleID
      )
      if (!res.ok) {
        this.error = res.msg
        return
      }
      this.setRoles()
      alert('Successfully delete role')
    }
  }
}
</script>

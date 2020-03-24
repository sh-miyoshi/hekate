<template>
  <div class="card">
    <div class="card-header">
      <h3>
        {{ currentRoleName }}
        <span class="icon">
          <i class="fa fa-trash" @click="deleteRoleConfirm"></i>
        </span>
      </h3>
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

      <div class="form-group row">
        <label for="name" class="col-sm-2 control-label">
          Name
        </label>
        <div class="col-sm-5">
          <input v-if="role" v-model="role.name" class="form-control" />
        </div>
      </div>

      <div class="card-footer">
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <button class="btn btn-primary" @click="updateRole">Update</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      currentRoleName: '',
      role: null,
      error: ''
    }
  },
  mounted() {
    this.setRole(this.$route.params.id)
  },
  methods: {
    async setRole(roleID) {
      const res = await this.$api.RoleGet(
        this.$store.state.current_project,
        roleID
      )
      if (res.ok) {
        this.role = res.data
        this.currentRoleName = res.data.name
      } else {
        console.log('Failed to get role: %o', res)
      }
    },
    deleteRoleConfirm() {
      this.$refs['confirm-delete-role'].show()
    },
    async deleteRole() {
      if (!this.role) {
        return
      }
      console.log('delete role id: ' + this.role.id)
      const res = await this.$api.RoleDelete(
        this.$store.state.current_project,
        this.role.id
      )
      if (!res.ok) {
        this.error = res.msg
        return
      }
      alert('Successfully delete role')
      this.$router.push('/role')
    },
    async updateRole() {
      if (!this.role) {
        return
      }

      const data = {
        name: this.role.name
      }

      const res = await this.$api.RoleUpdate(
        this.$store.state.current_project,
        this.role.id,
        data
      )

      if (!res.ok) {
        this.error = res.msg
        return
      }
      alert('Successfully update role')
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

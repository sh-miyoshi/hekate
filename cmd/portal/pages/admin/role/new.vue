<template>
  <div class="card">
    <div class="card-header">
      <h3>New Role Info</h3>
    </div>

    <div class="card-body">
      <div class="form-group row">
        <label for="id" class="col-sm-2 col-form-label">
          Name
          <span class="required">*</span>
        </label>
        <div class="col-md-5">
          <input
            v-model="name"
            type="text"
            class="form-control"
            :class="{ 'is-invalid': nameValidateError }"
            @blur="validateRoleName()"
          />
          <div class="invalid-feedback">
            {{ nameValidateError }}
          </div>
        </div>
      </div>

      <div class="card-footer">
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <button class="btn btn-primary mr-2" @click="create">Create</button>
        <nuxt-link to="/admin/role">Cancel</nuxt-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ValidateRoleName } from '~/plugins/validation'

export default {
  data() {
    return {
      name: '',
      nameValidateError: '',
      error: ''
    }
  },
  methods: {
    async create() {
      if (this.nameValidateError.length > 0) {
        this.error = 'Please fix validation error before create.'
        return
      }

      const data = {
        name: this.name
      }
      const projectName = this.$store.state.current_project
      const res = await this.$api.RoleCreate(projectName, data)
      console.log('role create result: %o', res)
      if (!res.ok) {
        this.error = res.message
        return
      }

      await this.$bvModal.msgBoxOk('successfully created.')
      this.$router.push('/admin/role')
    },
    validateRoleName() {
      const res = ValidateRoleName(this.name)
      if (!res.ok) {
        this.nameValidateError = res.message
      } else {
        this.nameValidateError = ''
      }
    }
  }
}
</script>

<style scoped>
.required {
  color: #ee2222;
  font-size: 18px;
}
</style>

<template>
  <div class="card">
    <div class="card-header">
      <h3>New User Info</h3>
    </div>

    <div class="card-body">
      <div class="form-group row">
        <label for="name" class="col-sm-2 col-form-label">
          Name
          <span class="required">*</span>
        </label>
        <div class="col-md-5">
          <input
            v-model="name"
            type="text"
            class="form-control"
            :class="{ 'is-invalid': nameValidateError }"
            @blur="validateUserName()"
          />
          <div class="invalid-feedback">
            {{ nameValidateError }}
          </div>
        </div>
      </div>

      <div class="form-group row">
        <label for="password" class="col-sm-2 col-form-label">
          Password
          <span class="required">*</span>
        </label>
        <div class="col-md-5">
          <input
            v-model="password"
            type="password"
            class="form-control"
            :class="{ 'is-invalid': passwordValidateError }"
          />
          <div class="invalid-feedback">
            {{ passwordValidateError }}
          </div>
        </div>
      </div>

      <div class="card-footer">
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <button class="btn btn-primary mr-2" @click="create">Create</button>
        <nuxt-link to="/user">Cancel</nuxt-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ValidateUserName } from '~/plugins/validation'

export default {
  data() {
    return {
      name: '',
      nameValidateError: '',
      password: '',
      passwordValidateError: '',
      error: ''
    }
  },
  methods: {
    async create() {
      if (
        this.nameValidateError.length > 0 ||
        this.passwordValidateError.length > 0
      ) {
        this.error = 'Please fix validation error before create.'
        return
      }

      const data = {
        name: this.name,
        password: this.password
      }
      const projectName = this.$store.state.current_project
      const res = await this.$api.UserCreate(projectName, data)
      console.log('user create result: %o', res)
      if (!res.ok) {
        this.error = res.message
        return
      }

      alert('successfully created.')
      this.$router.push('/user')
    },
    validateUserName() {
      const res = ValidateUserName(this.name)
      if (!res.ok) {
        this.nameValidateError = res.msg
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

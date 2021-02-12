<template>
  <div class="card">
    <div class="card-header">
      <h3>Password</h3>
    </div>
    <div class="card-body">
      <div class="form-group row">
        <label for="id" class="col-sm-3 control-label">
          New Password
          <span class="required">*</span>
        </label>
        <div class="col-sm-7">
          <input v-model="password" class="form-control" />
        </div>
      </div>
      <div class="form-group row">
        <label for="confirm" class="col-sm-3 control-label">
          Confirm
          <span class="required">*</span>
        </label>
        <div class="col-sm-7">
          <input v-model="confirm" class="form-control" />
        </div>
      </div>
    </div>
    <div class="card-footer">
      <div v-if="error" class="alert alert-danger">
        {{ error }}
      </div>

      <button class="btn btn-primary mr-2" @click="update">Update</button>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'user',
  middleware: 'userAuth',
  data() {
    return {
      error: '',
      password: '',
      confirm: ''
    }
  },
  methods: {
    update() {
      if (this.password !== this.confirm) {
        this.error = 'password and confirm are not same.'
        return
      }
      this.error = ''
      // TODO change password
    }
  }
}
</script>

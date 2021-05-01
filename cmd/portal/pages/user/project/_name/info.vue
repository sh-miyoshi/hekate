<template>
  <div class="card">
    <div class="card-header">
      <h3>User Info</h3>
    </div>
    <div class="card-body">
      <div class="form-group row">
        <label for="id" class="col-sm-2 control-label"> ID </label>
        <div class="col-sm-7">
          <input v-model="id" class="form-control" disabled />
        </div>
      </div>
      <div class="form-group row">
        <label for="name" class="col-sm-2 control-label"> Name </label>
        <div class="col-sm-7">
          <input v-model="name" class="form-control" disabled />
        </div>
      </div>
      <div class="form-group row">
        <label for="email" class="col-sm-2 control-label"> E-mail </label>
        <div class="col-sm-7">
          <input v-model="email" class="form-control" disabled />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import jwtdecode from 'jwt-decode'

export default {
  layout: 'user',
  middleware: 'userAuth',
  data() {
    return {
      id: '',
      name: '',
      email: ''
    }
  },
  mounted() {
    this.setUserInfo()
  },
  methods: {
    setUserInfo() {
      const token = window.localStorage.getItem('access_token')
      const data = jwtdecode(token)
      this.id = data.sub
      this.name = data.preferred_username
      this.email = data.email
    }
  }
}
</script>

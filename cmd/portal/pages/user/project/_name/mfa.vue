<template>
  <div class="card">
    <div class="card-header">
      <h3>Authenticator Application</h3>
    </div>
    <div v-if="isOTPEnabled()" class="card-body">
      <div class="form-group row">
        TODO: show OTP setting
      </div>
    </div>
    <div v-else class="card-body">
      <div v-if="qrcode">
        <div class="form-group row">
          <ol>
            <li>
              Install one of the following applications on your mobile:
              <ul class="otp-apps">
                <li>FreeOTP</li>
                <li>Google Authenticator</li>
              </ul>
            </li>
            <li>Open the application and scan the QR code:</li>
            <img :src="qrcode" />
            <li>
              Enter the one-time code provided by the application and click
              Submit button to finish the setup.
            </li>
          </ol>
        </div>
        <div class="form-group row">
          <label for="id" class="col-sm-2 control-label">
            One-time Code <span class="required">*</span>
          </label>
          <div class="col-sm-5">
            <input v-model="usercode" class="form-control" />
          </div>
        </div>
        <div class="form-group row">
          <button class="btn btn-primary ml-3" @click="submit">Submit</button>
          <button class="btn btn-link ml-3" @click="cancel">Cancel</button>
        </div>
      </div>
      <div v-else>
        <div class="form-group row col-md-8">
          Authenticator Application is not set up.
        </div>
        <div class="form-group row col-md-8">
          <button class="btn btn-link" @click="generateQRCode()">
            Set Application
          </button>
        </div>
      </div>
    </div>
    <div v-if="error" class="card-footer">
      <div class="alert alert-danger">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'user',
  middleware: 'userAuth',
  data() {
    return {
      user: null,
      qrcode: '',
      error: '',
      usercode: ''
    }
  },
  mounted() {
    this.setUser()
  },
  methods: {
    isOTPEnabled() {
      if (!this.user) {
        return false
      }
      return this.user.otp_info.enabled
    },
    async setUser() {
      const project = window.localStorage.getItem('login_project')
      const userID = window.localStorage.getItem('user_id')
      const res = await this.$api.UserAPIGetUser(project, userID)
      console.log('get user result: %o', res)
      if (!res.ok) {
        this.error = res.message
        return
      }
      this.user = res.data
    },
    async generateQRCode() {
      const project = window.localStorage.getItem('login_project')
      const userID = window.localStorage.getItem('user_id')
      const res = await this.$api.UserAPIOTPGenerate(project, userID)
      console.log('generate QR code result: %o', res)
      if (!res.ok) {
        this.error = res.message
        return
      }
      this.qrcode = 'data:image/png;base64,' + res.data.qrcode
    },
    cancel() {
      this.qrcode = ''
      this.error = ''
    },
    async submit() {
      let res = this.validateUserCode(this.usercode)
      if (res !== '') {
        this.error = res
        return
      }

      const project = window.localStorage.getItem('login_project')
      const userID = window.localStorage.getItem('user_id')
      res = await this.$api.UserAPIOTPVerify(project, userID, this.usercode)
      if (!res.ok) {
        this.error = res.message
        return
      }
      this.error = ''
      this.setUser()
    },
    validateUserCode(code) {
      if (typeof code !== 'string') {
        return 'The code is not string'
      }

      if (code.length !== 6) {
        return 'The length of code is not 6'
      }

      // all char is only digit
      for (let i = 0; i < code.length; i++) {
        if (code[i] < '0' || code[i] > '9') {
          return 'The code contains non-numeric characters'
        }
      }

      return ''
    }
  }
}
</script>

<style scoped>
.otp-apps {
  padding-left: 30px;
}
</style>

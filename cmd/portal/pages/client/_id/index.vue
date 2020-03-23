<template>
  <div class="card">
    <div class="card-header">
      <h3>
        {{ currentClientID }}
        <span class="icon">
          <i class="fa fa-trash" @click="deleteClientConfirm"></i>
        </span>
      </h3>
    </div>

    <div class="card-body">
      <div>
        <b-modal
          id="confirm-delete-client"
          ref="confirm-delete-client"
          title="Confirm"
          cancel-variant="outline-dark"
          ok-variant="danger"
          ok-title="Delete client"
          @ok="deleteClient"
        >
          <p class="mb-0">Are you sure to delete the client ?</p>
        </b-modal>
      </div>

      <div class="form-group row">
        <label for="id" class="col-sm-2 control-label">
          ID
        </label>
        <div class="col-sm-5">
          <input v-if="client" v-model="client.id" class="form-control" />
        </div>
      </div>

      <div class="form-group row">
        <label for="accessType" class="col-sm-2 col-form-label">
          Access Type
        </label>
        <div class="col-md-5">
          <select
            v-if="client"
            v-model="client.access_type"
            name="accessType"
            class="form-control"
          >
            <option>confidential</option>
            <option>public</option>
          </select>
        </div>
      </div>

      <div
        v-if="client && client.access_type === 'confidential'"
        class="form-group row"
      >
        <label for="secret" class="col-sm-2 col-form-label">
          Secret
        </label>
        <div class="col-md-5">
          <input
            v-if="client"
            v-model="client.secret"
            type="text"
            class="form-control"
            disabled="disabled"
          />
        </div>
        <div class="col-md-5">
          <button class="btn btn-dark mr-2" @click="generateSecret">
            Regenerate Secret
          </button>
        </div>
      </div>

      <!-- TODO AllowedCallbackURLs -->
      <div class="form-group row">
        <label for="callbacks" class="col-sm-2 col-form-label">
          Allowed Callback URL
        </label>
        <div v-if="client" class="col-md-5">
          <div
            v-for="(url, i) in client.allowed_callback_urls"
            :key="i"
            class="input-group mb-1"
          >
            <input
              v-model="client.allowed_callback_urls[i]"
              class="form-control"
              type="url"
            />
            <div class="input-group-append">
              <span class="input-group-text icon" @click="removeCallback(i)">
                <i class="fa fa-trash"></i>
              </span>
            </div>
          </div>
          <div class="input-group mb-1">
            <input v-model="newCallback" class="form-control" type="url" />
            <div class="input-group-append" @click="appendCallback">
              <span class="input-group-text icon">
                +
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="card-footer">
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <button class="btn btn-primary" @click="updateClient">Update</button>
      </div>
    </div>
  </div>
</template>

<script>
import { v4 as uuidv4 } from 'uuid'
import validator from 'validator'

export default {
  data() {
    return {
      currentClientID: '',
      client: null,
      error: '',
      newCallback: ''
    }
  },
  mounted() {
    this.setClient(this.$route.params.id)
  },
  methods: {
    async setClient(clientID) {
      const res = await this.$api.ClientGet(
        this.$store.state.current_project,
        clientID
      )
      if (res.ok) {
        this.client = res.data
        this.currentClientID = res.data.id
        console.log(this.client)
      } else {
        console.log('Failed to get client: %o', res)
      }
    },
    deleteClientConfirm() {
      this.$refs['confirm-delete-client'].show()
    },
    async deleteClient() {
      if (!this.client) {
        return
      }
      console.log('delete client id: ' + this.client.id)
      const res = await this.$api.ClientDelete(
        this.$store.state.current_project,
        this.client.id
      )
      if (!res.ok) {
        this.error = res.msg
        return
      }
      alert('Successfully delete client')
      this.$router.push('/client')
    },
    async updateClient() {
      if (!this.client) {
        return
      }

      const data = {
        secret: this.client.secret,
        access_type: this.client.access_type,
        allowed_callback_urls: this.client.allowed_callback_urls
      }

      const res = await this.$api.ClientUpdate(
        this.$store.state.current_project,
        this.client.id,
        data
      )

      if (!res.ok) {
        this.error = res.msg
        return
      }
      alert('Successfully update client')
    },
    generateSecret() {
      if (!this.client) {
        return
      }
      this.client.secret = uuidv4()
    },
    appendCallback() {
      if (!this.client) {
        return
      }

      if (!validator.isURL(this.newCallback, { require_tld: false })) {
        this.error = 'New callback url is invalid url format.'
        return
      }

      if (!this.client.allowed_callback_urls) {
        this.client.allowed_callback_urls = [this.newCallback]
        this.newCallback = ''
        return
      }

      if (this.client.allowed_callback_urls.includes(this.newCallback)) {
        this.error = 'The url ' + this.newCallback + ' was already appended'
        this.newCallback = ''
        return
      }

      this.client.allowed_callback_urls.push(this.newCallback)
      this.newCallback = ''
    },
    removeCallback(index) {
      if (!this.client) {
      }

      this.client.allowed_callback_urls.splice(index, 1)
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

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

export default {
  data() {
    return {
      currentClientID: '',
      client: null,
      error: ''
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
    deleteClient() {
      // TODO(implement this)
    },
    updateClient() {
      // TODO(implement this)
    },
    generateSecret() {
      if (!this.client) {
        return
      }
      this.client.secret = uuidv4()
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

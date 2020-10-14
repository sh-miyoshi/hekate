<template>
  <div class="card">
    <div class="card-header">
      <h3>Clients</h3>
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

      <table class="table table-responsive-sm">
        <thead>
          <tr>
            <td>Client ID</td>
            <td>Access Type</td>
            <td>Actions</td>
          </tr>
        </thead>
        <tbody>
          <tr v-for="client in clients" :key="client.id">
            <td>{{ client.id }}</td>
            <td>{{ client.access_type }}</td>
            <td>
              <button
                class="btn btn-primary"
                @click="$router.push('/admin/client/' + client.id)"
              >
                edit
              </button>
              <span v-if="client.id !== mainClientID" class="icon ml-2 h4">
                <i
                  class="fa fa-trash"
                  @click="deleteClientConfirm(client.id)"
                ></i>
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="error" class="alert alert-danger">
        {{ error }}
      </div>

      <button
        class="btn btn-primary"
        @click="$router.push('/admin/client/new')"
      >
        Add New Client
      </button>
    </div>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      mainClientID: process.env.CLIENT_ID,
      clients: [],
      deleteClientID: '',
      error: ''
    }
  },
  mounted() {
    this.setClients()
  },
  methods: {
    async setClients() {
      const res = await this.$api.ClientGetList(
        this.$store.state.current_project
      )
      if (res.ok) {
        this.clients = []
        for (const cli of res.data) {
          this.clients.push(cli)
        }
      } else {
        console.log('Failed to get client list: %o', res)
      }
    },
    deleteClientConfirm(id) {
      this.deleteClientID = id
      this.$refs['confirm-delete-client'].show()
    },
    async deleteClient() {
      console.log('delete client id: ' + this.deleteClientID)
      const res = await this.$api.ClientDelete(
        this.$store.state.current_project,
        this.deleteClientID
      )
      if (!res.ok) {
        this.error = res.message
        return
      }
      this.setClients()
      await this.$bvModal.msgBoxOk('Successfully delete client')
    }
  }
}
</script>

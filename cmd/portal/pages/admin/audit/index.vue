<template>
  <div class="card">
    <div class="card-header">
      <h3>Audit Events</h3>
    </div>
    <div class="card-body">
      <div>
        <b-modal
          id="view-error-message"
          ref="view-error-message"
          title="Error Message"
          hide-footer
        >
          <p class="text-break">
            {{ auditErrorMessage }}
          </p>
          <button class="btn btn-primary" @click="hideErrorMessage()">
            OK
          </button>
        </b-modal>
      </div>

      <table class="table table-responsive-sm">
        <thead>
          <tr>
            <td>Time</td>
            <td>Method</td>
            <td>Resource Type</td>
            <td>Path</td>
            <td>Details</td>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(audit, i) in audits" :key="i">
            <td>{{ audit.time }}</td>
            <td>{{ audit.method }}</td>
            <td>{{ audit.resource_type }}</td>
            <td>
              <span
                :id="'resource-path-' + i"
                class="d-inline-block text-truncate"
                style="max-width: 200px"
              >
                {{ audit.path }}
              </span>
              <b-tooltip triggers="hover" :target="'resource-path-' + i">
                {{ audit.path }}
              </b-tooltip>
            </td>
            <td>
              <img v-if="audit.success" src="~/assets/img/ok.png" />
              <img v-else src="~/assets/img/ng.png" />
              <button
                v-if="audit.message !== ''"
                class="btn btn-outline-primary"
                @click="showErrorMessage(audit.message)"
              >
                Show
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      audits: [],
      auditErrorMessage: ''
    }
  },
  mounted() {
    this.setAudits()
  },
  methods: {
    async setAudits() {
      const res = await this.$api.AuditGetList(
        this.$store.state.current_project
      )
      if (res.ok) {
        console.log('Audit Events: ', res.data)
        this.audits = res.data
      } else {
        console.log('Failed to get audit event list: %o', res)
      }
    },
    showErrorMessage(msg) {
      this.auditErrorMessage = msg
      this.$refs['view-error-message'].show()
    },
    hideErrorMessage() {
      this.$refs['view-error-message'].hide()
    }
  }
}
</script>

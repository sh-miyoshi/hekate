<template>
  <div class="wrapper content">
    <h3>
      <span class="project">
        {{ this.$store.state.current_project }}
      </span>
      <span class="trush">
        <i class="fa fa-trash" @click="trushConfirm"></i>
      </span>
    </h3>

    <div>
      <b-modal
        id="confirm-delete-project"
        ref="confirm-delete-project"
        title="Confirm"
        cancel-variant="outline-dark"
        ok-variant="danger"
        ok-title="Delete project"
        @ok="trush"
      >
        <p class="mb-0">Are you sure to delete the project ?</p>
      </b-modal>
    </div>

    <div class="form-panel">
      <div>
        <label for="accessTokenLifeSpan" class="col-sm-4 control-label elem">
          Access Token Life Span
        </label>
        <div class="col-sm-3 elem">
          <input
            v-model.number="accessTokenLifeSpan"
            type="number"
            class="form-control"
          />
        </div>
        <div class="col-sm-2 elem">
          <select
            v-model="accessTokenUnit"
            name="accessUnit"
            class="form-control"
          >
            <option v-for="unit in units" :key="unit" :value="unit">
              {{ unit }}
            </option>
          </select>
        </div>
      </div>

      <div>
        <label for="refreshTokenLifeSpan" class="col-sm-4 control-label elem">
          Refresh Token Life Span
        </label>
        <div class="col-sm-3 elem">
          <input
            v-model.number="refreshTokenLifeSpan"
            type="number"
            class="form-control"
          />
        </div>
        <div class="col-sm-2 elem">
          <select
            v-model="refreshTokenUnit"
            name="refreshUnit"
            class="form-control"
          >
            <option v-for="unit in units" :key="unit" :value="unit">
              {{ unit }}
            </option>
          </select>
        </div>
      </div>

      <div class="divider"></div>

      <div v-if="error" class="alert alert-danger">
        {{ error }}
      </div>

      <button class="btn btn-theme" @click="update">Update</button>
    </div>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      units: ['sec', 'minutes', 'hours', 'days'],
      error: '',
      accessTokenLifeSpan: 0,
      accessTokenUnit: 'sec',
      refreshTokenLifeSpan: 0,
      refreshTokenUnit: 'sec',
      signingAlgorithm: ''
    }
  },
  mounted() {
    this.getProject()
  },
  methods: {
    trushConfirm() {
      this.$refs['confirm-delete-project'].show()
    },
    async trush() {
      const res = await this.$api.ProjectDelete(
        this.$store.state.current_project
      )
      console.log('project delete result: %o', res)
      if (!res.ok) {
        this.error = res.message
        return
      }

      alert('successfully deleted.')
      this.$store.commit('setCurrentProject', 'master') // TODO(set correct project name)
      this.$router.push('/home')
    },
    async update() {
      const data = {
        tokenConfig: {
          accessTokenLifeSpan: this.getSpan(
            this.accessTokenLifeSpan,
            this.accessTokenUnit
          ),
          refreshTokenLifeSpan: this.getSpan(
            this.refreshTokenLifeSpan,
            this.refreshTokenUnit
          ),
          signingAlgorithm: 'RS256'
        }
      }
      const res = await this.$api.ProjectUpdate(
        this.$store.state.current_project,
        data
      )
      if (!res.ok) {
        this.error = res.message
        return
      }

      this.getProject()
      alert('successfully updated.')
    },
    async getProject() {
      const res = await this.$api.ProjectGet(this.$store.state.current_project)
      if (!res.ok) {
        this.error = res.message
        return
      }

      let t = this.setUnit(res.data.tokenConfig.accessTokenLifeSpan)
      this.accessTokenLifeSpan = t.span
      this.accessTokenUnit = t.unit
      t = this.setUnit(res.data.tokenConfig.refreshTokenLifeSpan)
      this.refreshTokenLifeSpan = t.span
      this.refreshTokenUnit = t.unit
      this.signingAlgorithm = res.data.tokenConfig.signingAlgorithm
    },
    setUnit(span) {
      let unit = 'sec'
      while (true) {
        switch (unit) {
          case 'sec':
            if (span >= 60 && span % 60 === 0) {
              span /= 60
              unit = 'minutes'
            } else {
              return { span, unit }
            }
            break
          case 'minutes':
            if (span >= 60 && span % 60 === 0) {
              span /= 60
              unit = 'hours'
            } else {
              return { span, unit }
            }
            break
          case 'hours':
            if (span >= 24 && span % 24 === 0) {
              span /= 24
              unit = 'days'
            }
            return { span, unit }
          case 'days':
            return { span, unit }
          default:
            console.log('unexpect unit %s', unit)
            return
        }
      }
    },
    getSpan(span, unit) {
      switch (unit) {
        case 'sec':
          return span
        case 'minutes':
          return span * 60
        case 'hours':
          return span * 60 * 60
        case 'days':
          return span * 60 * 60 * 24
        default:
          console.log('unexpect unit %s', unit)
          break
      }
      return span
    }
  }
}
</script>

<style scoped>
.trush:hover {
  cursor: pointer;
}

.project {
  padding-right: 20px;
}

.elem {
  float: left;
}

.divider {
  clear: both;
  border-bottom: 1px solid #eff2f7;
  padding-bottom: 15px;
  margin-bottom: 15px;
}
</style>

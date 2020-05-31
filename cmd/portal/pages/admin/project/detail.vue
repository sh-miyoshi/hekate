<template>
  <div class="card">
    <div class="card-header">
      <h3>
        <span class="project">
          {{ this.$store.state.current_project }}
        </span>
        <span class="icon">
          <i class="fa fa-trash" @click="deleteProjectConfirm"></i>
        </span>
      </h3>
    </div>

    <div class="card-body">
      <div>
        <b-modal
          id="confirm-delete-project"
          ref="confirm-delete-project"
          title="Confirm"
          cancel-variant="outline-dark"
          ok-variant="danger"
          ok-title="Delete project"
          @ok="deleteProject"
        >
          <p class="mb-0">Are you sure to delete the project ?</p>
        </b-modal>
      </div>

      <div class="form-group row">
        <label class="col-sm-4 control-label">
          User Login URL
        </label>
        <div class="col-sm-7">
          <input
            v-model="loginURL"
            type="text"
            disabled="disabled"
            class="form-control"
          />
        </div>
      </div>

      <div class="form-group">
        <button
          class="btn btn-link dropdown-toggle h5 ml-n2"
          @click="showTokenConfig = !showTokenConfig"
        >
          Token Config
        </button>
        <div v-if="showTokenConfig" class="card-body">
          <div class="form-group row">
            <label class="col-sm-4 control-label">
              Access Token Life Span
            </label>
            <div class="col-sm-3">
              <input
                v-model.number="tokenConfig.accessTokenLifeSpan"
                type="number"
                class="form-control"
              />
            </div>
            <div class="col-sm-2">
              <select
                v-model="tokenConfig.accessTokenUnit"
                name="accessUnit"
                class="form-control"
              >
                <option v-for="unit in units" :key="unit" :value="unit">
                  {{ unit }}
                </option>
              </select>
            </div>
          </div>

          <div class="form-group row">
            <label class="col-sm-4 control-label">
              Refresh Token Life Span
            </label>
            <div class="col-sm-3">
              <input
                v-model.number="tokenConfig.refreshTokenLifeSpan"
                type="number"
                class="form-control"
              />
            </div>
            <div class="col-sm-2">
              <select
                v-model="tokenConfig.refreshTokenUnit"
                name="refreshUnit"
                class="form-control"
              >
                <option v-for="unit in units" :key="unit" :value="unit">
                  {{ unit }}
                </option>
              </select>
            </div>
          </div>

          <div class="form-group row">
            <label class="col-sm-4 control-label">
              Token Signing Algorithm
            </label>
            <div class="col-sm-3">
              <select
                v-model="tokenConfig.signingAlgorithm"
                name="signingAlgorithm"
                class="form-control"
              >
                <option v-for="alg in algs" :key="alg" :value="alg">
                  {{ alg }}
                </option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <div class="form-group">
        <button
          class="btn btn-link dropdown-toggle h5 ml-n3"
          @click="showPasswordPolicy = !showPasswordPolicy"
        >
          Password Policy
        </button>
        <div v-if="showPasswordPolicy" class="card-body">
          <div class="form-group row">
            <label class="col-sm-4 control-label">
              Minimum Length
            </label>
            <div class="col-sm-3">
              <input
                v-model.number="passwordPolicy.minimumLength"
                type="number"
                class="form-control"
              />
            </div>
          </div>
          <div class="form-group row">
            <label class="col-sm-4 control-label">
              Not User Name
            </label>
            <div class="col-sm-3">
              <label
                class="c-switch c-switch-label c-switch-pill c-switch-primary"
              >
                <input
                  class="c-switch-input"
                  type="checkbox"
                  :checked="passwordPolicy.notUserName"
                />
                <span
                  class="c-switch-slider"
                  data-checked="On"
                  data-unchecked="Off"
                ></span>
              </label>
            </div>
          </div>
          <div class="form-group row">
            <label class="col-sm-4 control-label">
              BlackList
            </label>
            <div class="col-sm-3">value</div>
          </div>
          <div class="form-group row">
            <label class="col-sm-4 control-label">
              Includes
            </label>
            <div class="col-md-7 col-form-label">
              <div
                v-if="
                  !passwordPolicy.includes.lower.checked &&
                    !passwordPolicy.includes.upper.checked
                "
                class="form-check checkbox"
              >
                <input
                  class="form-check-input"
                  type="checkbox"
                  :checked="passwordPolicy.includes.caseInsensitive.checked"
                  @change="
                    passwordPolicy.includes.caseInsensitive.checked = !passwordPolicy
                      .includes.caseInsensitive.checked
                  "
                />
                <label class="form-check-label">
                  {{ passwordPolicy.includes.caseInsensitive.name }}
                </label>
              </div>
              <div class="form-check checkbox">
                <input
                  class="form-check-input"
                  type="checkbox"
                  :checked="passwordPolicy.includes.lower.checked"
                  @change="
                    passwordPolicy.includes.lower.checked = !passwordPolicy
                      .includes.lower.checked
                  "
                />
                <label class="form-check-label">
                  {{ passwordPolicy.includes.lower.name }}
                </label>
              </div>
              <div class="form-check checkbox">
                <input
                  class="form-check-input"
                  type="checkbox"
                  :checked="passwordPolicy.includes.upper.checked"
                  @change="
                    passwordPolicy.includes.upper.checked = !passwordPolicy
                      .includes.upper.checked
                  "
                />
                <label class="form-check-label">
                  {{ passwordPolicy.includes.upper.name }}
                </label>
              </div>
              <div class="form-check checkbox">
                <input
                  class="form-check-input"
                  type="checkbox"
                  :checked="passwordPolicy.includes.digit.checked"
                  @change="
                    passwordPolicy.includes.digit.checked = !passwordPolicy
                      .includes.digit.checked
                  "
                />
                <label class="form-check-label">
                  {{ passwordPolicy.includes.digit.name }}
                </label>
              </div>
              <div class="form-check checkbox">
                <input
                  class="form-check-input"
                  type="checkbox"
                  :checked="passwordPolicy.includes.special.checked"
                  @change="
                    passwordPolicy.includes.special.checked = !passwordPolicy
                      .includes.special.checked
                  "
                />
                <label class="form-check-label">
                  {{ passwordPolicy.includes.special.name }}
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="form-group row">
        <label for="refreshTokenLifeSpan" class="col-sm-4 control-label">
          Allow Grant Types
        </label>
        <div class="col-md-7 col-form-label">
          <div
            v-for="type in grantTypes"
            :key="type.value"
            class="form-check checkbox"
          >
            <input
              class="form-check-input"
              type="checkbox"
              :checked="type.checked"
              @change="type.checked = !type.checked"
            />
            <label class="form-check-label">{{ type.name }}</label>
          </div>
        </div>
      </div>

      <div class="card-footer">
        <div v-if="error" class="alert alert-danger">
          {{ error }}
        </div>

        <button class="btn btn-primary" @click="update">Update</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  data() {
    return {
      units: ['sec', 'minutes', 'hours', 'days'],
      algs: ['RS256'],
      error: '',
      loginURL: '',
      showTokenConfig: true,
      tokenConfig: {
        accessTokenLifeSpan: 0,
        accessTokenUnit: 'sec',
        refreshTokenLifeSpan: 0,
        refreshTokenUnit: 'sec',
        signingAlgorithm: 'RS256'
      },
      showPasswordPolicy: true,
      passwordPolicy: {
        minimumLength: 0,
        notUserName: false,
        blackList: [],
        includes: {
          caseInsensitive: {
            name: 'Characters (case-insensitive)',
            checked: false
          },
          lower: { name: 'Lowercase Characters', checked: false },
          upper: { name: 'Uppercase Characters', checked: false },
          digit: { name: 'Digits', checked: false },
          special: { name: 'Special Characters', checked: false }
        }
      },
      grantTypes: [
        {
          name: 'Authorization Code',
          value: 'authorization_code',
          checked: false
        },
        {
          name: 'Client Credentials',
          value: 'client_credentials',
          checked: false
        },
        {
          name: 'Refresh Token',
          value: 'refresh_token',
          checked: false
        },
        {
          name: 'Password',
          value: 'password',
          checked: false
        }
      ]
    }
  },
  mounted() {
    this.setProjectInfo()

    let protcol = 'https'
    if (!process.env.https) {
      protcol = 'http'
    }
    this.loginURL =
      protcol +
      '://' +
      process.env.HEKATE_PORTAL_HOST +
      ':' +
      process.env.HEKATE_PORTAL_PORT +
      '/user/project/' +
      this.$store.state.current_project
  },
  methods: {
    deleteProjectConfirm() {
      this.$refs['confirm-delete-project'].show()
    },
    async deleteProject() {
      const res = await this.$api.ProjectDelete(
        this.$store.state.current_project
      )
      console.log('project delete result: %o', res)
      if (!res.ok) {
        this.error = res.message
        return
      }

      alert('successfully deleted.')
      this.$store.commit('setCurrentProject', 'master')
      this.$router.push('/admin')
    },
    async update() {
      const grantTypes = []
      console.log(this.grantTypes)
      for (const type of this.grantTypes) {
        if (type.checked) {
          grantTypes.push(type.value)
        }
      }

      // TODO(set password policy)

      const data = {
        tokenConfig: {
          accessTokenLifeSpan: this.getSpan(
            this.tokenConfig.accessTokenLifeSpan,
            this.tokenConfig.accessTokenUnit
          ),
          refreshTokenLifeSpan: this.getSpan(
            this.tokenConfig.refreshTokenLifeSpan,
            this.tokenConfig.refreshTokenUnit
          ),
          signingAlgorithm: this.tokenConfig.signingAlgorithm
        },
        allowGrantTypes: grantTypes
      }
      console.log(data)

      const res = await this.$api.ProjectUpdate(
        this.$store.state.current_project,
        data
      )
      if (!res.ok) {
        this.error = res.message
        return
      }

      this.setProjectInfo()
      alert('successfully updated.')
    },
    async setProjectInfo() {
      const res = await this.$api.ProjectGet(this.$store.state.current_project)
      if (!res.ok) {
        this.error = res.message
        return
      }

      let t = this.setUnit(res.data.tokenConfig.accessTokenLifeSpan)
      this.tokenConfig.accessTokenLifeSpan = t.span
      this.tokenConfig.accessTokenUnit = t.unit
      t = this.setUnit(res.data.tokenConfig.refreshTokenLifeSpan)
      this.tokenConfig.refreshTokenLifeSpan = t.span
      this.tokenConfig.refreshTokenUnit = t.unit
      this.tokenConfig.signingAlgorithm = res.data.tokenConfig.signingAlgorithm

      // TODO(set password policy)

      // set allow grant types
      for (const type of res.data.allowGrantTypes) {
        for (const t of this.grantTypes) {
          if (type === t.value) {
            t.checked = true
          }
        }
      }
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
.project {
  padding-right: 20px;
}
</style>

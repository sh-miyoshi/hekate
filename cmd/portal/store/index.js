export const state = () => ({
  current_project: 'master',
  login_state: ''
})

export const mutations = {
  setCurrentProject(state, project) {
    state.current_project = project
  },
  setLoginState(state, loginState) {
    state.login_state = loginState
  }
}

export const state = () => ({
  current_project: process.env.LOGIN_PROJECT
})

export const mutations = {
  setCurrentProject(state, project) {
    state.current_project = project
  }
}

export const state = () => ({
  current_project: ''
})

export const mutations = {
  setCurrentProject(state, project) {
    state.current_project = project
  }
}

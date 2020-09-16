export const state = () => ({
  current_project: 'master'
})

export const mutations = {
  setCurrentProject(state, project) {
    state.current_project = project
  }
}

export const state = () => ({
  projects: ['master', 'newproject'],
  select_project: 'master'
})

export const getters = {
  getProjects: (state) => {
    return state.projects
  }
}

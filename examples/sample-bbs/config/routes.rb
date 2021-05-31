Rails.application.routes.draw do
  get 'bbs/index'
  post 'bbs/login'
  get 'bbs/show'
  post 'bbs/add'
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html

  root 'bbs#index'
end

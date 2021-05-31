class BbsController < ApplicationController
  def index; end

  def callback; end

  def show
    @messages = Message.all
  end

  def add
    Message.create(
      text: params[:text],
      userid: 0 # debug
    )

    redirect_to action: 'show'
  end
end
